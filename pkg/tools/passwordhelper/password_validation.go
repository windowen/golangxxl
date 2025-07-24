package passwordhelper

import (
	"bufio"
	"go.uber.org/zap"
	"os"
	"queueJob/pkg/zlogger"
)

const WEAK_FILE = "weak_password_list.txt"

// CheckWeakPassword
/*
1、单数字+单字母，共八位
（0-9）（26个）
例如：1111aaaa 2222bbbb 3333zzzz

2、单字母+单数字，共八位
（26个）（0-9）
例如：aaaa1111 bbbb1111 tttt2222

3、abcd + 单数字，共八位
例如：abcd1111 abcd2222

4、qwer + 单数字，共八位
例如：qwer1111 qwer2222

5、asdf + 单数字，共八位
例如：asdf1111 asdf2222

6、zxcv + 单数字，共八位
例如：zxcv1111 zxcv2222
反过来亦是,同上

7、叠字母 + 单数字，共八位
例如：aabb1111 aabb2222 aacc5555 hhjj4444 ttaa1111
反过来亦是,同上*/
func CheckWeakPassword(password string) bool {
	commonPass := LoadWeakFile() // Load weak passwords file
	for _, pass := range commonPass {
		if pass == password {
			return false
		}
	}
	return true
}

func LoadWeakFile() []string {
	// Loads weak password file, and returns a slice
	file, err := os.Open(WEAK_FILE)
	if err != nil {
		zlogger.Errorw("LoadWeakFile", zap.Error(err))
	}
	a := make([]string, 0)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			zlogger.Errorw("load file error ： ", zap.Error(err))
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		a = append(a, scanner.Text())
	}
	return a
}
