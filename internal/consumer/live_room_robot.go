package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"

	"go.uber.org/zap"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"queueJob/pkg/db/redisdb/redis"

	"queueJob/pkg/constant"
	constsR "queueJob/pkg/constant/redis"
	"queueJob/pkg/queue"
	"queueJob/pkg/rocketmq"
	"queueJob/pkg/tools/utils"
	"queueJob/pkg/zlogger"
)

var liveRoomRobot = &liveRoomRobotDelay{}

type liveRoomRobotDelay struct{}

func (ur *liveRoomRobotDelay) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			tmpStr := fmt.Sprintf("err=%v panic ==> %s\n", err, string(buf[:n]))
			zlogger.Error(tmpStr) // 记录到日志
		}
	}()

	for _, msg := range msgList {
		ur.robotDelay(ctx, msg)
	}

	return consumer.ConsumeSuccess, nil
}

func (ur *liveRoomRobotDelay) robotDelay(ctx context.Context, msg *primitive.MessageExt) {
	data := &queue.LiveRoomRobotDelayReq{}
	if err := json.Unmarshal(msg.Body, data); err != nil {
		zlogger.Errorw("liveRoomRobotDelay::handleMessages, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
		return
	}

	// 获取直播间缓存
	roomCacheInfo, err := redis.GetRoomCache(data.RoomId)
	if err != nil {
		zlogger.Errorw("liveRoomRobotDelay getRoomCache", zap.Int("roomId", data.RoomId), zap.Error(err))
		return
	}

	if roomCacheInfo.SceneHistoryId != data.SceneId || roomCacheInfo.LiveStatus != constant.RoomSceneStatusDo {
		return
	}

	// 是否配置机器人
	robotCfg := GetRobotConfig(data.RoomId)
	if robotCfg == nil {
		zlogger.Errorw("liveRoomRobotDelay Robot not configured", zap.Int("roomId", data.RoomId))
		return
	}

	_ = redis.DelKey(fmt.Sprintf(constsR.RoomRobotNum, roomCacheInfo.Id))

	stayTime := utils.Random(robotCfg.MinStayTime, robotCfg.MaxStayTime)
	enterTime := utils.Random(robotCfg.MinJoinInterval, robotCfg.MaxJoinInterval)

	rocketmq.PublishWithSeconds(rocketmq.LiveRoomRobotEnter, &queue.LiveRoomEnterReq{
		RoomId:    data.RoomId,
		SceneId:   data.SceneId,
		LeaveTime: stayTime,
	}, enterTime)
}
