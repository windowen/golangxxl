package context

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/rs/xid"

	"liveJob/pkg/zlogger"
)

type Context struct {
	Ctx        *context.Context
	CancelFunc context.CancelFunc
	logs       strings.Builder
	Trace      string
}

func Background(timeout time.Duration) *Context {
	background := context.Background()

	if timeout > 0 {
		ctx, cancelFunc := context.WithTimeout(background, timeout)
		return &Context{Ctx: &ctx, CancelFunc: cancelFunc, logs: strings.Builder{}, Trace: xid.New().String()}
	}
	ctx, cancelFunc := context.WithCancel(background)
	return &Context{Ctx: &ctx, CancelFunc: cancelFunc, logs: strings.Builder{}, Trace: xid.New().String()}
}

func (c *Context) Info(args ...interface{}) {
	zlogger.Info(args)
}

func (c *Context) Infof(msg string, args ...interface{}) {
	zlogger.Infof(msg, args)
}

func (c *Context) Error(args ...interface{}) {
	zlogger.Error(args)
}

func (c *Context) Errorf(msg string, args ...interface{}) {
	zlogger.Errorf(msg, args)
}

// Console Log 调度中心控制台输出，请勿打印敏感信息
func (c *Context) Console(msg string, args ...interface{}) {
	if len(args) == 0 {
		c.logs.WriteString(msg)
	} else {
		c.logs.WriteString(fmt.Sprintf(msg, args...))
	}
	// 换行
	c.logs.WriteString("<br>")
	// 文件同步到日志文件
	c.Infof(msg, args...)
}

// ConsoleErr 记录错误日志信息
func (c *Context) ConsoleErr(msg string, args ...interface{}) {
	c.logs.WriteString("<text style='color:red'>")
	if len(args) == 0 {
		c.logs.WriteString(msg)
	} else {
		c.logs.WriteString(fmt.Sprintf(msg, args...))
	}
	c.logs.WriteString("</text>")
	// 换行
	c.logs.WriteString("<br>")
	// 文件同步到日志文件
	c.Errorf(msg, args...)
}

// GetConsoleLog 获取调度中心控制台输出
func (c *Context) GetConsoleLog() string {
	return c.logs.String()
}

func funcName(trace string) string {
	pc, _, _, _ := runtime.Caller(2)
	functionName := runtime.FuncForPC(pc).Name()
	if trace == "" {
		return fmt.Sprintf(path.Base(functionName))
	}
	return fmt.Sprintf("%s  trace:%s", path.Base(functionName), trace)
}
