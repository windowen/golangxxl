package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	constsR "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/redisdb/redis"
	rpcClient "queueJob/pkg/rpcclient"
	"queueJob/pkg/tools/cast"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	Rv9 "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"queueJob/pkg/constant"
	"queueJob/pkg/queue"
	"queueJob/pkg/rocketmq"
	"queueJob/pkg/zlogger"
)

var liveRoomRobot = &liveRoomRobotDelay{}

type liveRoomRobotDelay struct{}

func (ur *liveRoomRobotDelay) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		zlogger.Debugw("liveRoomRobotDelay Received message", zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))

		data := &queue.LiveRoomRobotDelayReq{}
		if err := json.Unmarshal(msg.Body, data); err != nil {
			zlogger.Errorw("liveRoomRobotDelay::handleMessages, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
			continue
		}

		// 获取直播间缓存
		roomCacheInfo, err := getRoomCache(data.RoomId)
		if err != nil {
			zlogger.Errorf("liveRoomRobotDelay getRoomCache |roomId:%v| err: %v", data.RoomId, err)
			continue
		}

		if roomCacheInfo.SceneHistoryId != data.SceneId || roomCacheInfo.LiveStatus != constant.RoomSceneStatusDo {
			zlogger.Errorf("liveRoomRobotDelay |sceneHistoryId:%v,sceneId:%v| err: this show has been downloaded", roomCacheInfo.SceneHistoryId, data.SceneId)
			continue
		}

		// 是否配置机器人
		robotCfg := GetRobotConfig(data.RoomId)
		if robotCfg == nil {
			zlogger.Errorf("liveRoomRobotDelay |roomId:%v| err : Robot not configured", data.RoomId)
			continue
		}

		// 继续投递队列
		rocketmq.PublishWithDelayJson(
			rocketmq.LiveRoomRobotDelay,
			data,
			1,
		)

		// 获取直播间机器人列表
		robotList, err := redis.ZRangeWithScores(fmt.Sprintf(constsR.RoomRobotListZSet, data.RoomId))
		if err != nil {
			zlogger.Errorf("liveRoomRobotDelay ZRangeWithScores |roomId:%v| err : %v", data.SceneId, err)
			continue
		}

		// 获取机器人数量
		count, err := redis.ZCARD(fmt.Sprintf(constsR.RoomRobotListZSet, data.RoomId))
		if err != nil {
			zlogger.Infof("liveRoomRobotDelay ZCARD |roomId:%v| err : %v", data.SceneId, err)
			continue
		}

		// 检查机器人是否到期
		for _, robot := range robotList {
			expiredTime := cast.ToInt64(robot.Score)
			if time.Now().After(time.Unix(expiredTime, 0)) {
				err := ur.quitRoom(ctx, cast.ToInt(robot.Member), roomCacheInfo, robotCfg.QuitLessenViewerCount)
				if err != nil {
					zlogger.Errorf("liveRoomRobotDelay quitRoom |roomId:%v,userId:%v| err : %v", data.RoomId, robot.Member, err)
				}
			}
		}

		// 检查是否加入新机器人
		if cast.ToInt(count) < robotCfg.RoomMaxRobots {
			lastJoinTime := cast.ToInt64(time.Now().Add(-time.Second))
			if roomCacheInfo.RobotLastJoinTime > constant.Zero {
				lastJoinTime = cast.ToInt64(roomCacheInfo.RobotLastJoinTime)
			}

			if time.Now().After(time.Unix(lastJoinTime, 0)) {
				err := ur.joinRoom(ctx, roomCacheInfo, robotCfg)
				if err != nil {
					zlogger.Errorf("liveRoomRobotDelay joinRoom |roomId:%v| err : %v", data.RoomId, err)
				}
			}
		}
	}

	return consumer.ConsumeSuccess, nil
}

// 退出房间
func (ur *liveRoomRobotDelay) quitRoom(ctx context.Context, robotId int, roomCache *RoomCacheInfo, isNotify bool) error {
	// 回收机器人
	err := redis.ZRem(fmt.Sprintf(constsR.RoomRobotListZSet, roomCache.Id), cast.ToString(robotId))
	if err != nil {
		return err
	}

	_, err = redis.SAdd(constsR.RoomRobotSet, robotId)
	if err != nil {
		return err
	}

	// 发送通知
	if isNotify {
		err = rpcClient.ServiceClientsInstance.LiveClient.RobotQuitChatRoom(ctx, robotId, roomCache.ChatRoomId)
		if err != nil {
			return err
		}
	}

	return nil
}

// 加入房间
func (ur *liveRoomRobotDelay) joinRoom(ctx context.Context, roomCache *RoomCacheInfo, robtCfg *RobotConfig) error {
	// 随机获取机器人
	robot, err := RandomAvailableRobot()
	if err != nil {
		zlogger.Errorf("joinRoom RandomAvailableRobot |roomId:%v | err: %v", roomCache.Id, err)
		return err
	}

	// 加入房间
	stayTime := time.Now().Add(ur.randomTime(robtCfg.MinStayTime, robtCfg.MaxStayTime)).Unix()
	_, err = redis.ZAdd(fmt.Sprintf(constsR.RoomRobotListZSet, roomCache.Id), Rv9.Z{
		Score:  cast.ToFloat64(stayTime),
		Member: robot.UserId,
	})
	if err != nil {
		zlogger.Errorf("joinRoom ZAdd |roomId:%v | err: %v", roomCache.Id, err)
		return err
	}

	// 更新直播间机器人加入时间
	robotLastJoinTime := time.Now().Add(ur.randomTime(robtCfg.MinJoinInterval, robtCfg.MaxJoinInterval)).Unix()
	roomCache.RobotLastJoinTime = cast.ToInt(robotLastJoinTime)

	err = setRoomCache(roomCache.Id, roomCache)
	if err != nil {
		zlogger.Errorf("joinRoom setRoomCache |roomId:%v | err: %v", roomCache.Id, err)
		return err
	}

	// 发送通知
	err = rpcClient.ServiceClientsInstance.LiveClient.RobotJoinChatRoom(ctx, robot.UserId, roomCache.ChatRoomId)
	if err != nil {
		zlogger.Errorf("joinRoom RobotJoinChatRoom |roomId:%v | err: %v", roomCache.Id, err)
		return err
	}

	return nil
}

// 获取随机时间
func (ur *liveRoomRobotDelay) randomTime(min, max int) time.Duration {
	rand.NewSource(time.Now().UnixNano())
	randomSeconds := min + rand.Intn(max-min+1)

	return time.Duration(randomSeconds) * time.Second
}
