package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	redis2 "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/db/table"
	"queueJob/pkg/queue"
	"queueJob/pkg/zlogger"
)

var bChange = &balanceChange{}

type balanceChange struct{}

func (o *balanceChange) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		o.consumeMessage(ctx, msg)
	}
	return consumer.ConsumeSuccess, nil
}

func (o *balanceChange) consumeMessage(ctx context.Context, msg *primitive.MessageExt) {
	zlogger.Debugw("balanceChange::consumeMessage, message",
		zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))

	moneyChange := &queue.MoneyChange{}
	if err := json.Unmarshal(msg.Body, moneyChange); err != nil {
		zlogger.Errorw("balanceChange::consumeMessage, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
		return
	}

	dbTable := &table.MoneyChange{
		UserId:       moneyChange.UserId,
		CountryCode:  moneyChange.CountryCode,
		CountryName:  moneyChange.CountryName,
		GameOrderNo:  moneyChange.GameOrderNo,
		ChangeType:   moneyChange.ChangeType,
		ChangeAmount: moneyChange.ChangeAmount,
		BeforeAmount: moneyChange.BeforeAmount,
		AfterAmount:  moneyChange.AfterAmount,
		ExchangeRate: moneyChange.ExchangeRate,
		Remark:       moneyChange.Remark,
		GameProvider: moneyChange.GameProvider,
		GameType:     moneyChange.GameType,
		GameName:     moneyChange.GameName,
		TradeType:    moneyChange.TradeType,
		WagerCode:    moneyChange.WagerCode,
		CreatedAt:    moneyChange.CreatedAt,
	}

	if err := mysql.LiveDB.WithContext(ctx).Create(dbTable).Error; err != nil {
		zlogger.Errorw("balanceChange::consumeMessage, insert money change table fail", zap.String("msgID", msg.MsgId), zap.Any("dbTable", dbTable), zap.Error(err))
		return
	}

	key := fmt.Sprintf(redis2.UserWalletCache, moneyChange.UserId)
	balance, err := redis.HGet(key, "balance")
	if err != nil {
		zlogger.Errorw("balanceChange::consumeMessage get redis cache error", zap.Int("uid", moneyChange.UserId), zap.Error(err))
		return
	}

	curBalance, err := decimal.NewFromString(balance)
	if err != nil {
		zlogger.Errorw("balanceChange::consumeMessage parse balance error", zap.Int("uid", moneyChange.UserId), zap.String("balance", balance), zap.Error(err))
		return
	}

	curBalance = curBalance.Mul(decimal.NewFromFloat(0.01)).Truncate(2)

	// 更新余额
	dbResult := mysql.LiveDB.WithContext(ctx).Model(&table.SiteUserWallet{}).
		Where("user_id = ?", moneyChange.UserId).
		Update("balance", curBalance)
	if dbResult.Error != nil {
		zlogger.Errorw("balanceChange::consumeMessage, update user wallet info error",
			zap.Int("userId", moneyChange.UserId), zap.Error(dbResult.Error))
		return
	}
}
