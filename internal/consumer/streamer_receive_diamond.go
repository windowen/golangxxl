package consumer

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	redis3 "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	redis2 "liveJob/pkg/constant/redis"
	"liveJob/pkg/db/mysql"
	"liveJob/pkg/db/redisdb/redis"
	"liveJob/pkg/db/table"
	"liveJob/pkg/queue"
	"liveJob/pkg/zlogger"
)

var (
	receiveDiamond = &streamerReceiveDiamond{}
	LevelStreamer  = `[10,30,70,190,385,640,1040,1665,2715,4290,6690,9840,14840,20840,27340,34840,42000,49500,57500,66000,75000,84500,94500,105000,115500,126000,138000,150500,163000,176000,189000,203000,218000,234000,251100,268000,286000,305000,325000,346000,368000,391000,415000,440000,466000,493000,521000,550000,580000,611000]`
)

type streamerReceiveDiamond struct{}

func (o *streamerReceiveDiamond) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		zlogger.Debugw("streamer receive diamond message", zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))
		o.receiveDiamond(ctx, msg)
	}
	return consumer.ConsumeSuccess, nil
}

func (o *streamerReceiveDiamond) receiveDiamond(ctx context.Context, msg *primitive.MessageExt) {
	jsonData := &queue.StreamerReceiveDiamond{}
	if err := json.Unmarshal(msg.Body, jsonData); err != nil {
		zlogger.Errorw("receiveDiamond::receiveDiamond, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
		return
	}

	roomCache, err := getRoomCache(jsonData.RoomId)
	if err != nil {
		zlogger.Errorw("receiveDiamond::receiveDiamond, get room cache error", zap.String("msgID", msg.MsgId), zap.Error(err))
		return
	}

	// 累积主播收到的钻石
	dbResult := mysql.LiveDB.WithContext(ctx).Model(&table.LiveRoom{}).
		Where("user_id = ?", jsonData.AnchorId).
		Update("accumulation_diamond", gorm.Expr("accumulation_diamond + ? ", jsonData.ProjectTotal))
	if dbResult.Error != nil {
		zlogger.Errorw("receiveDiamond::receiveDiamond, update room info error",
			zap.Int("userId", jsonData.AnchorId), zap.Error(dbResult.Error))
		return
	}
	if dbResult.RowsAffected == 0 {
		zlogger.Infow("receiveDiamond::receiveDiamond, update room info 0 effect",
			zap.Int("userId", jsonData.AnchorId), zap.Error(dbResult.Error))
		return
	}

	liveData := &table.LiveRoom{}
	if err = mysql.LiveDB.WithContext(ctx).Where("user_id = ?", jsonData.AnchorId).
		First(liveData).Error; err != nil {
		zlogger.Errorw("receiveDiamond::receiveDiamond, get room level error",
			zap.Int("userId", jsonData.AnchorId), zap.Error(err))
		return
	}

	accLevel := o.getCacheLevel(ctx, int(liveData.AccumulationDiamond))
	if accLevel <= roomCache.LevelId {
		return
	}

	if err = mysql.LiveDB.WithContext(ctx).Model(&table.LiveRoom{}).
		Where("user_id = ?", jsonData.AnchorId).
		Update("level_id", accLevel).Error; err != nil {
		zlogger.Errorw("receiveDiamond::receiveDiamond, update room level error",
			zap.Int("userId", jsonData.AnchorId), zap.Error(err))
		return
	}

	roomCache.LevelId = accLevel
	if err = setRoomCache(jsonData.RoomId, roomCache); err != nil {
		zlogger.Errorw("receiveDiamond::receiveDiamond, update room cache level error",
			zap.Int("userId", jsonData.AnchorId), zap.Error(err))
		return
	}
}

func (o *streamerReceiveDiamond) getCacheLevel(ctx context.Context, diamond int) int {
	levelList, err := redis.ZRevRangeByScore(redis2.LevelStreamerConfig, "-inf", strconv.Itoa(diamond))
	if len(levelList) == 0 {
		siteCfg := &table.SysSiteConfig{}
		levelContext := ""
		if err = mysql.LiveDB.WithContext(ctx).Where("config_code = ?", "finance_cfg_streamer_level").
			First(siteCfg).Error; err != nil {
			zlogger.Errorw("streamerReceiveDiamond::getCacheLevel, get db config diamond error", zap.Error(err))
			levelContext = LevelStreamer
		} else {
			levelContext = siteCfg.Content
		}

		var levelStreamer []int
		if err = json.Unmarshal([]byte(levelContext), &levelStreamer); err != nil {
			zlogger.Errorw("streamerReceiveDiamond::getCacheLevel, unmarshal db config error", zap.Error(err))
			return 0
		}

		retLevel := 0
		// 遍历等级和对应的钻石消费，将其写入 Redis 的 ZSet
		members := make([]redis3.Z, 0, len(levelStreamer))
		for level, d := range levelStreamer {
			members = append(members, redis3.Z{
				Score:  float64(d),
				Member: level + 1,
			})

			if diamond >= d {
				retLevel = level + 1
			}
		}

		// 消费钻石数作为 score
		if _, err = redis.ZAdd(redis2.LevelStreamerConfig, members...); err != nil {
			zlogger.Errorw("streamerReceiveDiamond::getCacheLevel, set config cache error", zap.Error(err))
			return 0
		}

		return retLevel
	}
	if err != nil {
		zlogger.Errorw("streamerReceiveDiamond::getCacheLevel, get level diamond config fail",
			zap.Int("diamond", diamond), zap.Error(err))
		return 0
	}

	accLevel, err := strconv.Atoi(levelList[0])
	if err != nil {
		zlogger.Errorw("streamerReceiveDiamond::getCacheLevel, transform int value fail",
			zap.Int("diamond", diamond), zap.Error(err))
		return 0
	}

	return accLevel
}
