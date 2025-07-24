package i18nh

import (
	"liveJob/pkg/constant"
)

func LangToString(langId int) string {
	var langMap map[int]string
	langMap = make(map[int]string)
	langMap[1] = constant.LangChinese
	langMap[2] = constant.LangEnglish
	langMap[3] = constant.LangID
	langMap[4] = constant.LangVi
	langMap[5] = constant.LangPT
	if l, ok := langMap[langId]; ok {
		return l
	}
	return ""
}
