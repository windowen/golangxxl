package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/google/uuid"

	"queueJob/internal/tasks"
	"queueJob/pkg/common/config"
	internal "queueJob/pkg/context"
	"queueJob/pkg/db/cache"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/middleware"
	"queueJob/pkg/tools/utils"
	"queueJob/pkg/xxl"
	"queueJob/pkg/zlogger"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			tmpStr := fmt.Sprintf("err=%v panic ==> %s\n", err, string(buf[:n]))
			fmt.Println(tmpStr)
		}
	}()
	// 初始化配置
	configFile, logFile, err := config.FlagParse("queueJob")
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
	svcID := uuid.NewString()
	zlogger.SetGlobalFields(map[string]string{
		"ssid": svcID,
		"snm":  "queueJob",
	})
	zlogger.InitLogConfig(logFile)

	// 初始化redis
	err = redis.InitRedis()
	if err != nil {
		panic(err)
	}

	// 初始化直播数据库mysql
	err = mysql.InitLiveDB()
	if err != nil {
		panic(err)
	}

	// 初始化redis
	err = cache.NewRedis()
	if err != nil {
		panic(err)
	}

	// 初始化XXLJOb数据库mysql
	err = mysql.InitXXLJobDB()
	if err != nil {
		panic(err)
	}

	execute := xxl.CreateExecutor(
		xxl.ServerAddr(config.Config.XXLJob.AdminServer),    // xxl-job-admin 服务器地址
		xxl.AccessToken(config.Config.XXLJob.AccessToken),   // 请求令牌(默认为空)
		xxl.RegistryKey(config.Config.XXLJob.ExecutorName),  // 执行器名称 本项目使用 liveExecutor
		xxl.ExecutorIp(config.Config.XXLJob.ExecutorIp),     // 可自动获取
		xxl.ExecutorPort(config.Config.XXLJob.ExecutorPort), // 该脚本执行器端口号默认9999（非必填）
		xxl.SetLogger(internal.Context{}),                   // 自定义日志
	)
	execute.Init()
	// 设置使用自定义中间件
	execute.Use(middleware.CustomMiddleware)
	// 设置日志查看handler
	execute.LogHandler(xxl.GetDBLogHandle)
	// 注册任务列表
	tasks.RegisterExecutors(execute)

	log.Fatal(execute.Run())
}
