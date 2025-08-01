package main

import (
	"fmt"
	"os"
	"os/signal"
	"queueJob/internal/kafkaconsumer"
	"queueJob/pkg/kafka"
	"queueJob/pkg/rocketmq"
	"runtime"
	"syscall"

	"queueJob/pkg/agora"
	"queueJob/pkg/common/config"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/redisdb/redis"
	liveRpc "queueJob/pkg/rpcclient"
	"queueJob/pkg/service"
	"queueJob/pkg/tools/component"
	"queueJob/pkg/tools/utils"
	"queueJob/pkg/zlogger"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			tmpStr := fmt.Sprintf("err=%v panic ==> %s\n", err, string(buf[:n]))
			fmt.Println(tmpStr)
			zlogger.Error(tmpStr) // 记录到日志
		}
	}()

	// 初始化配置
	configFile, logFile, err := config.FlagParse("queue")
	if err != nil {
		panic(err)
	}

	if err = config.InitConfig(configFile); err != nil {
		panic(err)
	}

	// 设置全局时区为北京时间
	utils.SetGlobalTimeZone(utils.GetBjTimeLoc())

	// 初始化日志
	zlogger.InitLogConfig(logFile)

	// 检测 ETCD
	if err = component.ComponentCheck(configFile, config.Config.App.Discovery, true); err != nil {
		panic(err)
	}

	// 初始化 Redis
	if err = redis.InitRedis(); err != nil {
		panic(err)
	}

	// 初始化直播数据库 MySQL
	if err = mysql.InitLiveDB(); err != nil {
		panic(err)
	}

	// 初始化rtc客户端
	agora.NewRtcClient()

	// 初始化 RocketMQ 生产者
	rocketmq.Init()

	// 初始化 RocketMQ 消费者
	//consumer.Init()

	// 初始化 kafka 消费者
	kafkaconsumer.Init()

	// 初始化kafka生产者
	kafka.Init()

	// 初始化全局客户端
	service.RegisterService(liveRpc.NewServiceClients())

	service.Start("queue")
	defer service.Stop("queue")

	// 等待信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
}
