package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	goRedis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"queueJob/pkg/agora"
	"queueJob/pkg/agora/model"
	"queueJob/pkg/constant"
	constsR "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/queue"
	"queueJob/pkg/rocketmq"
	rpcClient "queueJob/pkg/rpcclient"
	"queueJob/pkg/tools/cast"
	"queueJob/pkg/tools/errs"
	"queueJob/pkg/tools/utils"
	"queueJob/pkg/zlogger"
)

var liveMPaid = &liveMinutePaid{}

type liveMinutePaid struct{}

func (lmp *liveMinutePaid) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	var timeNow = time.Now().Unix()
	for _, msg := range msgList {
		lmp.minutePay(ctx, msg, timeNow)
	}

	return consumer.ConsumeSuccess, nil
}

func (lmp *liveMinutePaid) minutePay(ctx context.Context, msg *primitive.MessageExt, now int64) {
	zlogger.Debugw("LiveMinutePaid Received message", zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))

	data := &queue.LiveRoomUserMinuteDelayPaid{}
	if err := json.Unmarshal(msg.Body, data); err != nil {
		zlogger.Errorw("liveMinutePaid::handleMessages, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
		return
	}

	// 获取直播间缓存
	roomCacheInfo, err := getRoomCache(data.RoomId)
	if err != nil {
		zlogger.Errorf("liveMinutePaid getRoomCache |roomId:%v| err: %v", data.RoomId, err)
		return
	}

	// 获取主播信息缓存
	anchorCacheInfo, err := getUserCache(data.AnchorId)
	if err != nil {
		zlogger.Errorf("liveMinutePaid getUserCache |roomId:%v| err: %v", data.RoomId, err)
		return
	}

	// 是否开播
	if roomCacheInfo.SceneHistoryId != data.SceneId || roomCacheInfo.LiveStatus != constant.RoomSceneStatusDo {
		zlogger.Infof("liveMinutePaid |sceneHistoryId:%v,sceneId:%v| this show has been downloaded", roomCacheInfo.SceneHistoryId, data.SceneId)
		return
	}

	// 用户是否在直播间
	if !lmp.isUserInRoom(data) {
		return
	}

	// 是否房管
	if lmp.isRoomManager(data) {
		return
	}

	// 继续投递每分钟扣款队列
	rocketmq.PublishWithDelayJson(rocketmq.LiveRoomStartFeeLive, data, 5)

	// 60秒内是否已支付
	if !lmp.isRepeatedDeduction(data, now) {
		return
	}

	// 扣费
	if err = rpcClient.ServiceClientsInstance.FinanceClient.RoomLiveMinuteDelayPaid(ctx,
		data.UserId,
		roomCacheInfo.UnitPrice,
		roomCacheInfo.Id,
		roomCacheInfo.UserId); err != nil {
		zlogger.Errorf("liveMinutePaid RoomLiveMinuteDelayPaid |roomId:%v,sceneId:%v,unitPrice:%v,userId:%v| err: %v", data.RoomId, data.SceneId, roomCacheInfo.UnitPrice, data.UserId, err)

		if errs.RpcErrCheck(err, errs.ErrDiamondNotEnough) {
			zlogger.Debugw("diamond not enough", zap.Int("uid", data.UserId))
			// 余额不足踢出直播间
			if err = agora.RtcClientInstance.RtcKickOutUser(model.RtcKickOutUserReq{
				UserId:   data.UserId,
				RoomId:   data.RoomId,
				Duration: 10, // 临时踢出10秒
			}); err != nil {
				zlogger.Errorf("liveMinutePaid RtcKickOutUser | err: %v", err)
				return
			}

			// 发送通知
			if err = rpcClient.ServiceClientsInstance.LiveClient.InsufficientBalance(ctx, data.UserId, roomCacheInfo.ChatRoomId); err != nil {
				zlogger.Errorf("liveMinutePaid InsufficientBalance | err: %v", err)
				return
			}
		}
		zlogger.Debugw("minute delay error", zap.Int("uid", data.UserId), zap.Error(err))

		return
	}

	// 发送收益通知
	if err = rpcClient.ServiceClientsInstance.LiveClient.LiveMinutePaidIncomeNotify(ctx,
		roomCacheInfo.Id,
		roomCacheInfo.UserId,
		data.UserId,
		roomCacheInfo.UnitPrice,
		roomCacheInfo.ChatRoomId); err != nil {
		zlogger.Errorf("liveMinutePaid LiveMinutePaidIncomeNotify | err: %v", err)
	}

	if _, err = redis.ZAdd(fmt.Sprintf(constsR.ScenePayUsers, roomCacheInfo.Id), goRedis.Z{
		Score:  float64(now),
		Member: data.UserId,
	}); err != nil {
		zlogger.Errorf("liveMinutePaid ZAdd |roomId:%v,userId:%v| err: %v", roomCacheInfo.Id, data.UserId, err)
	}

	// 发送消费纪录消息到队列
	rocketmq.PublishJson(rocketmq.LiveRoomSendGift, &queue.LiveRoomPayDiamond{
		BillNo:       utils.GetOrderNo("LR"),
		UserId:       data.UserId,
		RoomId:       roomCacheInfo.Id,
		FamilyId:     anchorCacheInfo.FamilyId,
		AnchorId:     roomCacheInfo.UserId,
		SceneId:      roomCacheInfo.SceneHistoryId,
		Category:     constant.IncomeTypeLive,
		ProjectId:    constant.No,
		ProjectNum:   1,
		UnitPrice:    roomCacheInfo.UnitPrice,
		ProjectTotal: roomCacheInfo.UnitPrice,
		IsDivide:     constant.Yes,
	})
}

// 是否离开直播间
func (lmp *liveMinutePaid) isUserInRoom(data *queue.LiveRoomUserMinuteDelayPaid) bool {
	isEx, err := redis.SIsMember(fmt.Sprintf(constsR.SceneUsersSet, data.RoomId), cast.ToString(data.UserId))
	if err != nil {
		zlogger.Errorf("isUserInRoom RoomManageCacheKey |roomId:%v,userId:%v| err: %v", data.RoomId, data.UserId, err)
		return false
	}

	if !isEx {
		zlogger.Debugf("isUserInRoom|roomId:%v,userId:%v| err: the user has left the live broadcast room", data.RoomId, data.UserId)
		return false
	}

	return true
}

// 是否房管
func (lmp *liveMinutePaid) isRoomManager(data *queue.LiveRoomUserMinuteDelayPaid) bool {
	isRoomManage, err := redis.SIsMember(fmt.Sprintf(constsR.RoomManageCacheKey, data.RoomId), cast.ToString(data.UserId))
	if err != nil {
		zlogger.Errorf("isRoomManager RoomManageCacheKey |roomId:%v,userId:%v| err: %v", data.RoomId, data.UserId, err)
		return false
	}

	if isRoomManage {
		zlogger.Infof("isRoomManager RoomManageCacheKey |roomId:%v,userId:%v| user is the house manager", data.RoomId, data.UserId)
		return true
	}

	return false
}

// 60秒内是否已支付
func (lmp *liveMinutePaid) isRepeatedDeduction(data *queue.LiveRoomUserMinuteDelayPaid, timeNow int64) bool {
	payTime, err := redis.ZScore(fmt.Sprintf(constsR.ScenePayUsers, data.RoomId), cast.ToString(data.UserId))
	if err != nil {
		zlogger.Errorf("isRepeatedDeduction ScenePayUsers |roomId:%v,userId:%v| err: %v", data.RoomId, data.UserId, err)
		return false
	}

	ret := timeNow - int64(payTime)
	if ret >= 60 {
		return true
	}

	zlogger.Infow("liveMinutePaid", zap.Int64("ret", ret), zap.Int("uid", data.UserId))
	return false
}
