package tasks

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"queueJob/pkg/common/config"
	jobStruct "queueJob/pkg/db/structs/job"
	"time"

	redisKey "queueJob/pkg/constant/redis"
	"queueJob/pkg/context"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/xxl"
	"queueJob/pkg/zlogger"
)

// JobGenerateIndex 定期生成首页
func JobGenerateIndex(cxt *context.Context, _ *xxl.RunReq) (msg string) {
	cxt.Trace = fmt.Sprintf("Job_Generate_Index_%s", cxt.Trace)

	zlogger.Infof("stats data JobGenerateIndex %v begin", time.Now())
	//start := time.Now()

	// 获取当前日期
	currentDate := time.Now().Format("2006-01-02")

	//endDate := BJNowTime()
	backGround := *cxt.Ctx

	// 检查 Redis 中存储的日期
	htmlRecords, err := redis.GetRecentItems(backGround, redisKey.RecentHtmlFilesId, redisKey.RecentHtmlFiles, 2000)
	if errors.Is(err, redis.Nil) {
		// 如果 Redis 中没有记录日期，初始化
		return "success"
	}
	if err != nil {
		zlogger.Errorf("request error: %v", err)
		return "failed"
	}

	// 确保目录存在
	err = os.MkdirAll(config.Config.Apk.TemplateJobIndexOneDir, os.ModePerm)
	if err != nil {
		zlogger.Errorf("request error: %v", err)
	}

	// 3. 加载模板
	tmpl, err := template.ParseFiles(config.Config.Apk.TemplateJobIndex)
	if err != nil {
		zlogger.Errorf("加载模板失败: %v", err)
	}

	// 构造文件路径，使用 punCompanyJob.Id 命名
	fileName := fmt.Sprintf(config.Config.Apk.TemplateJobIndexOneFile, currentDate)
	filePath := filepath.Join(config.Config.Apk.TemplateJobIndexOneDir, fileName) // 注意字段是 ID，不是 Id（Go 命名习惯）

	// 创建输出文件
	f, err := os.Create(filePath)
	if err != nil {
		zlogger.Errorf("创建HTML文件失败: %v", err)
	}
	defer f.Close()

	// 5. 渲染模板
	err = tmpl.Execute(f, htmlRecords)
	if err != nil {
		zlogger.Errorf("渲染模板失败：", err)
	}

	//zlogger.Infof("JobGenerateIndex start %v ", start)
	//zlogger.Infof("JobGenerateIndex currentDate %v ", currentDate)
	//zlogger.Infof("JobGenerateIndex htmlRecords %v ", htmlRecords)
	//zlogger.Infof("JobGenerateIndex endDate %v ", endDate)

	return "success"

}

func oldFlies4Test(user *jobStruct.TgMessageUser) string {

	// 获取当前可执行文件所在目录（更适用于生产环境）
	execPath, err := os.Getwd()
	if err != nil {
		panic("获取可执行文件路径失败：" + err.Error())
	}
	rootDir := filepath.Dir(execPath)

	// 拼接 public 目录路径
	publicDir := filepath.Join(rootDir, "../public")
	// 确保目录存在
	err = os.MkdirAll(publicDir, os.ModePerm)
	if err != nil {
		fmt.Printf("创建目录失败：%v\n", err)
	}

	// 加载模板
	tmplPath := filepath.Join(rootDir, "../template", "job_index.html")
	//tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		zlogger.Errorf("加载模板失败：%v", err)
	}

	// 构造文件路径，使用 punCompanyJob.Id 命名
	//filePath := filepath.Join(publicDir, fmt.Sprintf("job_index_%s.html", currentDate))
	//// 创建输出文件
	//f, err := os.Create(filePath)
	//if err != nil {
	//	zlogger.Errorf("创建HTML文件失败：%v", err)
	//}
	//defer f.Close()
	//
	//// 5. 渲染模板
	//err = tmpl.Execute(f, htmlRecords)
	//if err != nil {
	//	zlogger.Errorf("渲染模板失败：", err)
	//}
	//
	//zlogger.Infof("JobGenerateIndex start %v ", start)
	//zlogger.Infof("JobGenerateIndex currentDate %v ", currentDate)
	//zlogger.Infof("JobGenerateIndex htmlRecords %v ", htmlRecords)
	zlogger.Infof("JobGenerateIndex endDate %v ", tmplPath)
	return "success"
}
