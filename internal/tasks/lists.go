package tasks

import (
	"liveJob/internal/tasks/settlement"
	"liveJob/pkg/xxl"
)

const (
	// 测试任务
	demoRebate = "demoRebate"
	// 直播每周结算
	liveSettlement = "liveSettlement"
	// 用户坐骑过期检查
	userMountExpiredCheck = "userMountExpiredCheck"
	// 更新redis里的统计数据到mysql
	statsDataSyncMysql = "statsDataSyncMysql"
)

// RegisterExecutors  注册执行器列表
func RegisterExecutors(execute *xxl.Executor) {
	// 测试任务
	// execute.RegTask(demoRebate, settlement.DemoRebate)
	// 直播每周结算
	// execute.RegTask(liveSettlement, settlement.LiveSettlement)
	// 用户坐骑过期检查
	execute.RegTask(userMountExpiredCheck, settlement.UserMountExpiredCheck)
	// 更新redis里的统计数据到mysql
	execute.RegTask(statsDataSyncMysql, StatsSyncMysql)
}
