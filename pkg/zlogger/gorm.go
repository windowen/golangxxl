package zlogger

import (
	ctx "context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"liveJob/pkg/common/config"
)

type dbLog struct {
	logger.Config
	traceSwitch bool
}

func NewDBLog(config logger.Config) *dbLog {
	return &dbLog{Config: config, traceSwitch: true}
}

func (l *dbLog) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	newLogger.traceSwitch = true
	return &newLogger
}

func (l *dbLog) Info(ctx ctx.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Info {
		Info(fmt.Sprintf("DBLOG |"+s+"\n", i...))
	}
}

func (l *dbLog) Warn(ctx ctx.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Warn {
		Warn(fmt.Sprintf("DBLOG |"+s+"\n", i...))
	}
}

func (l *dbLog) Error(ctx ctx.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Error {
		Error(fmt.Sprintf("DBLOG |"+s+"\n", i...))
	}
}

func (l *dbLog) Trace(ctx ctx.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent || !l.traceSwitch {
		return
	}

	// 耗时
	timestamp := time.Since(begin)

	sql, rows := fc()

	// 如果报错直接打印日志
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		Errorf("DBLOG |err=%v |sql=%s |rows=%d |timestamp=%.2fs", err, sql, rows, timestamp.Seconds())
		return
	}

	// 产生慢查询
	if timestamp > l.SlowThreshold {
		// 记录不同标准的慢sql
		var standard = "PURPLE"
		switch {
		case timestamp > l.SlowThreshold*8:
			standard = "RED"
		case timestamp > l.SlowThreshold*4:
			standard = "ORANGE"
		case timestamp > l.SlowThreshold*2:
			standard = "YELLOW"
		}

		if timestamp > l.SlowThreshold*4 {
			Warnf("DBLOG |SLOW |standard=%s |sql=%s |rows=%d |timestamp=%.2fs", standard, sql, rows, timestamp.Seconds())
			return
		}

		Warnf("DBLOG |SLOW |standard=%s |sql=%s |rows=%d |timestamp=%.2fs", standard, sql, rows, timestamp.Seconds())

		return
	}

	// 如果不是慢日志，只有数据变更操作的日志才记录
	if rows > 0 && (strings.Contains(sql, "INSERT") || strings.Contains(sql, "UPDATE")) {
		Infof("DBLOG |sql=%s |rows=%d |timestamp=%.2fs", sql, rows, timestamp.Seconds())
		return
	}

	// 更新0的日志额外重点记录
	if strings.Contains(sql, "UPDATE") {
		Warnf("DBLOG | sql=%s |rows=%d |timestamp=%.2fs", sql, rows, timestamp.Seconds())
		return
	}

	// 开发环境打开SQL日志，方便查看SQL
	if strings.ToLower(config.Config.App.Env) == "dev" {
		Infof("DBLOG |sql=%s |timestamp=%.2fs", sql, timestamp.Seconds())
		return
	}
}
