package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"

	"liveJob/pkg/constant/redis"
	"liveJob/pkg/db/mysql"
	"liveJob/pkg/db/table"
	"liveJob/pkg/queue"
	"liveJob/pkg/zlogger"
)

var fCancel = &financeCancel{}

type financeCancel struct{}

func (o *financeCancel) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		o.cancelFinance(ctx, msg)
	}
	return consumer.ConsumeSuccess, nil
}

func (o *financeCancel) cancelFinance(ctx context.Context, msg *primitive.MessageExt) {
	zlogger.Debugw("finance::handleMessages, financeCancel message",
		zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))

	pMessage := &queue.PayInCache{}
	if err := json.Unmarshal(msg.Body, pMessage); err != nil {
		zlogger.Errorw("finance::handleMessages, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
		return
	}

	// 加锁防止重复结算
	lockSign := fmt.Sprintf(redis.UserPay, pMessage.UserId)
	isLock, retFun := tryGetDistributedLock(lockSign, lockSign, 10000, 10000)
	if !isLock {
		zlogger.Errorw("lock failed", zap.String("lock_sign", lockSign), zap.Int("uid", pMessage.UserId))
		return
	}

	defer retFun()

	dbRecord := &table.FinancePayRecord{}
	if err := mysql.LiveDB.WithContext(ctx).Where("bill_no", pMessage.OrderId).First(dbRecord).Error; err != nil {
		zlogger.Errorw("finance::handleMessages, get user record failed",
			zap.Int("uid", pMessage.UserId), zap.String("orderId", pMessage.OrderId), zap.Error(err))
		return
	}

	// 判断订单是否在处理中
	if dbRecord.Status != 1 {
		zlogger.Errorw("finance::handleMessages, check user record failed", zap.Int("uid", pMessage.UserId))
		return
	}

	if pMessage.UserId != dbRecord.UserId {
		zlogger.Errorw("finance::handleMessages, check user record failed", zap.Int("uid", pMessage.UserId))
		return
	}

	if dbResult := mysql.LiveDB.WithContext(ctx).Model(&table.FinancePayRecord{}).Where("id = ?", dbRecord.Id).Updates(map[string]interface{}{
		"updated_at": time.Now(),
		"status":     5,
	}); dbResult.Error != nil {
		zlogger.Errorw("finance::handleMessages, update pay record error", zap.Int("uid", pMessage.UserId), zap.Error(dbResult.Error))
		return
	}

	zlogger.Debugw("finance::handleMessages, update db record status", zap.Int("uid", pMessage.UserId), zap.String("orderId", pMessage.OrderId))
}
