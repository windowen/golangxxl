package i18nh

import (
	"context"
	"embed"
	"fmt"
	"log"
	"queueJob/pkg/common/mctx"
	"queueJob/pkg/zlogger"
	"strconv"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml"
	"golang.org/x/text/language"

	"queueJob/pkg/constant"
)

type (
	TranslateFunc = func(msgId string, a ...interface{}) string
)

//go:embed lang/*.toml
var LocaleFS embed.FS

var onceLangInstance sync.Once
var object *langInstance
var LangPack = []string{constant.LangChinese, constant.LangEnglish, constant.LangID, constant.LangVi, constant.LangPT}

type langInstance struct {
	localizes map[language.Tag]*i18n.Localizer
	functions map[language.Tag]TranslateFunc
	lange     string
}

func New() *langInstance {
	onceLangInstance.Do(func() {
		object = &langInstance{
			localizes: make(map[language.Tag]*i18n.Localizer),
			functions: make(map[language.Tag]TranslateFunc),
			lange:     "en",
		}
		err := object.load(LangPack)
		if err != nil {
			log.Fatal(err)
		}
	})
	return object
}

// Load 载入指定语言的翻译文件
func (obj *langInstance) load(langs []string) error {
	for _, lang := range langs {
		tag, err := language.All.Parse(lang)
		if err != nil {
			zlogger.Error("文件加载错误:", err)
			return fmt.Errorf("[local18n] Load language bundle '%s' failed, %s", lang, err)
		}

		b := i18n.NewBundle(tag)
		b.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		file := fmt.Sprintf("lang/%s.toml", lang)

		_, err = b.LoadMessageFileFS(LocaleFS, file)
		if err != nil {
			zlogger.Error("文件加载错误:", err)
			panic(err)
			//  return
		}

		obj.localizes[tag] = i18n.NewLocalizer(b, tag.String())
	}
	return nil
}

// T 返回指定语言的翻译结果
func T(c context.Context, msgId string, a ...interface{}) string {
	language := mctx.GetLanguage(c)
	for _, v := range LangPack {
		if v == language {
			object.lange = v
			break
		}
	}
	return object.lang(object.getLang())(msgId, a...)
}

// lang 返回指定语言的翻译函数
func (obj *langInstance) lang(lang string) TranslateFunc {
	tag, err := language.All.Parse(lang)
	if err != nil {
		return obj.defaultTranslate
	}
	local, ok := obj.localizes[tag]
	if !ok {
		return obj.defaultTranslate
	}
	f, ok := obj.functions[tag]
	if !ok {
		f = obj.createTranslateFunc(local)
		obj.functions[tag] = f
	}
	return f
}

func (obj *langInstance) getLang() string {
	return obj.lange
}

// // local
func (obj *langInstance) createTranslateFunc(l *i18n.Localizer) TranslateFunc {
	return func(msgId string, a ...interface{}) string {
		t := map[string]interface{}{}
		for i := range a {
			t["v"+strconv.Itoa(i+1)] = a[i]
		}
		r, err := l.Localize(&i18n.LocalizeConfig{MessageID: msgId, TemplateData: t})
		if err == nil {
			return r
		}
		return obj.defaultTranslate(msgId, a...)
	}
}

func (obj *langInstance) defaultTranslate(msg string, a ...interface{}) string {
	return fmt.Sprintf("#local18n miss# %s", msg)
}
