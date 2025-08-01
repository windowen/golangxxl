package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"queueJob/pkg/common/config"
	redisKey "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/redisdb/redis"
	jobStruct "queueJob/pkg/db/structs/job"
	"queueJob/pkg/db/table/job"
	"queueJob/pkg/tools/utils"
	"queueJob/pkg/zlogger"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 判断是否包含中文字符
func containChinese(text string) bool {
	re := regexp.MustCompile(`[\p{Han}]`)
	return re.MatchString(text)
}

// 过滤一些无意义的常用词
var excludedWord = map[string]bool{
	"ok": true, "OK": true, "yes": true, "no": true,
	"thanks": true, "thank you": true, "please": true,
	"good": true, "hello": true, "nice": true, "haha": true,
	"1": true, "2": true, "3": true, "4": true, "5": true,
	"6": true, "7": true, "8": true, "9": true, "10": true,
	"11": true,
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			tmpStr := fmt.Sprintf("err=%v panic ==> %s\n", err, string(buf[:n]))
			zlogger.Errorf(tmpStr)
		}
	}()

	ctx := context.Background()
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

	// 初始化XXLJOb数据库mysql
	err = mysql.InitXXLJobDB()
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI("7738584148:AAHMa9qeab148hWeJPCLD6w9VE9o6gBr_L4")
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	//fmt.Println("✅ 机器人已启动，等待消息...")

	// 打开（或创建）日志文件
	//logFile, err := os.OpenFile("app_info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	log.Fatalf("无法打开日志文件: %v", err)
	//}
	//defer logFile.Close()
	//
	//// 设置日志输出到文件，并带上日期时间
	//logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	//
	//// 示例日志
	//logger.Println("程序启动")
	//logger.Println("正在执行某些操作...")

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		text := strings.TrimSpace(update.Message.Text)

		// 忽略空白和排除词
		if len(text) <= 1 || excludedWord[strings.ToLower(text)] {
			continue
		}

		var translatedEn string

		if containChinese(text) {

			if err == nil {

				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					text+" En \n"+translatedEn+"\n\n")

				prettyJSON, _ := json.MarshalIndent(update.Message, "", "  ")

				var msgStruct jobStruct.TgMessage
				_ = json.Unmarshal(prettyJSON, &msgStruct)

				if isBot(msgStruct.ForwardFrom) {
					// 对于机器人消息 进行过滤
					continue
				}
				log.Printf("Message details:\n%s", string(prettyJSON))

				//bot.Send(msg)
				punCompanyJob := ConvertToPunCompanyJob(update.Message.Text)
				punCompanyJob.Uid = config.Config.Apk.UId
				punCompanyJob.ComName = msgStruct.From.Username
				if punCompanyJob.Name != "" && len(punCompanyJob.Name) > 0 {
					err = mysql.LiveDB.WithContext(ctx).Create(&punCompanyJob).Error
				}

				if err != nil {
					zlogger.Warnf("getGameIntegratorByCode query data failed | integratorCode:%v | err:%v", msg, err)
				}

				// 3. 加载模板
				tmpl, err := template.ParseFiles(config.Config.Apk.TemplateJob)
				if err != nil {
					log.Fatal("加载模板失败：", err)
				}

				// 确保目录存在
				err = os.MkdirAll("public", os.ModePerm)
				if err != nil {
					log.Fatal("创建目录失败：", err)
				}

				// 构造文件路径，使用 punCompanyJob.Id 命名
				filePath := fmt.Sprintf(config.Config.Apk.TemplateJobOne, punCompanyJob.Id) // 注意字段是 ID，不是 Id（Go 命名习惯）

				// 创建输出文件
				f, err := os.Create(filePath)
				if err != nil {
					log.Fatalf("创建HTML文件失败：%v", err)
				}
				defer f.Close()

				// 5. 渲染模板
				err = tmpl.Execute(f, punCompanyJob)
				if err != nil {
					log.Fatal("渲染模板失败：", err)
				}
				//fmt.Println("首页静态HTML生成成功！")

				SaveHTMLRecord(ctx, punCompanyJob, filePath)

			}
		}

		log.Printf("中文 → 柬埔寨语翻译后2: %s", text)
		//log.Printf("中文 → 柬埔寨语翻译后2: %v", update.Message)
		//log.Printf("翻译结果: ID=%d | Text=%s | Sender=%v",
		//	update.Message.ID,
		//	update.Message.Text,
		//	update.Message.Sender)

		//logger.Println("中文 → 柬埔寨语翻译后:", text)
	}
}
func ConvertToPunCompanyJob(text string) *job.PunCompanyJob {
	now := time.Now().Unix()
	job := &job.PunCompanyJob{
		// 设置默认值
		SDate:    int(now),
		EDate:    int(now + 30*24*60*60), // 默认30天后过期
		State:    0,                      // 审核状态（1通过，0未审核，2未通过）
		Status:   1,                      // 职位状态（1显示，0隐藏）
		Type:     1,                      // 全职
		Sex:      3,                      // 性别不限
		Source:   3,                      // apk来源
		LinkType: 1,                      // 默认联系方式
		Pr:       20,                     // 默认企业性质
		Mun:      30,                     // 默认企业规模
	}

	// 使用正则表达式从文本中提取关键信息
	re := regexp.MustCompile(`岗位名称：(.+)\n工作地点：(.+)\n薪资范围：(.+)\n人数需求：(.+)人`)
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 5 {
		job.Name = strings.TrimSpace(matches[1]) // 岗位名称
		//location := strings.TrimSpace(matches[2])  // 工作地点
		salary := strings.TrimSpace(matches[3])    // 薪资范围
		headcount := strings.TrimSpace(matches[4]) // 人数需求

		// 处理工作地点 - 暂时直接存储，实际应该解析为地区ID
		//job.ComName = location
		job.ComName = config.Config.Apk.ComName //公司名字

		// 处理薪资范围
		if salaryParts := strings.Split(salary, "-"); len(salaryParts) == 2 {
			min, _ := strconv.Atoi(strings.TrimSpace(strings.TrimRight(salaryParts[0], " 美金")))
			max, _ := strconv.Atoi(strings.TrimSpace(strings.TrimRight(salaryParts[1], " 美金")))
			job.MinSalary = min
			job.MaxSalary = max
		}

		// 处理人数需求
		if num, err := strconv.Atoi(headcount); err == nil {
			job.Number = num
		}
	}

	// 提取岗位职责和要求部分
	if sections := strings.Split(text, "岗位职责："); len(sections) > 1 {
		parts := strings.Split(sections[1], "岗位要求：")
		if len(parts) >= 2 {
			responsibilities := strings.TrimSpace(parts[0])
			requirements := strings.TrimSpace(parts[1])
			job.Description = fmt.Sprintf("岗位职责:\n%s\n\n岗位要求:\n%s",
				responsibilities, requirements)

			// 从要求中提取经验年限
			if expMatch := regexp.MustCompile(`(\d+)年以上`).FindStringSubmatch(requirements); len(expMatch) > 1 {
				if exp, err := strconv.Atoi(expMatch[1]); err == nil {
					job.Exp = exp
					job.ExpReq = fmt.Sprintf("%s年以上工作经验", expMatch[1])
				}
			}
		}
	}

	return job
}

func isBot(user *jobStruct.TgMessageUser) bool {
	if user == nil {
		return false
	}
	return user.IsBot
}

func SaveHTMLRecord(ctx context.Context, job *job.PunCompanyJob, path string) error {

	record := jobStruct.HTMLRecord{
		Title:     job.Name,
		Id:        strconv.Itoa(job.Id),
		Path:      path,
		Author:    job.ComName,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	//jsonData, err := json.Marshal(record)
	//if err != nil {
	//	return fmt.Errorf("JSON序列化失败: %v", err)
	//}

	// 使用流水线操作：先 LPUSH，再 LTRIM 保持最多 100 条
	// 封装调用
	//data := map[string]string{"name": "Alice", "action": "login"}
	//err = redis.PushAndTrimList(ctx, redisKey.RecentHtmlFilesId, redisKey.RecentHtmlFiles, record.Id, string(jsonData), 100)
	err := redis.PushAndTrimList(ctx, redisKey.RecentHtmlFilesId, redisKey.RecentHtmlFiles, record.Id, record, config.Config.Apk.MaxJobIndex, time.Duration(config.Config.Apk.MaxDays)*24*time.Hour)
	if err != nil {
		fmt.Println("操作失败:", err)
		return fmt.Errorf("JSON序列化失败: %v", err)
	}
	return err
}

//
//func GetRecentHTMLRecords(ctx context.Context, rdb *redis.Client) ([]HTMLRecord, error) {
//	const key = "recent_html_files"
//	results, err := rdb.LRange(ctx, key, 0, -1).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	var records []HTMLRecord
//	for _, item := range results {
//		var record HTMLRecord
//		if err := json.Unmarshal([]byte(item), &record); err == nil {
//			records = append(records, record)
//		}
//	}
//	return records, nil
//}

//// 更新或插入数据
//if report.Id > 0 {
//err = mysql.LiveDB.WithContext(ctx).Model(&punCompanyJob).Updates(&report).Error
//} else {
//err = mysql.LiveDB.WithContext(ctx).Create(&punCompanyJob).Error
//}
//db := mysql.LiveDB.WithContext(ctx)
//
//err = db.Where("code = ?", integratorCode).First(&integratorData).Error
