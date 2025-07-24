package kafkaconsumer

import (
	"fmt"
	"testing"
	"time"

	"queueJob/pkg/common/config"
	"queueJob/pkg/zlogger"
)

func TestSendFacebookEvent(t *testing.T) {
	Id := "1380876816448536"
	token := "EAAOCZA9A0dcQBO4vkZAupopdIFu4FQZAlgy07KiiFfXvdyKtjWTJudvaciJC1vE7Vx7VYJnQS4Hnp90Lyrie9MjNpUoMZBgZBZBTID5ZCttAYXHG1bXc1VVAV9ZCxf5KeLIK7s1mGHoMh37V3654NWI5sasakM03BEwdR7GwHpwanPEhbmJiRZBWq02mXhajeFKYAJQZDZD"

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

	// if err = SendFacebookEvent(Id, token, &FacebookEvent{
	// 	EventName:      "Lead",
	// 	EventTime:      time.Now().Unix(),
	// 	ActionSource:   "website",
	// 	EventSourceURL: "https://www.facebook.com",
	// 	UserData: &UserData{
	// 		Emails:          []string{"mmm@gmail.com"},
	// 		ClientIP:        "10.0.0.1",
	// 		ClientUserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	// 	},
	// }); err != nil {
	// 	fmt.Println(err)
	// }

	if err = SendFacebookEvent(Id, token, &FacebookEvent{
		EventName:      "Purchase",
		EventTime:      time.Now().Unix(),
		ActionSource:   "website",
		EventSourceURL: "https://www.facebook.com",
		UserData: &UserData{
			Emails:          []string{"ccc@gmail.com"},
			ClientIP:        "10.0.0.1",
			ClientUserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		},
		CustomData: &CustomData{
			Currency: "usd",
			Value:    12,
		},
	}); err != nil {
		fmt.Println(err)
	}

	if err = SendFacebookEvent(Id, token, &FacebookEvent{
		EventName:      "firstDeposit",
		EventTime:      time.Now().Unix(),
		ActionSource:   "website",
		EventSourceURL: "https://www.facebook.com",
		UserData: &UserData{
			Emails:          []string{"ccc@gmail.com"},
			ClientIP:        "10.0.0.1",
			ClientUserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		},
		CustomData: &CustomData{
			Currency: "usd",
			Value:    12,
		},
	}); err != nil {
		fmt.Println(err)
	}
}
