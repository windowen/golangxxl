package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/bregydoc/gtranslate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 判断是否包含中文字符
func containsChinese(text string) bool {
	re := regexp.MustCompile(`[\p{Han}]`)
	return re.MatchString(text)
}

// 判断是否可能是柬埔寨语（高棉文）
func isProbablyKhmer(text string) bool {
	re := regexp.MustCompile("[\u1780-\u17FF]")
	return re.MatchString(text)
}

// 过滤一些无意义的常用词
var excludedWords = map[string]bool{
	"ok": true, "OK": true, "yes": true, "no": true,
	"thanks": true, "thank you": true, "please": true,
	"good": true, "hello": true, "nice": true, "haha": true,
	"1": true, "2": true, "3": true, "4": true, "5": true,
	"6": true, "7": true, "8": true, "9": true, "10": true,
	"11": true,
}

func main3() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("❌ 加载 .env 文件失败")
	// }
	//
	// token := os.Getenv("TELEGRAM_TOKEN")
	// if token == "" {
	// 	log.Fatal("❌ TELEGRAM_TOKEN 未设置")
	// }

	bot, err := tgbotapi.NewBotAPI("7652002179:AAGXJlF7vNq2njzeuQOg6Y4SKM-wKirOleU")
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	fmt.Println("✅ 机器人已启动，等待消息...")

	// 打开（或创建）日志文件
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer logFile.Close()

	// 设置日志输出到文件，并带上日期时间
	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// 示例日志
	logger.Println("程序启动")
	logger.Println("正在执行某些操作...")

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		text := strings.TrimSpace(update.Message.Text)

		// 忽略空白和排除词
		if len(text) <= 1 || excludedWords[strings.ToLower(text)] {
			continue
		}

		var translated, translatedEn string

		if containsChinese(text) {
			// 中文 → 柬埔寨语
			translated, err = gtranslate.TranslateWithParams(
				text,
				gtranslate.TranslationParams{
					From: "zh-CN",
					To:   "km",
				},
			)
			translatedEn, err = gtranslate.TranslateWithParams(
				text,
				gtranslate.TranslationParams{
					From: "zh-CN",
					To:   "en",
				},
			)
			if err == nil {
				msgEn := tgbotapi.NewMessage(update.Message.Chat.ID,
					" Cambodian:\n"+translated)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					text+" En \n"+translatedEn+"\n\n"+msgEn.Text)

				log.Printf("中文 → 柬埔寨语翻译后: %s", msg)

				logger.Println("中文 → 柬埔寨语翻译后: %s", msg.Text)
				bot.Send(msg)
			}

			log.Printf("中文 → 柬埔寨语: %s", text)
		} else if isProbablyKhmer(text) {
			// 柬埔寨语 → 中文
			translated, err = gtranslate.TranslateWithParams(
				text,
				gtranslate.TranslationParams{
					From: "km",
					To:   "zh-CN",
				},
			)
			translatedEn, err = gtranslate.TranslateWithParams(
				text,
				gtranslate.TranslationParams{
					From: "km",
					To:   "en",
				},
			)
			if err == nil {
				// msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				// 	"🇰🇭 Original "+text+"⬆️\n\n🀄️ Chinese:\n"+translated)
				// log.Printf("柬埔寨语 855→ 中文 翻译 后: %s", msg)

				msgEn := tgbotapi.NewMessage(update.Message.Chat.ID,
					" Cambodian:\n"+translated)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					text+" En \n"+translatedEn+"\n\n"+msgEn.Text)
				logger.Println("柬埔寨语 → 中文 翻译: %s", msg.Text)
				bot.Send(msg)
			}

			// if err == nil {
			// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			// 		"🇰🇭 Original "+text+"⬆️\n\n🀄️ Chinese:\n"+translated)
			// 	log.Printf("柬埔寨语 855→ 中文 翻译 后: %s", msg)
			// 	bot.Send(msg)
			// }
			log.Printf(" 柬埔寨语 → 中文 翻译  : %s", text)
		}
	}
}
