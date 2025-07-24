package consumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"liveJob/pkg/db/table"
	"liveJob/pkg/queue"
	"liveJob/pkg/zlogger"
)

func Test_liveRoomStop_handleMessages(t *testing.T) {
	nowTime := time.Now()
	liveRoom := &queue.LiveRoomStop{
		RoomId:   1009528,
		AnchorId: 132,
		SceneId:  250,
	}

	var (
		LiveDB *gorm.DB
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		"liveuser",
		"7hr56dv3mysql+4@8abi_^W",
		"18.167.121.240:3306",
		"live_m5")

	db, err := gorm.Open(
		mysql.Open(dsn), &gorm.Config{
			SkipDefaultTransaction:                   true, // 禁用默认事务
			DisableForeignKeyConstraintWhenMigrating: true, // 外键
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					Colorful: true,
					LogLevel: logger.Info,
				}),
		},
	)
	if err != nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(60))
	sqlDB.SetMaxOpenConns(1000)
	sqlDB.SetMaxIdleConns(100)

	LiveDB = db

	err = LiveDB.Transaction(func(tx *gorm.DB) error {
		// 查询这场直播的总收入，并插入到live_room_income_settlement表中
		insertQuery := fmt.Sprintf(`INSERT INTO live_room_income_settlement
				(owner_id, family_id, family_master_id, room_id, scene_id, country_name, streamer_income, platform_income, family_income, settlement_time)
				SELECT 
					anchor_id as owner_id, 
					family_id, 
					family_master_id, 
					room_id, 
					scene_id,
					country_name, 
					IFNULL(SUM(anchor_income),0)/10000 as streamer_income, 
					-IFNULL(SUM(platform_income),0)/10000 as platform_income, 
					-IFNULL(SUM(family_income),0)/10000 as family_income, 
					? as settlement_time 
				FROM 
					live_room_income_details
				WHERE 
					anchor_id = ? AND settlement_status = 0 AND scene_id = ?
				GROUP BY 
					anchor_id, family_id, family_master_id, room_id, scene_id, country_name`)
		// 执行插入查询
		streamerResult := tx.Exec(insertQuery, nowTime, liveRoom.AnchorId, liveRoom.SceneId)
		if streamerResult.Error != nil {
			zlogger.Errorw("liveRoomStop::handleMessages, 主播结算任务报错", zap.String("sql", insertQuery), zap.Error(streamerResult.Error))
			return streamerResult.Error
		}

		// 获取最后插入的自增ID
		var lastInsertID int
		if err := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&lastInsertID).Error; err != nil {
			zlogger.Errorw("liveRoomStop::handleMessages, get last_insert_id error", zap.Error(err))
			return err
		}

		// 获取查询的结果
		result := &table.LiveRoomIncomeSettlement{}
		if err := tx.WithContext(context.Background()).Model(&table.LiveRoomIncomeSettlement{}).
			Select("streamer_income, family_income").
			Where("id = ?", lastInsertID).
			First(result).Error; err != nil {
			zlogger.Errorw("liveRoomStop::handleMessages, get streamer income error", zap.Error(err))
			return err
		}

		// 更新主播的结算收入
		if err := tx.Model(&table.SiteUserWallet{}).
			Where("user_id = ?", result.OwnerId).
			Update("settlement_diamond", gorm.Expr("settlement_diamond + ?", result.StreamerIncome)).Error; err != nil {
			zlogger.Errorw("liveRoomStop::handleMessages, update site user wallet error", zap.Error(err))
			return err
		}

		// 更新家族长收入
		if result.FamilyMasterId > 0 {
			if err := tx.Model(&table.SiteUserWallet{}).
				Where("user_id = ?", result.FamilyMasterId).
				Update("settlement_diamond", gorm.Expr("settlement_diamond + ?", result.FamilyIncome)).Error; err != nil {
				zlogger.Errorw("liveRoomStop::handleMessages, update site family user wallet error", zap.Error(err))
				return err
			}
		}

		// 更新结算标识为已经结算
		if err := tx.Model(&table.LiveRoomIncomeDetails{}).
			Where("anchor_id = ? AND scene_id = ? AND settlement_status = 0", liveRoom.AnchorId, liveRoom.SceneId).
			Update("settlement_status", 1).Error; err != nil {
			zlogger.Errorw("liveRoomStop::handleMessages, 更新结算标识任务报错", zap.Error(err))
			return err
		}
		return nil
	})

	if err != nil {
		return
	}
}
