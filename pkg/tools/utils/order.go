package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GetPayStatus 充值订单状态
func GetPayStatus(status int) string {
	var result string

	switch status {
	case 0:
		result = "待处理"
	case 1:
		result = "处理中"
	case 2:
		result = "失败"
	case 3:
		result = "已成功"
	}

	return result
}

// GetWithdrawStatus 提现订单状态
func GetWithdrawStatus(status int) string {
	var result string

	switch status {
	case 0:
		result = "待处理"
	case 1:
		result = "处理中"
	case 3:
		result = "拒绝"
	case 4:
		result = "已成功"
	}

	return result
}

// GetConsumeName 通过类型Id获取消费名称
func GetConsumeName(category int) string {
	var result string

	switch category {
	case 1:
		result = "礼物消费"
	case 2:
		result = "弹幕消费"
	case 3:
		result = "游戏消费"
	case 4:
		result = "直播收费消费"
	case 5:
		result = "道具消费"
	}

	return result
}

// GetOrderNo 获取订单号
func GetOrderNo(prefix string) string {
	// 获取当前时间
	now := BJNowTime()
	// 获取毫秒部分
	milliseconds := now.Nanosecond() / 1e6
	// 格式化毫秒部分并拼接到时间字符串
	timeStr := fmt.Sprintf("%s%03d", now.Format("20060102150405"), milliseconds)

	// 生成4位随机数字
	randomNumber := rand.Intn(10000) // 生成0到9999之间的随机数

	// 拼接时间字符串和随机数字
	finalStr := fmt.Sprintf("%s%s%04d", prefix, timeStr, randomNumber)
	return finalStr
}

// GenerateSerialNumber generates the serial number based on the rules
func GenerateSerialNumber(category int) string {
	// Step 1: Generate the prefix "M"
	prefix := "M"

	// Step 2: Format the category as a 3-digit string (e.g., "001", "002")
	categoryStr := fmt.Sprintf("%03d", category)

	// Step 3: Get the current date and time in yyyyMMddHHmmssSSS format (17 characters, without dot)
	now := time.Now()
	milliseconds := now.Nanosecond() / 1e6
	timestamp := fmt.Sprintf("%s%03d", now.Format("20060102150405"), milliseconds)

	// Step 4: Generate a 4-digit random number
	randomNum := rand.Intn(10000)

	// Combine all parts into the final serial number
	serialNumber := fmt.Sprintf("%s%s%s%d", prefix, categoryStr, timestamp, randomNum)

	return serialNumber
}
