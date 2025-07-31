package tasks

import (
	"errors"
	"fmt"
	"time"

	redis2 "queueJob/pkg/constant/redis"
	"queueJob/pkg/context"
	"queueJob/pkg/db/cache"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/xxl"
	"queueJob/pkg/zlogger"
)

// JobGenerateIndex 定期生成首页
func JobGenerateIndex(cxt *context.Context, _ *xxl.RunReq) (msg string) {
	cxt.Trace = fmt.Sprintf("Job_Generate_Index_%s", cxt.Trace)

	zlogger.Infof("stats data sync mysql %v begin", time.Now())
	start := time.Now()

	// 获取当前日期
	currentDate := time.Now().Format("2006-01-02")

	endDate := BJNowTime()
	lastDates := endDate - 1000 // 默认从当前时间前推1秒开始

	backGround := *cxt.Ctx

	// 检查 Redis 中存储的日期
	lastDate, err := cache.RedisClient.Get(backGround, redis2.StatsSyncTime).Result()
	if errors.Is(err, redis.Nil) {
		// 如果 Redis 中没有记录日期，初始化
		if setStatus := cache.RedisClient.Set(backGround, redis2.StatsSyncTime, currentDate, 0); setStatus.Err() != nil {
			return "failed"
		}
	} else if err != nil {
		return "failed"
	} else {
	}
	zlogger.Infof("AceLotteryGameRecord begin %v ", start)
	zlogger.Infof("AceLotteryGameRecord begin %v ", currentDate)
	zlogger.Infof("AceLotteryGameRecord begin %v ", lastDate)
	zlogger.Infof("AceLotteryGameRecord begin %v ", lastDates)

	if err != nil {
		zlogger.Errorf("request error: %v", err)
		return "failed"
	}
	return "success"

}
