package kafkaconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"

	constsR "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/redisdb/redis"

	"queueJob/pkg/db/table"
	"queueJob/pkg/queue"
	"queueJob/pkg/zlogger"
)

var checkBonusMessage = &checkBonus{}

type checkBonus struct{}

func (o *checkBonus) handleMessages(msg []byte) {
	ctx := context.Background()
	pMessage := &queue.CheckRegisterBonus{}
	if err := json.Unmarshal(msg, pMessage); err != nil {
		zlogger.Errorw("checkBonusMessage::handleMessages, unmarshal msg fail", zap.Error(err))
		return
	}

	zlogger.Debugw("checkBonusMessage::handleMessages", zap.Stringer("bonus check", pMessage))

	// 加锁同步操作用户钱包数据
	lockSign := fmt.Sprintf(constsR.UserPay, pMessage.UserId)
	isLock, retFun := redis.TryGetDistributedLock(lockSign, lockSign, 10000, 10000)
	if !isLock {
		zlogger.Errorw("lock failed", zap.String("lock_sign", lockSign), zap.Int("uid", pMessage.UserId))
		return
	}

	defer retFun()

	userId := pMessage.UserId
	// 检查是否满足充值条件
	config := &table.ActRegisterWelfareConfig{}
	if err := mysql.LiveDB.WithContext(ctx).First(config).Error; err != nil {
		zlogger.Errorw("checkBonusMessage::handleMessages, get ActRegisterWelfareConfig error",
			zap.Int("userId", userId),
			zap.Error(err),
		)
		return
	}

	// 查询钱包
	dbWallet := &table.SiteUserWallet{}
	if err := mysql.LiveDB.WithContext(ctx).Where("user_id = ?", userId).First(dbWallet).Error; err != nil {
		zlogger.Errorw("checkBonusMessage::handleMessages, get user wallet record failed",
			zap.Int("uid", userId), zap.Error(err))
		return
	}

	if dbWallet.TotalRecharge.LessThan(decimal.NewFromInt(int64(config.RechargeAmount))) {
		return
	}

	bonusWallet := &table.ActBonusWallet{}
	if err := mysql.LiveDB.WithContext(ctx).Where("user_id = ?", userId).First(bonusWallet).Error; err != nil {
		zlogger.Errorw("checkBonusMessage::handleMessages, get ActBonusWallet error",
			zap.Int("userId", userId),
			zap.Error(err),
		)
		return
	}

	if !(bonusWallet.FlowAmount.IsZero() && bonusWallet.ReachUsers == 0) &&
		bonusWallet.Balance.LessThan(config.WithdrawThresholdAmount) {
		zlogger.Infow("checkBonusMessage::handleMessages, check info",
			zap.Int("userId", userId),
			zap.Stringer("flowAmount", bonusWallet.FlowAmount),
			zap.Int("reachUsers", bonusWallet.ReachUsers),
			zap.Stringer("balance", bonusWallet.Balance),
			zap.Stringer("withdrawAmount", config.WithdrawThresholdAmount),
		)
		return
	}

	// 以下逻辑是把活动钱包的金额转移到用户余额中
	if innerErr := redis.DelKey(fmt.Sprintf(constsR.UserBonusWalletCache, userId)); innerErr != nil {
		zlogger.Errorw("checkBonusMessage::handleMessages, update user bonus wallet cache error",
			zap.Int("userId", userId), zap.Error(innerErr))
		return
	}

	// 更新钱包缓存
	incrResult, innerRet := redis.HIncrBy(fmt.Sprintf(constsR.UserWalletCache, userId),
		"balance", bonusWallet.Balance.Mul(decimal.NewFromInt(10000)).IntPart())
	if innerRet != nil {
		zlogger.Errorw("checkBonusMessage::handleMessages, update user wallet cache err",
			zap.Int("userId", userId), zap.Error(innerRet))
		return
	}

	curBalance := decimal.NewFromInt(incrResult).Mul(decimal.NewFromFloat(0.0001)).Truncate(4)

	parentCache, innerErr := redis.GetUserCache(userId)
	if innerErr != nil {
		zlogger.Errorw("checkBonusMessage::handleMessages, get user cache error",
			zap.Int("uid", userId), zap.Error(innerErr))
		return
	}

	nowTime := time.Now()
	err := mysql.LiveDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? ", userId).Model(&table.ActBonusWallet{}).
			Updates(map[string]interface{}{
				"balance":     decimal.Zero,
				"flow_amount": decimal.Zero,
			}).Error; err != nil {
			return err
		}

		// 更新余额
		if err := tx.Model(&table.SiteUserWallet{}).
			Where("user_id = ?", userId).
			Updates(map[string]interface{}{"balance": curBalance}).
			Error; err != nil {
			zlogger.Errorw("checkBonusMessage::handleMessages, update user wallet info error",
				zap.Int("userId", userId), zap.Error(err))
			return err
		}

		zlogger.Infow("checkBonusMessage::handleMessages, update user wallet info",
			zap.Int("userId", userId), zap.Stringer("curBalance", curBalance))

		// 活动彩金转余额账变
		dbTable := &table.MoneyChange{
			UserId:          userId,
			CountryCode:     parentCache.CountryCode,
			CountryName:     parentCache.Country,
			ChangeType:      7, // 注册活动彩金转余额
			ChangeAmount:    bonusWallet.Balance,
			BeforeAmount:    dbWallet.Balance,
			AfterAmount:     curBalance,
			AfterFlowAmount: dbWallet.FlowAmount,
			ActAfterAmount:  decimal.Zero,
			CreatedAt:       nowTime,
			CreatedDay:      nowTime.Format(time.DateOnly),
		}
		if err := tx.Table(dbTable.TableName()).Create(dbTable).Error; err != nil {
			zlogger.Errorw("checkBonusMessage::handleMessages, insert money change table fail",
				zap.Any("dbTable", dbTable), zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		zlogger.Errorw("checkBonusMessage::handleMessages, catch error", zap.Int("uid", userId), zap.Error(err))
	}
}
