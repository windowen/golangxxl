package consumer

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	redis3 "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	redis2 "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/db/table"
	"queueJob/pkg/message"
	rpcClient "queueJob/pkg/rpcclient"
	"queueJob/pkg/zlogger"
)

var (
	vipLevel     = &vipLevelUp{}
	LevelDiamond = `[10,30,70,190,385,640,1040,1665,2715,4290,6690,9840,14840,20840,27340,34840,42000,49500,57500,66000,75000,84500,94500,105000,115500,126000,138000,150500,163000,176000,189000,203000,218000,234000,251100,268000,286000,305000,325000,346000,368000,391000,415000,440000,466000,493000,521000,550000,580000,611000]`
)

type vipLevelUp struct{}

func (o *vipLevelUp) getCacheLevel(ctx context.Context, diamond int) int {
	levelList, err := redis.ZRevRangeByScore(redis2.LevelDiamondConfig, "-inf", strconv.Itoa(diamond))
	if len(levelList) == 0 {
		siteCfg := &table.SysSiteConfig{}
		levelContext := ""
		if err = mysql.LiveDB.WithContext(ctx).Where("config_code = ?", "finance_cfg_diamond_level").
			First(siteCfg).Error; err != nil {
			zlogger.Errorw("vipLevelUp::getCacheLevel, get db config diamond error", zap.Error(err))
			levelContext = LevelDiamond
		} else {
			levelContext = siteCfg.Content
		}

		var levelDiamond []int
		if err = json.Unmarshal([]byte(levelContext), &levelDiamond); err != nil {
			zlogger.Errorw("vipLevelUp::getCacheLevel, unmarshal db config error", zap.Error(err))
			return 0
		}

		retLevel := 0
		// 遍历等级和对应的钻石消费，将其写入 Redis 的 ZSet
		members := make([]redis3.Z, 0, len(levelDiamond))
		for level, d := range levelDiamond {
			members = append(members, redis3.Z{
				Score:  float64(d),
				Member: level + 1,
			})

			if diamond >= d {
				retLevel = level + 1
			}
		}

		// 消费钻石数作为 score
		if _, err = redis.ZAdd(redis2.LevelDiamondConfig, members...); err != nil {
			zlogger.Errorw("vipLevelUp::getCacheLevel, set config cache error", zap.Error(err))
			return 0
		}

		return retLevel
	}
	if err != nil {
		zlogger.Errorw("vipLevelUp::getCacheLevel, get level diamond config fail",
			zap.Int("diamond", diamond), zap.Error(err))
		return 0
	}

	accLevel, err := strconv.Atoi(levelList[0])
	if err != nil {
		zlogger.Errorw("vipLevelUp::getCacheLevel, transform int value fail",
			zap.Int("diamond", diamond), zap.Error(err))
		return 0
	}

	return accLevel
}

func (o *vipLevelUp) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		zlogger.Debugw("vipLevelUp::handleMessages, vip level up message",
			zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))

		pMessage := &message.AccumulationDiamond{}
		if err := json.Unmarshal(msg.Body, pMessage); err != nil {
			zlogger.Errorw("vipLevelUp::handleMessages, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
			continue
		}

		userCache, err := getUserCache(pMessage.UserId)
		if err != nil {
			zlogger.Errorw("vipLevelUp::handleMessages, get user cache error", zap.String("msgID", msg.MsgId), zap.Error(err))
			continue
		}

		userWallet := &table.SiteUserWallet{}
		if err = mysql.LiveDB.WithContext(ctx).Where("user_id = ?", pMessage.UserId).
			First(userWallet).Error; err != nil {
			zlogger.Errorw("vipLevelUp::handleMessages, get user wallet error",
				zap.String("msgID", msg.MsgId), zap.Int("userId", pMessage.UserId), zap.Error(err))
			continue
		}

		accLevel := o.getCacheLevel(ctx, userWallet.AccumulationDiamond)
		if accLevel <= userCache.LevelId {
			continue
		}

		if err = mysql.LiveDB.WithContext(ctx).Model(&table.User{}).
			Where("id = ?", pMessage.UserId).
			Update("level_id", accLevel).Error; err != nil {
			zlogger.Errorw("vipLevelUp::handleMessages, update user level error",
				zap.String("msgID", msg.MsgId), zap.Int("userId", pMessage.UserId), zap.Error(err))
			continue
		}

		userCache.LevelId = accLevel
		if err = setUserCache(pMessage.UserId, userCache); err != nil {
			zlogger.Errorw("vipLevelUp::handleMessages, update user cache level error",
				zap.String("msgID", msg.MsgId), zap.Int("userId", pMessage.UserId), zap.Error(err))
			continue
		}

		// 发送房间升级通知
		roomCache, err := getRoomCache(pMessage.RoomId)
		if err == nil {
			if err = rpcClient.ServiceClientsInstance.LiveClient.UpgradeNotifyWrap(ctx, pMessage.UserId, accLevel, roomCache.Id); err != nil {
				zlogger.Errorw("vipLevelUp::handleMessages, notify error", zap.Int("userId", pMessage.UserId), zap.Error(err))
				continue
			}
		}

		zlogger.Infow("vipLevelUp::handleMessages, update user level ok",
			zap.Int("userId", pMessage.UserId), zap.Int("level", accLevel))
	}
	return consumer.ConsumeSuccess, nil
}
