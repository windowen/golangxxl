package main

import (
	"fmt"
	"os"
	"strings"
)

func specifyDirectory() []string {
	return []string{"live", "game", "finance"}
}

// 生成pb时，执行该文件
func main() {
	for _, dirName := range specifyDirectory() {
		dir, err := os.ReadDir("./" + dirName + "/")
		if err != nil {
			return
		}
		for _, file := range dir {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			if !strings.Contains(fileName, "pb.go") {
				continue
			}
			filepath := ("./" + dirName + "/") + fileName
			fileData, err := os.ReadFile(filepath)
			if err != nil {
				fmt.Printf("ReadFile err:%s \n", err.Error())
				return
			}
			data := string(fileData)

			// fmt.Println("需要替换的文件如下:" + fileName)
			data = strings.ReplaceAll(data, ",omitempty", "")
			fileData = []byte(data)
			err = os.WriteFile(filepath, fileData, 0644)
			if err != nil {
				fmt.Printf("WriteFile err:%s \n ", err.Error())
				return
			}
		}
	}

}
