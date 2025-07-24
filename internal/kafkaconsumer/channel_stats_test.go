package kafkaconsumer

import (
	"encoding/json"
	"testing"
	"time"

	"queueJob/pkg/common/config"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/redisdb/redis"

	"queueJob/pkg/queue"
	"queueJob/pkg/zlogger"
)

func Test_channelStats_handleMessages(t *testing.T) {
	// 初始化配置
	configFile, logFile, err := config.FlagParsePath("queue", "../../config/config.yaml")
	if err != nil {
		panic(err)
	}

	if err = config.InitConfig(configFile); err != nil {
		panic(err)
	}

	// 初始化日志
	zlogger.InitLogConfig(logFile)

	// 初始化 Redis
	if err = redis.InitRedis(); err != nil {
		panic(err)
	}

	// 初始化直播数据库 MySQL
	if err = mysql.InitLiveDB(); err != nil {
		panic(err)
	}

	// 逻辑测试
	o := &channelStats{}

	jsonData := &queue.ChannelStats{
		StatsType:    6,
		UserId:       0,
		ChannelCode:  10086,
		DeviceId:     "RjvAzVvfXBwUT2hY0yWq",
		RegisterTime: time.Now(),
		Active1:      1,
	}

	data, err := json.Marshal(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	o.handleMessages(data)
}
