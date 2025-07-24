package kafkaconsumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	constsR "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/db/table"
	"queueJob/pkg/queue"
	"queueJob/pkg/tools/utils"
	"queueJob/pkg/zlogger"
)

var cStats = &channelStats{}

type channelStats struct{}

func (o *channelStats) handleMessages(msg []byte) {
	ctx := context.Background()
	jsonData := &queue.ChannelStats{}
	if err := json.Unmarshal(msg, jsonData); err != nil {
		zlogger.Errorw("channelStats::handleMessages, unmarshal msg fail", zap.Error(err))
		return
	}

	channelManage := &table.ChannelManage{}
	if jsonData.ChannelCode != 0 {
		// 查询渠道配置
		if err := mysql.LiveDB.WithContext(ctx).Where("channel_code = ?", jsonData.ChannelCode).First(channelManage).Error; err != nil {
			zlogger.Errorw("channelStats::handleMessages, get channel data error",
				zap.Int("uid", jsonData.UserId), zap.Int("channel_code", jsonData.ChannelCode), zap.Error(err))
			return
		}

		if channelManage.Status == 0 {
			zlogger.Debugw("channelStats::handleMessages, channel is ban",
				zap.Int("uid", jsonData.UserId), zap.Int("channel_code", jsonData.ChannelCode))
			return
		}
	}

	switch {
	// 按注册量扣除
	case channelManage.DeductionType == 3 && jsonData.StatsType == queue.StatsUserRegistrations:
		register, err := o.registerDevice(jsonData.DeviceId)
		if err != nil {
			zlogger.Errorw("channelStats::handleMessages, cache device id error",
				zap.Int("uid", jsonData.UserId), zap.Any("stats", jsonData), zap.Error(err))
			return
		}

		if !register {
			zlogger.Infow("channelStats::handleMessages, device is register", zap.Any("stats", jsonData))
			return
		}
	case jsonData.StatsType == queue.DeviceActive:
		// 判断该设备今天是否统计过活跃
		key := fmt.Sprintf(constsR.DeviceKey, jsonData.DeviceId)
		if redis.HasKeyToday(key) {
			return
		}
	case jsonData.StatsType == queue.UserActive:
		if err := mysql.LiveDB.WithContext(ctx).
			Table("site_user").
			Where("id = ?", jsonData.UserId).
			Update("last_active_at", time.Now()).Error; err != nil {
			zlogger.Errorw("channelStats::handleMessages, get record error", zap.Int("uid", jsonData.UserId), zap.Error(err))
			return
		}
		return
	}

	// 加锁防止重复
	lockSign := fmt.Sprintf(redis.ChannelReportLock, channelManage.ChannelCode)
	isLock, retFun := redis.TryGetDistributedLock(lockSign, lockSign, 30000, 30000)
	if !isLock {
		zlogger.Errorw("lock failed", zap.String("lock_sign", lockSign), zap.Int("channelCode", channelManage.ChannelCode))
		return
	}

	defer retFun()

	nowTime := time.Now()

	today := nowTime.Format(time.DateOnly)

	// 查询是否已有该渠道的今日统计数据
	report := &table.ChannelReport{}
	err := mysql.LiveDB.WithContext(ctx).Where("channel_code = ? AND created_date = ?", jsonData.ChannelCode, today).First(report).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zlogger.Errorw("channelStats::handleMessages, get record error", zap.Int("uid", jsonData.UserId), zap.Error(err))
		return
	}

	// 如果不存在，创建新记录
	if errors.Is(err, gorm.ErrRecordNotFound) {
		report = &table.ChannelReport{
			ChannelCode:        jsonData.ChannelCode,
			ChannelPartnerId:   channelManage.ChannelPartnerId,
			ChannelPartnerName: channelManage.ChannelPartnerName,
			ChannelUrl:         channelManage.ChannelUrl,
			CreatedDate:        nowTime,
		}
	}

	err = mysql.LiveDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		switch jsonData.StatsType {
		case queue.StatsUserRegistrations:
			err = handleUserRegistration(ctx, tx, jsonData, report, channelManage, nowTime)
		case queue.StatsUserFirstRecharge:
			err = handleUserFirstRecharge(ctx, tx, jsonData, report, channelManage, nowTime)
		case queue.StatsUserRecharge:
			err = handleUserRecharge(ctx, tx, jsonData, report, channelManage)
		case queue.StatsWithdraw:
			report.WithdrawMemberCount++
			report.WithdrawSumMoney = report.WithdrawSumMoney.Add(jsonData.RechargeNum)
		case queue.StatsBet:
			report.ValidSumBet = report.ValidSumBet.Add(jsonData.BetNum)
		case queue.DeviceActive:
			err = handleDeviceActive(ctx, tx, jsonData, report, today, nowTime)
		case queue.StatsExchangeDiamond:
			report.ExchangeDiamondCount += jsonData.DiamondNum
			key := fmt.Sprintf(constsR.ExchangeDiamondStats, jsonData.UserId)
			if redis.HasKeyToday(key) {
				break
			}
			report.ExchangeDiamondTimes++
			redis.MarkKeyToday(key)
		case queue.StatsPlayGame:
			key := fmt.Sprintf(constsR.PlayGameStats, jsonData.UserId)
			if redis.HasKeyToday(key) {
				break
			}
			report.PlayGameTimes++
			redis.MarkKeyToday(key)
		default:
			zlogger.Errorw("unhandled default case")
		}

		// 最终更新 report 数据
		if report.Id > 0 {
			return tx.WithContext(ctx).Model(&report).Updates(report).Error
		}
		return tx.WithContext(ctx).Create(report).Error
	})

	if err != nil {
		zlogger.Errorw("channelStats::handleMessages, create report fail", zap.Int("uid", jsonData.UserId), zap.Error(err))
		return
	}
	zlogger.Infow("channelStats::handleMessages success", zap.Any("jsonData", jsonData))
}

// 尝试注册设备，如果是首次注册返回 true，否则 false
func (o *channelStats) registerDevice(deviceID string) (bool, error) {
	added, err := redis.SAdd(constsR.UserDeviceId, deviceID)
	if err != nil {
		return false, err
	}
	return added == 1, nil
}

func handleUserRegistration(ctx context.Context,
	tx *gorm.DB,
	jsonData *queue.ChannelStats,
	report *table.ChannelReport,
	channelManage *table.ChannelManage,
	nowTime time.Time,
) error {
	switch jsonData.Platform {
	case "ios":
		report.RegisterSum++
	case "android":
		report.AndroidRegisterSum++
	case "h5":
		report.H5RegisterSum++
	}

	registerCnt := report.RegisterSum + report.AndroidRegisterSum + report.H5RegisterSum
	report.RegisterCount = registerCnt
	report.CRegisterCount = registerCnt

	isDeduct := false
	if report.CRegisterCount > channelManage.Threshold && channelManage.DeductionType == 3 && channelManage.Per > 0 {
		report.CRegisterCount -= (report.CRegisterCount - channelManage.Threshold) / channelManage.Per
		isDeduct = true
	}

	var h5Active table.ChannelH5Active
	err := tx.WithContext(ctx).Where("device_id = ?", jsonData.DeviceId).First(&h5Active).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = tx.WithContext(ctx).Create(&table.ChannelH5Active{
			ChannelCode:  jsonData.ChannelCode,
			DeviceId:     jsonData.DeviceId,
			ActiveTime:   nowTime,
			RegisterTime: nowTime,
		}).Error
	}

	if err != nil {
		zlogger.Errorw("channelStats::handleUserRegistration, set h5 channel active record error",
			zap.Int("uid", jsonData.UserId), zap.Error(err))
	}

	if isDeduct {
		return nil
	}

	if channelManage.PixelId == "" {
		return nil
	}

	if err = SendFacebookEvent(channelManage.PixelId, channelManage.PixelToken, &FacebookEvent{
		EventName:      "Lead",
		EventTime:      time.Now().Unix(),
		ActionSource:   "website",
		EventSourceURL: channelManage.ChannelUrl,
		UserData: &UserData{
			Emails:          []string{jsonData.Email},
			ClientIP:        jsonData.Ip,
			ClientUserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		},
	}); err != nil {
		zlogger.Errorw("channelStats::handleUserRegistration, send event error",
			zap.Int("uid", jsonData.UserId), zap.Error(err))
	}
	return err
}

func updateChannelStatusAndRechargeStats(
	ctx context.Context,
	tx *gorm.DB,
	jsonData *queue.ChannelStats,
	report *table.ChannelReport,
	channelManage *table.ChannelManage,
	isFirstRecharge bool,
) error {
	var siteUser table.User
	if err := tx.WithContext(ctx).Where("id = ?", jsonData.UserId).First(&siteUser).Error; err != nil {
		return err
	}

	isDeduct := false
	if siteUser.ChannelStatus != 1 {
		report.ValidRechargeTimes++
		excess := report.ValidRechargeTimes - channelManage.Threshold
		if excess > 0 && channelManage.Per > 0 && excess%channelManage.Per == 0 {
			// 对于渠道商来说，充值状态为无效，充值数据不需要让渠道商看到
			if siteUser.ChannelStatus != 2 {
				if err := tx.WithContext(ctx).Model(&table.FinancePayRecord{}).Where("id = ?", jsonData.RecordId).Updates(map[string]interface{}{
					"channel_code":       channelManage.ChannelCode,
					"channel_partner_id": channelManage.ChannelPartnerId,
					"channel_status":     2,
				}).Error; err != nil {
					return err
				}

				isDeduct = true

				if err := tx.WithContext(ctx).
					Table("site_user").
					Where("id = ?", jsonData.UserId).
					Update("channel_status", 2).Error; err != nil {
					return err
				}
			}
		} else {
			if err := tx.WithContext(ctx).Model(&table.FinancePayRecord{}).Where("id = ?", jsonData.RecordId).Updates(map[string]interface{}{
				"channel_code":       channelManage.ChannelCode,
				"channel_partner_id": channelManage.ChannelPartnerId,
				"channel_status":     1,
			}).Error; err != nil {
				return err
			}

			// 对于渠道商来说，充值状态为有效，充值数据需要让渠道商看到
			if err := tx.WithContext(ctx).
				Table("site_user").
				Where("id = ?", jsonData.UserId).
				Update("channel_status", 1).Error; err != nil {
				return err
			}
			report.CRechargeSumMoney = report.CRechargeSumMoney.Add(jsonData.RechargeNum)
			report.CRechargeMemberTimes++
		}
	} else {
		if err := tx.WithContext(ctx).Model(&table.FinancePayRecord{}).Where("id = ?", jsonData.RecordId).Updates(map[string]interface{}{
			"channel_code":       channelManage.ChannelCode,
			"channel_partner_id": channelManage.ChannelPartnerId,
			"channel_status":     1,
		}).Error; err != nil {
			return err
		}
		// ChannelStatus 是1, 对于渠道商来说，充值有效，需要让渠道商看到
		report.CRechargeSumMoney = report.CRechargeSumMoney.Add(jsonData.RechargeNum)
		report.CRechargeMemberTimes++
	}

	if isDeduct {
		return nil
	}

	if channelManage.PixelId == "" {
		return nil
	}

	if isFirstRecharge {
		if err := SendFacebookEvent(channelManage.PixelId, channelManage.PixelToken, &FacebookEvent{
			EventName:      "Purchase",
			EventTime:      time.Now().Unix(),
			ActionSource:   "website",
			EventSourceURL: channelManage.ChannelUrl,
			UserData: &UserData{
				Emails:          []string{jsonData.Email},
				ClientIP:        jsonData.Ip,
				ClientUserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			},
			CustomData: &CustomData{
				Currency: "usd",
				Value:    jsonData.RechargeNum.InexactFloat64(),
				// ContentType: "product",
				// ContentIds:  []string{productId},
			},
		}); err != nil {
			zlogger.Errorw("channelStats::handleUserRegistration, send event error",
				zap.Int("uid", jsonData.UserId), zap.Error(err))
		}
	}

	return nil
}

func handleUserFirstRecharge(
	ctx context.Context,
	tx *gorm.DB,
	jsonData *queue.ChannelStats,
	report *table.ChannelReport,
	channelManage *table.ChannelManage,
	nowTime time.Time,
) error {
	report.FirstRechargeMemberCount++
	report.FirstRechargeMoney = report.FirstRechargeMoney.Add(jsonData.RechargeNum)
	report.RechargeSumMoney = report.RechargeSumMoney.Add(jsonData.RechargeNum)
	report.RechargeMemberTimes++

	key := fmt.Sprintf(constsR.UserRecharge, jsonData.UserId)
	if !redis.HasUserRechargedToday(key) {
		report.RechargeMemberCount++
		redis.MarkUserRechargedToday(key)
	}

	userCache, err := redis.GetUserCache(jsonData.UserId)
	if err == nil && utils.IsSameDay(userCache.RegisterTime, nowTime) {
		report.TodayRegisterRechargeCount++
	}

	if channelManage.DeductionType != 1 {
		return nil
	}

	return updateChannelStatusAndRechargeStats(ctx, tx, jsonData, report, channelManage, true)
}

func handleUserRecharge(
	ctx context.Context,
	tx *gorm.DB,
	jsonData *queue.ChannelStats,
	report *table.ChannelReport,
	channelManage *table.ChannelManage,
) error {
	report.RechargeSumMoney = report.RechargeSumMoney.Add(jsonData.RechargeNum)
	report.RechargeMemberTimes++

	key := fmt.Sprintf(constsR.UserRecharge, jsonData.UserId)
	if !redis.HasUserRechargedToday(key) {
		report.RechargeMemberCount++
		redis.MarkUserRechargedToday(key)
	}

	if channelManage.DeductionType != 1 {
		return nil
	}

	return updateChannelStatusAndRechargeStats(ctx, tx, jsonData, report, channelManage, false)
}

func handleDeviceActive(
	ctx context.Context,
	tx *gorm.DB,
	jsonData *queue.ChannelStats,
	report *table.ChannelReport,
	today string,
	nowTime time.Time,
) error {

	redis.MarkKeyToday(fmt.Sprintf(constsR.DeviceKey, jsonData.DeviceId))
	report.H5ActiveUsers++
	redis.HIncr(fmt.Sprintf(constsR.LiveStats, today), constsR.EventUserH5Active)

	if jsonData.IsNewDevice == 1 {
		report.H5NewDevices++
	}

	updateRetention := func(daysAgo int, rateColumn, numColumn string) error {
		// date := nowTime.AddDate(0, 0, -daysAgo).Format(time.DateOnly)
		date := nowTime.Format(time.DateOnly)

		if err := tx.WithContext(ctx).Model(&table.ChannelDailyUserStats{}).
			Where("report_date = ?", date).
			Updates(map[string]interface{}{rateColumn: gorm.Expr(fmt.Sprintf("%s + 1", rateColumn))}).Error; err != nil {
			return err
		}

		if err := tx.Model(&table.ChannelReport{}).
			Where("channel_code = ? AND created_date = ?", jsonData.ChannelCode, date).
			Updates(map[string]interface{}{numColumn: gorm.Expr(fmt.Sprintf("%s + 1", numColumn))}).Error; err != nil {
			return err
		}

		return nil
	}

	switch {
	case jsonData.Active1 == 1:
		if err := updateRetention(1, "h5_retention_rate", "h5_retention_num"); err != nil {
			return err
		}
	case jsonData.Active7 == 1:
		if err := updateRetention(1, "h5_retention_rate", "h5_retention_num"); err != nil {
			return err
		}
		if err := updateRetention(7, "h5_retention_rate_7", "h5_retention_num_7"); err != nil {
			return err
		}
	case jsonData.Active30 == 1:
		if err := updateRetention(1, "h5_retention_rate", "h5_retention_num"); err != nil {
			return err
		}
		if err := updateRetention(7, "h5_retention_rate_7", "h5_retention_num_7"); err != nil {
			return err
		}
		if err := updateRetention(30, "h5_retention_rate_30", "h5_retention_num_30"); err != nil {
			return err
		}
	}

	return nil
}
