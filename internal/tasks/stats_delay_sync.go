package tasks

import (
	ctx "context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	redis2 "liveJob/pkg/constant/redis"
	"liveJob/pkg/context"
	"liveJob/pkg/db/cache"
	"liveJob/pkg/db/mysql"
	"liveJob/pkg/db/redisdb/redis"
	"liveJob/pkg/db/table"
	"liveJob/pkg/xxl"
	"liveJob/pkg/zlogger"
)

// StatsSyncMysql 统计数据延时同步到mysql
func StatsSyncMysql(cxt *context.Context, _ *xxl.RunReq) (msg string) {
	cxt.Trace = fmt.Sprintf("Stats_Sync_Mysql_%s", cxt.Trace)

	zlogger.Infof("stats data sync mysql %v begin", time.Now())
	start := time.Now()

	// 获取当前日期
	currentDate := time.Now().Format("2006-01-02")

	backGround := *cxt.Ctx

	isNewDay := false
	// 检查 Redis 中存储的日期
	lastDate, err := cache.RedisClient.Get(backGround, redis2.StatsSyncTime).Result()
	if errors.Is(err, redis.Nil) {
		// 如果 Redis 中没有记录日期，初始化
		if setStatus := cache.RedisClient.Set(backGround, redis2.StatsSyncTime, currentDate, 0); setStatus.Err() != nil {
			return "failed"
		}
		isNewDay = true
	} else if err != nil {
		return "failed"
	} else {
		isNewDay = lastDate != currentDate
	}

	zlogger.Debugw("StatsSyncMysql, is new day", zap.Bool("isNewDay", isNewDay))

	if isNewDay {
		// 如果跨天，把昨天的redis记录更新到mysql
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02") // 当前日期减去1天
		ok, dbData := getRedisData(backGround, yesterday)
		if ok {
			if !isRecordExists(backGround, yesterday) {
				if err = mysql.LiveDB.WithContext(backGround).Create(dbData).Error; err != nil {
					zlogger.Errorw("StatsSyncMysql, create db data failed", zap.Error(err))
					return "failed"
				}
			} else {
				if err = mysql.LiveDB.WithContext(backGround).Model(&table.ChannelDailyUserStats{}).Where("report_date = ?", yesterday).Updates(dbData).Error; err != nil {
					zlogger.Errorw("StatsSyncMysql, update db data failed", zap.Error(err))
					return "failed"
				}
			}
			clearRedisData(backGround, yesterday)
		}
	}

	// 更新redis里的数据到mysql
	ok, dbData := getRedisData(backGround, currentDate)
	if ok {
		if !isRecordExists(backGround, currentDate) {
			if err = mysql.LiveDB.WithContext(backGround).Create(dbData).Error; err != nil {
				zlogger.Errorw("StatsSyncMysql, create db data failed", zap.Error(err))
				return "failed"
			}
		} else {
			if err = mysql.LiveDB.WithContext(backGround).Model(&table.ChannelDailyUserStats{}).Where("report_date = ?", currentDate).Updates(dbData).Error; err != nil {
				zlogger.Errorw("StatsSyncMysql, update db data failed", zap.Error(err))
				return "failed"
			}
		}
	}

	end := time.Since(start)

	if setStatus := cache.RedisClient.Set(backGround, redis2.StatsSyncTime, currentDate, 0); setStatus.Err() != nil {
		zlogger.Errorw("StatsSyncMysql, set db data failed", zap.Error(setStatus.Err()))
		return "failed"
	}

	zlogger.Infof("stats data sync mysql end：%v", end)
	return "success"
}

func getRedisData(background ctx.Context, day string) (bool, *table.ChannelDailyUserStats) {
	dayTime, err := time.Parse("2006-01-02", day)
	if err != nil {
		zlogger.Errorw("StatsSyncMysql, parse day time failed", zap.String("day", day), zap.Error(err))
		return false, nil
	}

	record := &table.ChannelDailyUserStats{
		ReportDate: dayTime,
	}

	hGetAll := cache.RedisClient.HGetAll(background, fmt.Sprintf(redis2.LiveStats, day))
	if hGetAll.Err() != nil {
		return false, nil
	}

	for key, val := range hGetAll.Val() {
		count, err := strconv.Atoi(val)
		if err != nil {
			return false, nil
		}

		switch key {
		case redis2.EventUserRegistrations:
			record.UserRegistrations = count
		case fmt.Sprintf(redis2.EventPageStayTime, "Game"):
			record.PageStayTime = count
		case redis2.EventGameLaunchCount:
			record.GameLaunchCount = count
		case redis2.EventGameAwardCount:
			record.GameAwardCount = count
		case redis2.EventHomepageBannerClicks:
			record.HomepageBannerClicks = count
		case redis2.EventRecommendedBannerClicks:
			record.RecommendedBannerClicks = count
		case redis2.EventPopularBannerClicks:
			record.PopularBannerClicks = count
		}
	}

	return true, record
}

// 判断记录是否存在
func isRecordExists(background ctx.Context, date string) bool {
	var count int64

	// 查询记录数量
	if err := mysql.LiveDB.WithContext(background).Model(&table.ChannelDailyUserStats{}).Where("report_date = ?", date).Count(&count).Error; err != nil {
		zlogger.Errorw("StatsSyncMysql, count error", zap.Error(err))
		return false
	}

	return count > 0
}

// 清空上一天的数据
func clearRedisData(background ctx.Context, day string) bool {
	_, err := cache.RedisClient.Del(background, fmt.Sprintf(redis2.LiveStats, day)).Result()
	if err != nil {
		zlogger.Errorw("StatsSyncMysql, clear db data failed", zap.String("day", day), zap.Error(err))
	}
	return true
}
