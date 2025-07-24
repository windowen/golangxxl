package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"

	"liveJob/pkg/db/mysql"
	"liveJob/pkg/db/table"
	"liveJob/pkg/queue"
	"liveJob/pkg/rocketmq"
	"liveJob/pkg/tools/utils"
	"liveJob/pkg/zlogger"
)

var (
	gift = &giftConsumer{}
)

type giftConsumer struct{}

func (o *giftConsumer) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		zlogger.Infow("gift message", zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))
		o.sendGift(ctx, msg)
	}
	return consumer.ConsumeSuccess, nil
}

func (o *giftConsumer) sendGift(ctx context.Context, msg *primitive.MessageExt) {
	jsonData := &queue.LiveRoomPayDiamond{}
	if err := json.Unmarshal(msg.Body, jsonData); err != nil {
		zlogger.Errorw("giftConsumer::sendGift, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
		return
	}

	details := &table.LiveRoomIncomeDetails{
		BillNo:       jsonData.BillNo,
		UserId:       jsonData.UserId,
		RoomId:       jsonData.RoomId,
		FamilyId:     jsonData.FamilyId,
		AnchorId:     jsonData.AnchorId,
		SceneId:      jsonData.SceneId,
		Category:     jsonData.Category,
		ProjectId:    jsonData.ProjectId,
		ProjectNum:   jsonData.ProjectNum,
		UnitPrice:    jsonData.UnitPrice,
		ProjectTotal: jsonData.ProjectTotal,
		IsDivide:     jsonData.IsDivide,
	}

	userCache, err := getUserCache(details.AnchorId)
	if err != nil {
		zlogger.Errorw("giftConsumer::sendGift, get user cache error", zap.String("msgID", msg.MsgId), zap.Error(err))
		return
	}
	details.CountryCode = userCache.CountryCode
	if userCache.IsFamilyMaster == 1 {
		details.FamilyMasterId = details.AnchorId
	} else {
		details.FamilyMasterId = userCache.FamilyMasterId
	}

	roomCache, err := getRoomCache(details.RoomId)
	if err != nil {
		zlogger.Errorw("giftConsumer::sendGift, get room cache error", zap.String("msgID", msg.MsgId), zap.Error(err))
		return
	}
	details.GiftRatio = roomCache.GiftRatio
	details.PlatformRatio = roomCache.PlatformRatio
	details.FamilyRatio = roomCache.FamilyRatio

	total := details.ProjectTotal * 10000

	if details.IsDivide == 1 { // 需要分成
		// 计算并分配收入
		assignIncome := func(ratio int, remaining *int) int64 {
			// ratio 是乘过100的数值, 始终保持原数值*10000
			income := ratio * details.ProjectTotal * 100
			if *remaining-income <= 0 {
				income = *remaining
				*remaining = 0
			} else {
				*remaining -= income
			}
			return int64(income)
		}

		// 分配主播收入
		details.AnchorIncome = assignIncome(details.GiftRatio, &total)

		// 分配平台收入
		details.PlatformIncome = assignIncome(details.PlatformRatio, &total)

		// 分配家族收入
		details.FamilyIncome = assignIncome(details.FamilyRatio, &total)

		// 如果有剩余收入，增加到平台收入
		if total > 0 {
			details.PlatformIncome += int64(total)
		}
	} else {
		details.AnchorIncome = int64(total)
	}

	details.CreatedAt = time.Now()
	details.BillSerial = utils.GenerateSerialNumber(details.Category)

	if err = mysql.LiveDB.WithContext(ctx).Create(details).Error; err != nil {
		zlogger.Errorw("giftConsumer::sendGift, insert detail table fail", zap.String("msgID", msg.MsgId), zap.Any("detail", details), zap.Error(err))
		return
	}

	rocketmq.PublishJson(rocketmq.StreamerReceiveDiamond, &queue.StreamerReceiveDiamond{
		AnchorId:     jsonData.AnchorId,
		RoomId:       jsonData.RoomId,
		ProjectTotal: jsonData.ProjectTotal,
	})
}
