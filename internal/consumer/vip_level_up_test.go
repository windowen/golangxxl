package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	redis3 "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	redis2 "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/db/table"
	"queueJob/pkg/message"
)

var (
	LiveDB *gorm.DB
)

func getCacheLevel(ctx context.Context, diamond int) int {
	levelList, err := redis.ZRevRangeByScore(redis2.LevelDiamondConfig, "-inf", strconv.Itoa(diamond))
	if len(levelList) == 0 {
		siteCfg := &table.SysSiteConfig{}
		levelContext := ""
		if err = LiveDB.WithContext(ctx).Where("config_code = ?", "finance_cfg_diamond_level").
			First(siteCfg).Error; err != nil {
			levelContext = LevelDiamond
		} else {
			levelContext = siteCfg.Content
		}

		var levelDiamond []int
		err = json.Unmarshal([]byte(levelContext), &levelDiamond)
		if err != nil {
			return 0
		}

		retLevel := 0
		// 遍历等级和对应的钻石消费，将其写入 Redis 的 ZSet
		for level, d := range levelDiamond {
			// 消费钻石数作为 score
			score := float64(d)
			_, err = redis.ZAdd(redis2.LevelDiamondConfig, redis3.Z{
				Score:  score,
				Member: level + 1,
			})
			if err != nil {
				return 0
			}

			if diamond >= d {
				retLevel = level
			}
		}

		return retLevel
	}
	if err != nil {
		return 0
	}

	accLevel, err := strconv.Atoi(levelList[0])
	if err != nil {
		return 0
	}

	return accLevel
}

func Test_vipLevelUp_handleMessages(t *testing.T) {

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
	ctx := context.Background()

	err = redis.InitRedisByConfig([]string{"18.167.121.240:6379"}, "ws20240726v3+4@8abi_^W")
	if err != nil {
		panic(err)
	}

	pMessage := &message.AccumulationDiamond{
		UserId: 8884468,
	}

	userCache, err := getUserCache(pMessage.UserId)
	if err != nil {
		return
	}

	userWallet := &table.SiteUserWallet{}
	if err = LiveDB.WithContext(ctx).Where("user_id = ?", pMessage.UserId).
		First(userWallet).Error; err != nil {
		return
	}

	accLevel := getCacheLevel(ctx, userWallet.AccumulationDiamond)
	if accLevel <= userCache.LevelId {
		return
	}

	if err = LiveDB.WithContext(ctx).Model(&table.User{}).
		Where("id = ?", pMessage.UserId).
		Update("level_id", accLevel).Error; err != nil {
		return
	}

	userCache.LevelId = accLevel
	err = setUserCache(pMessage.UserId, userCache)
	if err != nil {
		return
	}
}
