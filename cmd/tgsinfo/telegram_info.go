package main

import (
	"fmt"
	"log"
	"queueJob/pkg/common/config"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/tools/utils"
	"queueJob/pkg/zlogger"
	"regexp"
	"runtime"
	"strings"

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

				bot.Send(msg)
			}
		}

		log.Printf("中文 → 柬埔寨语翻译后2: %s", text)

		//logger.Println("中文 → 柬埔寨语翻译后:", text)
	}
}
