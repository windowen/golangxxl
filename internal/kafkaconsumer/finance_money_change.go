package kafkaconsumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"

	redis2 "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/redisdb/redis"

	"queueJob/pkg/db/table"
	"queueJob/pkg/kafka"
	"queueJob/pkg/queue"
	"queueJob/pkg/zlogger"
)

var bChange = &balanceChange{}

type balanceChange struct{}

func (o *balanceChange) handleMessages(msg []byte) {
	ctx := context.Background()
	moneyChange := &queue.MoneyChange{}
	if err := json.Unmarshal(msg, moneyChange); err != nil {
		zlogger.Errorw("balanceChange::handleMessages, unmarshal msg fail", zap.Error(err))
		return
	}

	zlogger.Debugw("balanceChange::handleMessages", zap.Stringer("moneyChange", moneyChange))

	if !moneyChange.ReqAmount.Truncate(4).IsZero() || moneyChange.ChangeType != 4 {
		dbTable := &table.MoneyChange{
			UserId:          moneyChange.UserId,
			CountryCode:     moneyChange.CountryCode,
			CountryName:     moneyChange.CountryName,
			GameOrderNo:     moneyChange.GameOrderNo,
			TransNo:         moneyChange.TransNo,
			ChangeType:      moneyChange.ChangeType,
			ChangeAmount:    moneyChange.ChangeAmount,
			BeforeAmount:    moneyChange.BeforeAmount,
			AfterAmount:     moneyChange.AfterAmount,
			ReqAmount:       moneyChange.ReqAmount,
			ExchangeRate:    moneyChange.ExchangeRate,
			Currency:        moneyChange.Currency,
			Remark:          moneyChange.Remark,
			GameProvider:    moneyChange.GameProvider,
			GameType:        moneyChange.GameType,
			GameName:        moneyChange.GameName,
			TradeType:       moneyChange.TradeType,
			WagerCode:       moneyChange.WagerCode,
			BatchNum:        moneyChange.BatchNum,
			ActChangeAmount: moneyChange.ActChangeAmount,
			ActBeforeAmount: moneyChange.ActBeforeAmount,
			ActAfterAmount:  moneyChange.ActAfterAmount,
			CreatedAt:       moneyChange.CreatedAt,
			CreatedDay:      moneyChange.CreatedAt.Format(time.DateOnly),
		}

		if !moneyChange.AfterFlowAmount.IsZero() {
			dbTable.AfterFlowAmount = moneyChange.AfterFlowAmount
		}

		if !moneyChange.ActAfterFlowAmount.IsZero() {
			dbTable.ActAfterFlowAmount = moneyChange.ActAfterFlowAmount
		}

		if err := mysql.LiveDB.WithContext(ctx).Table(dbTable.TableName()).Create(dbTable).Error; err != nil {
			zlogger.Errorw("balanceChange::handleMessages, insert money change table fail",
				zap.Any("dbTable", dbTable), zap.Error(err))
			return
		}
		moneyChange.Id = dbTable.Id
	}

	if moneyChange.ChangeType != 4 {
		// 不是游戏下注账变，直接返回
		return
	}

	userId := moneyChange.UserId
	// 获取下注信息缓存
	wagerCache, err := redis.GetWagerCache(userId, moneyChange.WagerCode)
	if err != nil {
		zlogger.Errorw("balanceChange::handleMessages, get wager cache error",
			zap.Int("uid", userId),
			zap.String("wager code", moneyChange.WagerCode),
			zap.Error(err))
		return
	}

	if wagerCache.Exchange.IsZero() {
		zlogger.Errorw("balanceChange::handleMessages, exchange is zero",
			zap.Stringer("cache", wagerCache),
			zap.Int("uid", wagerCache.UserId))
		return
	}

	// 加锁同步操作用户钱包数据
	lockSign := fmt.Sprintf(redis2.UserPay, userId)
	isLock, retFun := redis.TryGetDistributedLock(lockSign, lockSign, 10000, 10000)
	if !isLock {
		zlogger.Errorw("lock failed", zap.String("lock_sign", lockSign), zap.Int("uid", userId))
		return
	}

	defer retFun()
	// 查询钱包
	dbWallet := &table.SiteUserWallet{}
	if err = mysql.LiveDB.WithContext(ctx).Where("user_id = ?", userId).First(dbWallet).Error; err != nil {
		zlogger.Errorw("balanceChange::handleMessages, get user wallet record failed",
			zap.Int("uid", userId), zap.Error(err))
		return
	}

	key := fmt.Sprintf(redis2.UserWalletCache, userId)
	walletItem, err := redis.HGetAll(key)
	if err != nil {
		zlogger.Errorw("balanceChange::handleMessages get redis cache error",
			zap.Int("uid", userId), zap.Error(err))
		return
	}

	balance, ok := walletItem["balance"]
	if !ok {
		zlogger.Errorw("balanceChange::handleMessages get balance cache error",
			zap.Int("uid", userId), zap.Error(err))
		return
	}

	curBalance, err := decimal.NewFromString(balance)
	if err != nil {
		zlogger.Errorw("balanceChange::handleMessages parse balance error",
			zap.Int("uid", userId), zap.String("balance", balance), zap.Error(err))
		return
	}

	curBalance = curBalance.Mul(decimal.NewFromFloat(0.0001)).Truncate(4)

	totalBet, ok := walletItem["totalBet"]
	if !ok {
		zlogger.Errorw("balanceChange::handleMessages get total bet cache error",
			zap.Int("uid", userId), zap.Error(err))
		return
	}

	curTotalBet, err := decimal.NewFromString(totalBet)
	if err != nil {
		zlogger.Errorw("balanceChange::handleMessages parse totalBet error",
			zap.Int("uid", userId), zap.String("totalBet", totalBet), zap.Error(err))
		return
	}

	curTotalBet = curTotalBet.Mul(decimal.NewFromFloat(0.0001)).Truncate(4)

	updates := map[string]interface{}{
		"balance":   curBalance,
		"total_bet": curTotalBet,
	}

	dbWallet.Balance = curBalance
	dbWallet.TotalBet = curTotalBet

	balanceFlow := dbWallet.FlowAmount
	if wagerCache.ChangeAmount.GreaterThan(decimal.Zero) {
		balanceFlow = dbWallet.FlowAmount.Sub(wagerCache.ChangeAmount)
		if moneyChange.TradeType == 1 {
			// 更新现金表流水额
			if balanceFlow.LessThan(decimal.Zero) {
				balanceFlow = decimal.Zero
			}
			updates["flow_amount"] = balanceFlow
			dbWallet.FlowAmount = balanceFlow
		}
	}

	err = mysql.LiveDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新钱包表数据
		dbResult := tx.Model(&table.SiteUserWallet{}).
			Where("user_id = ?", userId).Updates(updates)
		if dbResult.Error != nil {
			return dbResult.Error
		}

		if moneyChange.TradeType == 1 {
			model := &table.MoneyChange{UserId: userId}
			if innerErr := tx.Table(model.TableName()).
				Where("id = ?", moneyChange.Id).
				Update("after_flow_amount", balanceFlow).Error; innerErr != nil {
				return innerErr
			}
		}

		return nil
	})

	if err != nil {
		zlogger.Errorw("balanceChange::handleMessages, update user wallet info error",
			zap.Int("uid", userId), zap.Error(err))
		return
	}

	zlogger.Infow("balanceChange::handleMessages, update user wallet info",
		zap.Int("uid", userId), zap.Any("updates", updates))

	if moneyChange.TradeType == 1 {
		// 累积下注数据
		o.cacheWager(wagerCache, moneyChange.CountryCode)
	}

	// 判断是否触发预警
	if moneyChange.TradeType == 2 {
		// 结算账变
		o.checkMonitor(ctx, wagerCache)
	}

	// 不需要改活动钱包的数据，直接返回
	if wagerCache.ActChangeAmount.IsZero() {
		return
	}

	// 处理活动钱包相关
	switch moneyChange.TradeType {
	case 1: // 投注
		o.wagerMessage(ctx, moneyChange, wagerCache, dbWallet)
	case 2: // 结算
		o.settleMessage(ctx, moneyChange, wagerCache, dbWallet)
	case 3: // 回滚
		o.rollbackMessage(ctx, moneyChange, wagerCache)
	}
}

func (o *balanceChange) cacheWager(wagerCache *redis.WagerCacheInfo, countryCode string) {
	// 获取当前日期字符串
	now := time.Now()

	// 计算当天 24 点的时间
	midnight := time.Date(now.Year(), now.Month(), now.Day()+2, 1, 0, 0, 0, now.Location())
	ttl := time.Until(midnight)

	wager := wagerCache.WagerAmount.Div(wagerCache.Exchange).Truncate(4)
	// 累积投注 美元
	wagerInt := wager.Mul(decimal.NewFromInt(10000)).IntPart()
	currDate := now.Format(time.DateOnly)

	redisKey := fmt.Sprintf(redis2.UserGameWager, wagerCache.UserId, currDate)
	_, err := redis.IncrBy(redisKey, wagerInt)
	if err != nil {
		zlogger.Errorw("balanceChange::cacheWager, cache user wager error",
			zap.Int("uid", wagerCache.UserId), zap.Error(err))
		return
	}

	err = redis.SetExpireKey(redisKey, ttl)
	if err != nil {
		zlogger.Errorw("balanceChange::cacheWager, set expire error",
			zap.Int("uid", wagerCache.UserId), zap.Error(err))
		return
	}

	key := fmt.Sprintf("system_stats_wager_%d", wagerCache.UserId)
	if !redis.HasKeyToday(key) {
		redis.HIncr(fmt.Sprintf(redis2.SystemStatsGameKey, currDate), countryCode)
		redis.MarkKeyToday(key)
	}

	_, err = redis.HIncrBy(fmt.Sprintf(redis2.SystemStatsWagerKey, currDate), countryCode, wagerInt)
	if err != nil {
		zlogger.Errorw("balanceChange::cacheWager, hIncrBy fail", zap.Int("uid", wagerCache.UserId),
			zap.String("countryCode", countryCode), zap.Error(err))
		return
	}
}

// 缓存用户累计下注和获利数据，并返回每笔赌注的盈利，下注，当日累计盈利，当日累计下注
func (o *balanceChange) cacheProfit(wagerCache *redis.WagerCacheInfo, countryCode string) (profit, allProfit, allWager decimal.Decimal, err error) {
	profit = decimal.Zero
	allProfit = decimal.Zero
	allWager = decimal.Zero

	if wagerCache.Exchange.IsZero() {
		zlogger.Errorw("balanceChange::cacheProfit, exchange is zero", zap.Int("uid", wagerCache.UserId))
		return
	}

	// 累积结算，单位: 美元, 乘以 10000 转为 int64
	settlement := wagerCache.SettlementAmount.Div(wagerCache.Exchange).Mul(decimal.NewFromInt(10000)).IntPart()

	// 获取当前日期字符串
	now := time.Now()
	currDate := now.Format(time.DateOnly)

	_, err = redis.HIncrBy(fmt.Sprintf(redis2.SystemStatsSettleKey, currDate), countryCode, settlement)
	if err != nil {
		zlogger.Errorw("balanceChange::cacheProfit, hIncrBy fail", zap.Int("uid", wagerCache.UserId), zap.Error(err))
		return
	}

	redisKey := fmt.Sprintf(redis2.UserGameProfit, wagerCache.UserId, currDate)
	settlementCache, err := redis.IncrBy(redisKey, settlement)
	if err != nil {
		zlogger.Errorw("balanceChange::cacheProfit, cache user profit error",
			zap.Int("uid", wagerCache.UserId), zap.Error(err))
		return
	}

	allProfit = decimal.NewFromInt(settlementCache).Mul(decimal.NewFromFloat(0.0001))

	// 计算过期的时间
	midnight := time.Date(now.Year(), now.Month(), now.Day()+2, 1, 0, 0, 0, now.Location())
	ttl := time.Until(midnight)
	err = redis.SetExpireKey(redisKey, ttl)
	if err != nil {
		zlogger.Errorw("balanceChange::cacheProfit, set expire error",
			zap.Int("uid", wagerCache.UserId), zap.Error(err))
		return
	}

	betString, err := redis.Get(fmt.Sprintf(redis2.UserGameWager, wagerCache.UserId, now.Format(time.DateOnly)))
	if err != nil && !errors.Is(err, redis.Nil) {
		zlogger.Errorw("balanceChange::cacheProfit, cache user wager error",
			zap.Int("uid", wagerCache.UserId), zap.Error(err))
		return
	}

	betCache, err := decimal.NewFromString(betString)
	if err != nil {
		zlogger.Errorw("balanceChange::cacheProfit, cache user wager error",
			zap.Int("uid", wagerCache.UserId), zap.Error(err))
		return
	}

	// 累积投注 美元
	allWager = betCache.Mul(decimal.NewFromFloat(0.0001))

	// 玩家中奖
	profit = wagerCache.SettlementAmount
	if profit.LessThanOrEqual(decimal.Zero) {
		return
	}

	profit = profit.Div(wagerCache.Exchange).Truncate(4)

	return
}

// 判断是否触发预警
func (o *balanceChange) checkMonitor(ctx context.Context, wagerCache *redis.WagerCacheInfo) {
	userCache, err := redis.GetUserCache(wagerCache.UserId)
	if err != nil {
		zlogger.Errorw("balanceChange::checkMonitor, get user cache error", zap.Int("uid", wagerCache.UserId), zap.Error(err))
		return
	}

	profit, allProfit, allWager, err := o.cacheProfit(wagerCache, userCache.CountryCode)
	if err != nil {
		zlogger.Errorw("balanceChange::checkMonitor, cache wager error", zap.Stringer("cache", wagerCache), zap.Error(err))
		return
	}

	if profit.LessThanOrEqual(decimal.Zero) {
		zlogger.Debugw("balanceChange::checkMonitor, profit is zero", zap.Stringer("cache", wagerCache))
		return
	}

	wager := wagerCache.WagerAmount.Div(wagerCache.Exchange)

	monitorConfig := redis.GetProfitMonitor(ctx)
	if monitorConfig == nil {
		zlogger.Errorw("balanceChange::checkMonitor, get profit monitor error")
		return
	}

	nowTime := time.Now()

	if profit.GreaterThanOrEqual(monitorConfig.LargePrize) {
		// 大额中奖
		mysql.LiveDB.WithContext(ctx).Create(&table.GameProfitMonitorUser{
			UserId:      userCache.Id,
			UserAccount: userCache.Email,
			UserTag: fmt.Sprintf("%d,%d,%d,%d",
				userCache.AgentRebateStatus, userCache.InviteRebateStatus, userCache.BetStatus, userCache.DrawStatus),
			LoginIp:     userCache.LoginIp,
			RegSource:   userCache.ChannelName,
			MonitorType: "大额中奖",
			ActualValue: profit,
			BetOrderNo:  wagerCache.WagerCode,
			Status:      0,
			CreatedDate: nowTime,
			CreatedAt:   nowTime,
			UpdatedAt:   nowTime,
		})

		if monitorConfig.RiskWarning == 1 {
			if incr, err := redis.Incr(redis2.MonitorNotify); err == nil {
				zlogger.Infow("balanceChange::checkMonitor, monitor notify", zap.Int64("incr", incr))
			}
		}
	}

	if monitorConfig.HighMultiplierJackpot == 1 && !wager.IsZero() {
		// 检查高倍爆奖
		if profit.GreaterThanOrEqual(monitorConfig.HighMultiplierJackpotMoney) &&
			profit.Div(wager).GreaterThanOrEqual(decimal.NewFromInt(int64(monitorConfig.HighMultiplierJackpotMultiple))) {
			mysql.LiveDB.WithContext(ctx).Create(&table.GameProfitMonitorUser{
				UserId:      userCache.Id,
				UserAccount: userCache.Email,
				UserTag: fmt.Sprintf("%d,%d,%d,%d",
					userCache.AgentRebateStatus, userCache.InviteRebateStatus, userCache.BetStatus, userCache.DrawStatus),
				LoginIp:     userCache.LoginIp,
				RegSource:   userCache.ChannelName,
				MonitorType: "高倍爆奖",
				ActualValue: profit,
				BetOrderNo:  wagerCache.WagerCode,
				Status:      0,
				CreatedDate: nowTime,
				CreatedAt:   nowTime,
				UpdatedAt:   nowTime,
			})

			if monitorConfig.RiskWarning == 1 {
				if incr, err := redis.Incr(redis2.MonitorNotify); err == nil {
					zlogger.Infow("balanceChange::checkMonitor, monitor notify", zap.Int64("incr", incr))
				}
			}
		}
	}

	if monitorConfig.ProfitMargin == 1 && !allWager.IsZero() {
		// 检查当日会员获利比
		if allProfit.GreaterThanOrEqual(monitorConfig.DailyProfitMarginQuota) &&
			allProfit.Div(allWager).GreaterThanOrEqual(decimal.NewFromInt(int64(monitorConfig.DailyProfitMargin))) {
			mysql.LiveDB.WithContext(ctx).Create(&table.GameProfitMonitorUser{
				UserId:      userCache.Id,
				UserAccount: userCache.Email,
				UserTag: fmt.Sprintf("%d,%d,%d,%d",
					userCache.AgentRebateStatus, userCache.InviteRebateStatus, userCache.BetStatus, userCache.DrawStatus),
				LoginIp:     userCache.LoginIp,
				RegSource:   userCache.ChannelName,
				MonitorType: "会员当天获利比",
				ActualValue: profit,
				BetOrderNo:  wagerCache.WagerCode,
				Status:      0,
				CreatedDate: nowTime,
				CreatedAt:   nowTime,
				UpdatedAt:   nowTime,
			})

			if monitorConfig.RiskWarning == 1 {
				if incr, err := redis.Incr(redis2.MonitorNotify); err == nil {
					zlogger.Infow("balanceChange::checkMonitor, monitor notify", zap.Int64("incr", incr))
				}
			}
		}
	}
}

func (o *balanceChange) settleMessage(
	ctx context.Context,
	moneyChange *queue.MoneyChange,
	wagerCache *redis.WagerCacheInfo,
	dbWallet *table.SiteUserWallet,
) {
	actKey := fmt.Sprintf(moneyChange.WalletKey, moneyChange.UserId)
	field := wagerCache.BatchNum

	actBalance := o.getAndScaleActBalance(actKey, field)
	zlogger.Infow("settleMessage get scale act balance",
		zap.Int("uid", moneyChange.UserId),
		zap.String("actKey", actKey),
		zap.String("field", field),
		zap.Stringer("actBalance", actBalance))

	switch moneyChange.WalletKey {
	case redis2.UserBonusWalletCache:
		o.settleBonusWallet(ctx, moneyChange, actBalance)
	case redis2.UserActWalletCache:
		o.settleActWallet(ctx, moneyChange, actBalance, wagerCache.BatchNum, actKey, dbWallet)
	}
}

func (o *balanceChange) getAndScaleActBalance(actKey, field string) decimal.Decimal {
	s, err := redis.HGet(actKey, field)
	if err != nil {
		return decimal.Zero
	}
	v, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return v.Mul(decimal.NewFromFloat(0.0001)).Truncate(4)
}

func (o *balanceChange) settleBonusWallet(
	ctx context.Context,
	m *queue.MoneyChange,
	actBalance decimal.Decimal,
) {
	dbRec := &table.ActBonusWallet{}
	if err := mysql.LiveDB.WithContext(ctx).Where("user_id = ?", m.UserId).First(dbRec).Error; err != nil {
		zlogger.Errorw("balanceChange::settleBonusWallet, get ActBonusWallet error",
			zap.Int("uid", m.UserId),
			zap.Error(err),
		)
		return
	}

	updates := o.prepareBalanceUpdate(actBalance)
	_ = o.updateMoneyChangeActFlow(ctx, m.Id, m.UserId, dbRec.FlowAmount)

	if (dbRec.FlowAmount.IsZero() && dbRec.ReachUsers == 0) ||
		actBalance.GreaterThanOrEqual(dbRec.ReachAmount) {
		kafka.PublicKey(strconv.Itoa(m.UserId), kafka.CheckRegisterBonus, &queue.CheckRegisterBonus{UserId: m.UserId})
		return
	}

	mysql.LiveDB.WithContext(ctx).Model(dbRec).Where("user_id = ?", m.UserId).Updates(updates)
	return
}

func (o *balanceChange) settleActWallet(
	ctx context.Context,
	m *queue.MoneyChange,
	actBalance decimal.Decimal,
	batchNum, actKey string,
	dbWallet *table.SiteUserWallet,
) {
	dbRec := &table.ActUserWallet{}
	if err := mysql.LiveDB.WithContext(ctx).
		Where("user_id = ? AND game_bonus_type_config_batch_num = ?", m.UserId, batchNum).
		First(dbRec).Error; err != nil {
		return
	}

	updates := o.prepareBalanceUpdate(actBalance)
	if !dbRec.FlowAmount.IsZero() {
		mysql.LiveDB.WithContext(ctx).
			Model(dbRec).
			Where("user_id = ? AND game_bonus_type_config_batch_num = ?", m.UserId, batchNum).
			Updates(updates)

		_ = o.updateMoneyChangeActFlow(ctx, m.Id, m.UserId, dbRec.FlowAmount)
		return
	}

	// 流水已完成，触发转余额
	_ = redis.HDel(actKey, batchNum)
	updates["balance"] = decimal.Zero
	_ = o.processCompletion(ctx, m, actBalance, dbWallet, updates, 6, true)
}

func (o *balanceChange) prepareBalanceUpdate(actBalance decimal.Decimal) map[string]interface{} {
	if actBalance.LessThanOrEqual(decimal.Zero) {
		return map[string]interface{}{
			"balance": decimal.Zero,
		}
	}
	return map[string]interface{}{
		"balance": actBalance,
	}
}

func (o *balanceChange) processCompletion(
	ctx context.Context,
	m *queue.MoneyChange,
	actBalance decimal.Decimal,
	dbWallet *table.SiteUserWallet,
	updates map[string]interface{},
	changeType int,
	swapFromZeroFlow bool,
) error {
	return mysql.LiveDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? ", m.UserId).Model(&table.ActUserWallet{}).Updates(updates).Error; err != nil {
			return err
		}

		rec := o.buildMoneyChangeRecord(m, actBalance, dbWallet, changeType)
		if swapFromZeroFlow {
			rec.ChangeAmount = m.ChangeAmount.Add(m.ActChangeAmount)
			rec.AfterAmount = m.AfterAmount.Add(m.ActChangeAmount)
			rec.ActChangeAmount = decimal.Zero
			rec.ActBeforeAmount = decimal.Zero
		}
		if err := tx.Table(rec.TableName()).Create(rec).Error; err != nil {
			return err
		}

		key := fmt.Sprintf(redis2.UserWalletCache, m.UserId)
		incr, err := redis.HIncrBy(key, "balance", actBalance.Mul(decimal.NewFromInt(10000)).IntPart())
		if err != nil {
			return err
		}
		curBal := decimal.NewFromInt(incr).Mul(decimal.NewFromFloat(0.0001)).Truncate(4)
		tx.Model(&table.SiteUserWallet{}).Where("user_id = ?", m.UserId).Update("balance", curBal)
		return nil
	})
}

// 下注逻辑
func (o *balanceChange) wagerMessage(
	ctx context.Context,
	moneyChange *queue.MoneyChange,
	wagerCache *redis.WagerCacheInfo,
	dbWallet *table.SiteUserWallet,
) {
	actKey := fmt.Sprintf(moneyChange.WalletKey, moneyChange.UserId)
	field := wagerCache.BatchNum

	actBalance, err := o.getActBalanceFromCache(actKey, field, moneyChange.UserId)
	if err != nil {
		zlogger.Errorw("wagerMessage get act wallet cache error",
			zap.Int("uid", moneyChange.UserId), zap.Error(err))
		return
	}
	actBalance = actBalance.Mul(decimal.NewFromFloat(0.0001)).Truncate(4)

	switch wagerCache.WalletKey {
	case redis2.UserActWalletCache:
		o.handleActWallet(ctx, moneyChange, wagerCache, actBalance, dbWallet)
	case redis2.UserBonusWalletCache:
		o.handleBonusWallet(ctx, moneyChange, actBalance)
	}
}

func (o *balanceChange) handleActWallet(
	ctx context.Context,
	moneyChange *queue.MoneyChange,
	wagerCache *redis.WagerCacheInfo,
	actBalance decimal.Decimal,
	dbWallet *table.SiteUserWallet,
) {
	dbRecord := &table.ActUserWallet{}
	err := mysql.LiveDB.WithContext(ctx).
		Where("user_id = ? AND game_bonus_type_config_batch_num = ?", moneyChange.UserId, wagerCache.BatchNum).
		First(dbRecord).Error
	if err != nil {
		return
	}

	flowAmount := dbRecord.FlowAmount.Add(moneyChange.ActChangeAmount)
	if flowAmount.LessThan(decimal.Zero) {
		flowAmount = decimal.Zero
	}

	if flowAmount.IsZero() && actBalance.GreaterThan(decimal.Zero) {
		_ = o.zeroMoneyChangeActFlow(ctx, moneyChange.Id, moneyChange.UserId)

		if err = o.transferActToBalance(ctx, moneyChange, wagerCache, actBalance, dbWallet, 6); err != nil {
			zlogger.Errorw("handleActWallet transfer error", zap.Int("uid", moneyChange.UserId), zap.Error(err))
		}
		return
	}

	o.updateActWalletFlow(ctx, moneyChange, wagerCache.BatchNum, actBalance, flowAmount)
}

func (o *balanceChange) handleBonusWallet(
	ctx context.Context,
	moneyChange *queue.MoneyChange,
	actBalance decimal.Decimal,
) {
	dbRecord := &table.ActBonusWallet{}
	err := mysql.LiveDB.WithContext(ctx).
		Where("user_id = ?", moneyChange.UserId).
		First(dbRecord).Error
	if err != nil {
		return
	}

	flowAmount := dbRecord.FlowAmount.Add(moneyChange.ActChangeAmount)
	if flowAmount.LessThan(decimal.Zero) {
		flowAmount = decimal.Zero
	}

	if flowAmount.IsZero() && actBalance.GreaterThan(decimal.Zero) && dbRecord.ReachUsers == 0 {
		kafka.PublicKey(strconv.Itoa(moneyChange.UserId), kafka.CheckRegisterBonus,
			&queue.CheckRegisterBonus{UserId: moneyChange.UserId})
		return
	}

	o.updateBonusWalletFlow(ctx, moneyChange, actBalance, flowAmount)
}

func (o *balanceChange) updateActWalletFlow(
	ctx context.Context,
	moneyChange *queue.MoneyChange,
	batchNum string,
	actBalance, flowAmount decimal.Decimal,
) {
	err := mysql.LiveDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新活动钱包表
		if err := tx.Model(&table.ActUserWallet{}).
			Where("user_id = ? AND game_bonus_type_config_batch_num = ?", moneyChange.UserId, batchNum).
			Updates(map[string]interface{}{
				"balance":     actBalance,
				"flow_amount": flowAmount,
			}).Error; err != nil {
			zlogger.Errorw("updateActWalletFlow, update act wallet error",
				zap.Int("uid", moneyChange.UserId), zap.Error(err))
			return err
		}

		// 更新 money_change 中的 act_after_flow_amount 字段
		if err := tx.Table((&table.MoneyChange{UserId: moneyChange.UserId}).TableName()).
			Where("id = ?", moneyChange.Id).
			Update("act_after_flow_amount", flowAmount).Error; err != nil {
			zlogger.Errorw("updateActWalletFlow, update money_change error",
				zap.Int("uid", moneyChange.UserId), zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		zlogger.Errorw("updateActWalletFlow transaction failed",
			zap.Int("uid", moneyChange.UserId), zap.Error(err))
	}
}

func (o *balanceChange) updateBonusWalletFlow(ctx context.Context, moneyChange *queue.MoneyChange, actBalance, flowAmount decimal.Decimal) {
	err := mysql.LiveDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新奖金钱包表
		if err := tx.Model(&table.ActBonusWallet{}).
			Where("user_id = ?", moneyChange.UserId).
			Updates(map[string]interface{}{
				"balance":     actBalance,
				"flow_amount": flowAmount,
			}).Error; err != nil {
			zlogger.Errorw("updateBonusWalletFlow, update bonus wallet error",
				zap.Int("uid", moneyChange.UserId), zap.Error(err))
			return err
		}

		// 更新 money_change 表中的 act_after_flow_amount 字段
		if err := tx.Table((&table.MoneyChange{UserId: moneyChange.UserId}).TableName()).
			Where("id = ?", moneyChange.Id).
			Update("act_after_flow_amount", flowAmount).Error; err != nil {
			zlogger.Errorw("updateBonusWalletFlow, update money_change error",
				zap.Int("uid", moneyChange.UserId), zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		zlogger.Errorw("updateBonusWalletFlow transaction failed",
			zap.Int("uid", moneyChange.UserId), zap.Error(err))
	}
}

func (o *balanceChange) transferActToBalance(
	ctx context.Context,
	moneyChange *queue.MoneyChange,
	wagerCache *redis.WagerCacheInfo,
	actBalance decimal.Decimal,
	dbWallet *table.SiteUserWallet,
	changeType int,
) error {
	return o.transferToMainBalance(ctx, moneyChange, actBalance, dbWallet, changeType, func(tx *gorm.DB) error {
		return tx.Model(&table.ActUserWallet{}).
			Where("user_id = ? AND game_bonus_type_config_batch_num = ?", moneyChange.UserId, wagerCache.BatchNum).
			Updates(map[string]interface{}{"balance": decimal.Zero, "flow_amount": decimal.Zero}).Error
	})
}

func (o *balanceChange) transferBonusToBalance(
	ctx context.Context,
	moneyChange *queue.MoneyChange,
	actBalance decimal.Decimal,
	dbWallet *table.SiteUserWallet,
	changeType int,
) error {
	return o.transferToMainBalance(ctx, moneyChange, actBalance, dbWallet, changeType, func(tx *gorm.DB) error {
		return tx.Model(&table.ActBonusWallet{}).
			Where("user_id = ?", moneyChange.UserId).
			Updates(map[string]interface{}{"balance": decimal.Zero, "flow_amount": decimal.Zero}).Error
	})
}

func (o *balanceChange) transferToMainBalance(
	ctx context.Context,
	moneyChange *queue.MoneyChange,
	actBalance decimal.Decimal,
	dbWallet *table.SiteUserWallet,
	changeType int,
	updateActWalletFn func(tx *gorm.DB) error,
) error {
	return mysql.LiveDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := updateActWalletFn(tx); err != nil {
			return err
		}

		if err := redis.HDel(fmt.Sprintf(moneyChange.WalletKey, moneyChange.UserId), moneyChange.BatchNum); err != nil {
			return err
		}

		key := fmt.Sprintf(redis2.UserWalletCache, moneyChange.UserId)
		incrResult, err := redis.HIncrBy(key, "balance", actBalance.Mul(decimal.NewFromInt(10000)).IntPart())
		if err != nil {
			return err
		}
		curBalance := decimal.NewFromInt(incrResult).Mul(decimal.NewFromFloat(0.0001)).Truncate(4)

		if err = tx.Model(&table.SiteUserWallet{}).
			Where("user_id = ?", moneyChange.UserId).
			Update("balance", curBalance).Error; err != nil {
			return err
		}

		dbTable := o.buildMoneyChangeRecord(moneyChange, actBalance, dbWallet, changeType)
		return tx.Table(dbTable.TableName()).Create(dbTable).Error
	})
}

func (o *balanceChange) buildMoneyChangeRecord(
	m *queue.MoneyChange,
	actBalance decimal.Decimal,
	dbWallet *table.SiteUserWallet,
	changeType int,
) *table.MoneyChange {
	return &table.MoneyChange{
		UserId:             m.UserId,
		CountryCode:        m.CountryCode,
		CountryName:        m.CountryName,
		GameOrderNo:        m.GameOrderNo,
		TransNo:            m.TransNo,
		ChangeType:         changeType,
		ChangeAmount:       actBalance,
		BeforeAmount:       dbWallet.Balance,
		AfterAmount:        dbWallet.Balance.Add(actBalance),
		AfterFlowAmount:    dbWallet.FlowAmount,
		ReqAmount:          m.ReqAmount,
		ExchangeRate:       m.ExchangeRate,
		Currency:           m.Currency,
		Remark:             m.Remark,
		GameProvider:       m.GameProvider,
		GameType:           m.GameType,
		GameName:           m.GameName,
		TradeType:          m.TradeType,
		WagerCode:          m.WagerCode,
		BatchNum:           m.BatchNum,
		ActChangeAmount:    actBalance.Neg(),
		ActBeforeAmount:    actBalance,
		ActAfterAmount:     decimal.Zero,
		ActAfterFlowAmount: decimal.Zero,
		CreatedAt:          m.CreatedAt,
		CreatedDay:         m.CreatedAt.Format(time.DateOnly),
	}
}

func (o *balanceChange) zeroMoneyChangeActFlow(ctx context.Context, id, userId int) error {
	return mysql.LiveDB.WithContext(ctx).
		Table((&table.MoneyChange{UserId: userId}).TableName()).
		Where("id = ?", id).
		Update("act_after_flow_amount", decimal.Zero).Error
}

func (o *balanceChange) updateMoneyChangeActFlow(ctx context.Context, id, userId int, flowAmount decimal.Decimal) error {
	return mysql.LiveDB.WithContext(ctx).
		Table((&table.MoneyChange{UserId: userId}).TableName()).
		Where("id = ?", id).
		Update("act_after_flow_amount", flowAmount).Error
}

func (o *balanceChange) getActBalanceFromCache(actKey, field string, userId int) (decimal.Decimal, error) {
	actBalanceStr, err := redis.HGet(actKey, field)
	if err != nil {
		zlogger.Errorw("getActBalanceFromCache error", zap.Int("uid", userId), zap.Error(err))
		return decimal.Zero, err
	}

	actBalance, err := decimal.NewFromString(actBalanceStr)
	if err != nil {
		zlogger.Errorw("parse act balance error", zap.Int("uid", userId), zap.Error(err))
		return decimal.Zero, err
	}
	return actBalance, nil
}

// 回滚逻辑, 只更新剩余打码量数据
func (o *balanceChange) rollbackMessage(ctx context.Context, moneyChange *queue.MoneyChange, wagerCache *redis.WagerCacheInfo) {
	actFlowAmount := decimal.Zero

	switch moneyChange.WalletKey {
	case redis2.UserActWalletCache:
		actWallet := &table.ActUserWallet{}
		if innerErr := mysql.LiveDB.WithContext(ctx).
			Where("user_id = ? AND game_bonus_type_config_batch_num = ?", moneyChange.UserId, wagerCache.BatchNum).
			First(actWallet).Error; innerErr == nil {
			actFlowAmount = actWallet.FlowAmount
		}
	case redis2.UserBonusWalletCache:
		actWallet := &table.ActBonusWallet{}
		if innerErr := mysql.LiveDB.WithContext(ctx).
			Where("user_id = ?", moneyChange.UserId).
			First(actWallet).Error; innerErr == nil {
			actFlowAmount = actWallet.FlowAmount
		}
	}

	// 查询钱包
	dbWallet := &table.SiteUserWallet{}
	if innerErr := mysql.LiveDB.WithContext(ctx).
		Where("user_id = ?", moneyChange.UserId).
		First(dbWallet).Error; innerErr != nil {
		return
	}

	err := mysql.LiveDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := &table.MoneyChange{UserId: moneyChange.UserId}
		if innerErr := tx.Table(model.TableName()).
			Where("id = ?", moneyChange.Id).
			Updates(map[string]interface{}{
				"after_flow_amount":     dbWallet.FlowAmount,
				"act_after_flow_amount": actFlowAmount,
			}).Error; innerErr != nil {
			return innerErr
		}

		return nil
	})

	if err != nil {
		zlogger.Errorw("balanceChange::rollbackMessage, update user wallet info error",
			zap.Int("uid", moneyChange.UserId), zap.Error(err))
		return
	}
}
