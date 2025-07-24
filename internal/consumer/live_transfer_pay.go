package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"

	"liveJob/pkg/agora"
	"liveJob/pkg/agora/model"
	"liveJob/pkg/constant"
	constsR "liveJob/pkg/constant/redis"
	"liveJob/pkg/db/redisdb/redis"
	"liveJob/pkg/queue"
	"liveJob/pkg/tools/cast"
	"liveJob/pkg/zlogger"
)

var liveRoomTransferPay = &liveTransferPay{}

type liveTransferPay struct{}

func (ur *liveTransferPay) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		zlogger.Debugw("LiveTransferPay Received message", zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))

		data := &queue.LiveRoomTransferPayDelay{}
		if err := json.Unmarshal(msg.Body, data); err != nil {
			zlogger.Errorf("LiveTransferPay |msgID:%v| err: %v", msg.MsgId, err)
			continue
		}

		// 加锁防止重复处理
		lockSign := fmt.Sprintf("room_start_delay_%d_%d_%d", data.RoomId, data.AnchorId, data.SceneId)
		isLock, retFun := tryGetDistributedLock(lockSign, lockSign, 10000, 10000)
		if !isLock {
			zlogger.Errorf("LiveTransferPay tryGetDistributedLock |roomId:%v,sceneId:%v| err: failed to acquire lock", data.RoomId, data.SceneId)
			continue
		}

		// 获取直播间缓存
		roomCacheInfo, err := getRoomCache(data.RoomId)
		if err != nil {
			// 释放锁
			retFun()
			zlogger.Errorf("LiveTransferPay getRoomCache |roomId:%v| err: %v", data.RoomId, err)
			continue
		}

		if roomCacheInfo.SceneHistoryId != data.SceneId || roomCacheInfo.LiveStatus != constant.RoomSceneStatusDo {
			// 释放锁
			retFun()
			zlogger.Infof("LiveTransferPay |sceneHistoryId:%v,sceneId:%v| this show has been downloaded", roomCacheInfo.SceneHistoryId, data.SceneId)
			continue
		}

		// 查询直播间使用
		ids, err := redis.SMembers(fmt.Sprintf(constsR.SceneUsersSet, data.RoomId))
		if err != nil {
			// 释放锁
			retFun()
			zlogger.Errorf("LiveTransferPay SceneUsersSet |roomId:%v,sceneId:%v| err: %v", data.RoomId, data.SceneId, err)
			continue
		}

		if len(ids) == 0 {
			// 释放锁
			retFun()
			continue
		}

		for _, userId := range ids {
			// 获取扣费时间
			payTime, err := redis.ZScore(fmt.Sprintf(constsR.ScenePayUsers, data.RoomId), cast.ToString(userId))
			if err != nil {
				zlogger.Errorf("LiveTransferPay ZScore ScenePayUsers |roomId:%v,userId:%v| err: %v", roomCacheInfo.Id, userId, err)
				continue
			}

			// 60内未扣费
			if cast.ToInt64(payTime) < time.Now().Add(-60*time.Second).Unix() {
				// 踢出直播
				err := agora.RtcClientInstance.RtcKickOutUser(model.RtcKickOutUserReq{
					UserId:   cast.ToInt(userId),
					RoomId:   data.RoomId,
					Duration: 10, // 临时踢出10秒
				})
				if err != nil {
					zlogger.Errorf("LiveTransferPay RtcKickOutUser | err: %v", err)
					continue
				}
			}
		}

		// 释放锁
		retFun()
	}

	return consumer.ConsumeSuccess, nil
}
