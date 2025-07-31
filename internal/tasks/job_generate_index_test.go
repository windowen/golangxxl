package tasks

import (
	"queueJob/pkg/common/config"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/kafka"
	"queueJob/pkg/service"
	"queueJob/pkg/tools/utils"
	"testing"
	"time"

	"queueJob/pkg/context"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/zlogger"
)

func TestJobGenerateIndex(t *testing.T) {
	// 初始化配置
	configFile, logFile, err := config.FlagParsePath("queueJob", "../../config/config.yaml")
	if err != nil {
		panic(err)
	}
	err = config.InitConfig(configFile)
	if err != nil {
		panic(err)
	}

	// 设置全局时区为北京时间
	utils.SetGlobalTimeZone(utils.GetBjTimeLoc())

	// 初始化日志
	zlogger.InitLogConfig(logFile)

	// 初始化redis
	err = redis.InitRedis()
	if err != nil {
		panic(err)
	}

	// 初始化kafka生产者
	kafka.Init()

	// 初始化直播数据库mysql
	err = mysql.InitLiveDB()
	if err != nil {
		panic(err)
	}
	service.Start("job")
	defer service.Stop("job")

	JobGenerateIndex(context.Background(10*time.Minute), nil)
}
