package utils

import (
	"crypto/rand"
	"math/big"
	mathRand "math/rand"
	"strconv"
)

func SixDigitNum() string {
	// 生成一个六位随机数字
	num, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return ""
	}

	// 确保数字总是六位长
	return strconv.FormatInt(num.Int64()+100000, 10)
}

// Random 获取随机时间
func Random(min, max int) int {
	randInt := min + mathRand.Intn(max-min+1)

	return randInt
}
