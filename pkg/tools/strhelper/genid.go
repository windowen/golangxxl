package strhelper

import (
	"crypto/rand"

	"liveJob/pkg/common/config"
)

func GenUserID() string {
	const l = 10
	data := make([]byte, l)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		if i == 0 {
			data[i] = chars[1:][data[i]%9]
		} else {
			data[i] = chars[data[i]%10]
		}
	}
	return string(data)
}

func GenVerifyCode() string {
	data := make([]byte, config.Config.Sms.CodeLen)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		data[i] = chars[data[i]%10]
	}
	return string(data)
}
