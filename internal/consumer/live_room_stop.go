package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/table"
	"queueJob/pkg/queue"
	"queueJob/pkg/zlogger"
)

var live = &liveRoomStop{}

type liveRoomStop struct{}

func (o *liveRoomStop) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		zlogger.Debugw("live room stop message", zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))

		liveRoom := &queue.LiveRoomStop{}
		if err := json.Unmarshal(msg.Body, liveRoom); err != nil {
			zlogger.Errorw("liveRoomStop::handleMessages, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
			continue
		}

		nowTime := time.Now()

		// 加锁防止重复结算
		lockSign := fmt.Sprintf("room_stop_%d_%d_%d", liveRoom.RoomId, liveRoom.AnchorId, liveRoom.SceneId)
		isLock, retFun := tryGetDistributedLock(lockSign, lockSign, 10000, 10000)
		if !isLock {
			zlogger.Errorw("room stop lock failed", zap.String("lock_sign", lockSign), zap.Int("uid", liveRoom.AnchorId))
			continue
		}

		err := mysql.LiveDB.Transaction(func(tx *gorm.DB) error {
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
					IFNULL(SUM(platform_income),0)/10000 as platform_income, 
					IFNULL(SUM(family_income),0)/10000 as family_income, 
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
			if streamerResult.RowsAffected == 0 {
				zlogger.Infow("liveRoomStop::handleMessages, 主播没有礼物数据要结算", zap.String("sql", insertQuery))
				return errors.New("no rows affected")
			}

			// 获取查询的结果
			result := &table.LiveRoomIncomeSettlement{}
			if err := tx.WithContext(ctx).Model(&table.LiveRoomIncomeSettlement{}).
				Select("owner_id, family_master_id, streamer_income, family_income").
				Where("scene_id = ?", liveRoom.SceneId).
				First(result).Error; err != nil {
				zlogger.Errorw("liveRoomStop::handleMessages, get streamer income error", zap.Error(err))
				return err
			}

			// 更新主播的结算收入
			if err := tx.Model(&table.SiteUserWallet{}).
				Where("user_id = ?", result.OwnerId).
				Updates(map[string]interface{}{
					"settlement_diamond": gorm.Expr("settlement_diamond + ?", result.StreamerIncome),
					// "balance":            gorm.Expr("balance + ?", decimal.NewFromInt(result.StreamerIncome).Div(decimal.NewFromInt32(100)).Truncate(2)),
				}).Error; err != nil {
				zlogger.Errorw("liveRoomStop::handleMessages, update site user wallet error", zap.Error(err))
				return err
			}

			// 更新家族长收入
			if result.FamilyMasterId > 0 {
				if err := tx.Model(&table.SiteUserWallet{}).
					Where("user_id = ?", result.FamilyMasterId).
					Updates(map[string]interface{}{
						"settlement_diamond": gorm.Expr("settlement_diamond + ?", result.FamilyIncome),
						// "balance":            gorm.Expr("balance + ?", decimal.NewFromInt(result.FamilyIncome).Div(decimal.NewFromInt32(100)).Truncate(2)),
					}).Error; err != nil {
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

		retFun()

		if err != nil {
			zlogger.Infow("live room stop message result", zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic),
				zap.String("body", string(msg.Body)), zap.Error(err))
			continue
		}
	}
	return consumer.ConsumeSuccess, nil
}
