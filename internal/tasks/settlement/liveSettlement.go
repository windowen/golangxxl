package settlement

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"liveJob/pkg/constant"
	constsR "liveJob/pkg/constant/redis"
	"liveJob/pkg/context"
	"liveJob/pkg/db/cache"
	"liveJob/pkg/db/mysql"
	"liveJob/pkg/db/table"
	"liveJob/pkg/tools/cast"
	"liveJob/pkg/tools/utils"
	"liveJob/pkg/xxl"
	"liveJob/pkg/zlogger"
)

// DemoRebate 测试任务
func DemoRebate(cxt *context.Context, _ *xxl.RunReq) (msg string) {
	cxt.Trace = fmt.Sprintf("DemoRebate_%s", cxt.Trace)
	zlogger.Infof("【直播结算DemoRebate任务】%v 开始执行", time.Now())

	return "success"
}

// LiveSettlement 直播结算
func LiveSettlement(cxt *context.Context, _ *xxl.RunReq) (msg string) {
	cxt.Trace = fmt.Sprintf("LiveSettlement_%s", cxt.Trace)

	zlogger.Infof("【直播结算LiveSettlement任务】开始执行")

	start := time.Now()

	lastWeekStart, lastWeekEnd := utils.GetLastWeekTime()
	lastWeekIndex := utils.GetLastWeekIndex()

	err := mysql.LiveDB.Transaction(func(tx *gorm.DB) error {
		// 只有主播的
		insertQuery := fmt.Sprintf(`INSERT INTO live_room_income_settlement (owner_id, settlement_cycle, family_id, family_master_id, room_id, country_name, streamer_income, platform_income, family_income, settlement_time)
		SELECT 
			anchor_id as owner_id, 
			%d AS settlement_cycle, -- 使用上周的索引
			family_id, 
			family_master_id, 
			room_id, 
			country_name, 
			IFNULL(SUM(anchor_income),0)/10000 as streamer_income, 
			-IFNULL(SUM(platform_income),0)/10000 as platform_income, 
			-IFNULL(SUM(family_income),0)/10000 as family_income, 
			? as settlement_time 
		FROM 
			live_room_income_details
		WHERE 
			settlement_status = 0 AND family_master_id != anchor_id AND created_at BETWEEN ? AND ?
		GROUP BY 
			anchor_id, family_id, family_master_id, room_id, country_name`, lastWeekIndex)
		// 执行插入查询
		streamerResult := tx.WithContext(*cxt.Ctx).Exec(insertQuery, start, lastWeekStart, lastWeekEnd)
		if streamerResult.Error != nil {
			zlogger.Errorw("主播结算任务报错 LiveSettlement",
				zap.String("sql", insertQuery),
				zap.Time("startTime", lastWeekStart),
				zap.Time("endTime", lastWeekEnd),
				zap.Error(streamerResult.Error),
			)
			return streamerResult.Error
		}

		// 只有家族长的
		familyInsertQuery := fmt.Sprintf(`INSERT INTO live_room_income_settlement (owner_id, settlement_cycle, family_id, family_master_id, room_id, country_name, streamer_income, platform_income, family_income, settlement_time)
		SELECT 
			anchor_id as owner_id, 
			%d AS settlement_cycle, -- 使用上周的索引
			family_id, 
			family_master_id, 
			room_id, 
			country_name, 
			IFNULL(SUM(anchor_income),0)/10000 as streamer_income, 
			-IFNULL(SUM(platform_income),0)/10000 as platform_income, 
			-IFNULL(SUM(family_income),0)/10000 as family_income, 
			? as settlement_time 
		FROM 
			live_room_income_details
		WHERE 
			settlement_status = 0 AND family_master_id = anchor_id AND created_at BETWEEN ? AND ? 
		GROUP BY 
			family_master_id, family_id, room_id, country_name`, lastWeekIndex)
		// 执行插入查询
		if err := tx.WithContext(*cxt.Ctx).Exec(familyInsertQuery, start, lastWeekStart, lastWeekEnd).Error; err != nil {
			zlogger.Errorw("家族长主播结算任务报错 LiveSettlement",
				zap.String("sql", familyInsertQuery),
				zap.Time("startTime", lastWeekStart),
				zap.Time("endTime", lastWeekEnd),
				zap.Error(err),
			)
			return err
		}

		if streamerResult.RowsAffected > 0 {
			var familySettlements []table.LiveRoomIncomeSettlement
			// 查询出家族长的结算记录
			tx.WithContext(*cxt.Ctx).Table("live_room_income_settlement").
				Select("family_master_id, IFNULL(SUM(family_income,0)) AS family_income").
				Where("family_master_id != owner_id AND family_master_id != 0 AND settlement_cycle = ?", lastWeekIndex).
				Group("family_master_id").
				Scan(&familySettlements)
			for _, settlement := range familySettlements {
				if err := tx.WithContext(*cxt.Ctx).Model(&table.LiveRoomIncomeSettlement{}).
					Where("owner_id = ?", settlement.FamilyMasterId).
					Where("settlement_cycle = ?", lastWeekIndex).
					Update("family_income", gorm.Expr("family_income + ?", -settlement.FamilyIncome)).Error; err != nil {
					zlogger.Errorw("家族长主播结算任务报错 LiveSettlement",
						zap.Time("startTime", lastWeekStart),
						zap.Time("endTime", lastWeekEnd),
						zap.Error(err),
					)
					return err
				}
			}
		}

		if err := tx.WithContext(*cxt.Ctx).Model(&table.LiveRoomIncomeDetails{}).
			Where("created_at >= ? AND created_at <= ? AND settlement_status = 0", lastWeekStart, lastWeekEnd).
			Update("settlement_status", 1).Error; err != nil {
			zlogger.Errorw("更新结算标识任务报错 LiveSettlement",
				zap.Time("startTime", lastWeekStart),
				zap.Time("endTime", lastWeekEnd),
				zap.Error(err),
			)
			return err
		}
		return nil
	})

	if err != nil {
		zlogger.Errorf("直播结算LiveSettlement任务 error:%v", err)
	}

	end := time.Since(start)
	zlogger.Infof("【直播结算LiveSettlement任务】执行结束，耗时：%v", end)
	return "success"
}

// UserMountExpiredCheck 用户坐骑过期检查
func UserMountExpiredCheck(cxt *context.Context, _ *xxl.RunReq) (msg string) {
	cxt.Trace = fmt.Sprintf("UserMountExpiredCheck_%s", cxt.Trace)
	// zlogger.Infof("【坐骑过期检查UserMountExpiredCheck任务】%v 开始执行", time.Now())
	// start := time.Now()

	var userMounts []*table.UserMounts
	if err := mysql.LiveDB.WithContext(*cxt.Ctx).
		Model(&table.UserMounts{}).
		Where("status = ?", constant.StatusNormal).
		Where("expired_time <= ?", time.Now()).
		Find(&userMounts).Error; err != nil {
		zlogger.Errorf("UserMountExpiredCheck 查询用户坐骑 | err: %v", err)
	}

	// todo 考虑使用goroutine
	for _, mount := range userMounts {
		err := mysql.LiveDB.WithContext(*cxt.Ctx).Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&table.UserMounts{}).Where("id = ?", mount.Id).Updates(map[string]interface{}{"is_selected": constant.No, "status": constant.StatusDel}).Error
			if err != nil {
				return err
			}

			err = cache.RedisClient.HDel(*cxt.Ctx, fmt.Sprintf(constsR.UserMountCacheHash, mount.UserId), cast.ToString(mount.MountsId)).Err()
			if err != nil {
				return err
			}

			// 获取用户信息
			userCacheInfo, err := getUserCache(mount.UserId)
			if err != nil {
				return err
			}

			// 取消用户使用坐骑
			if userCacheInfo.MountsId == mount.MountsId {
				userCacheInfo.MountsId = constant.Zero

				err = setUserCache(mount.UserId, userCacheInfo)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			zlogger.Errorf("UserMountExpiredCheck |userId:%v,mountId:%v| err: %v", mount.UserId, mount.MountsId, err)
		}
	}

	// end := time.Since(start)
	// zlogger.Infof("【坐骑过期检查UserMountExpiredCheck任务】执行结束，耗时：%v", end)

	return "success"
}
