package zlogger

import (
	ctx "context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbLog struct {
	logger.Config
	traceSwitch bool
}

func NewDBLog(config logger.Config) *DbLog {
	return &DbLog{
		Config:      config,
		traceSwitch: true,
	}
}

func (l *DbLog) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	newLogger.traceSwitch = true
	return &newLogger
}

func (l *DbLog) Info(ctx ctx.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		Infof(msg, args...)
	}
}

func (l *DbLog) Warn(ctx ctx.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		Warnf(msg, args...)
	}
}

func (l *DbLog) Error(ctx ctx.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		Errorf(msg, args...)
	}
}

func (l *DbLog) Trace(ctx ctx.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent || !l.traceSwitch {
		return
	}

	// 计算执行耗时
	duration := time.Since(begin)

	sql, rows := fc()

	// 如果出现错误且不是记录未找到，直接记录错误日志
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		Errorw("trace",
			zap.String("pos", funcName4Gorm()),
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("timestamp", duration.Seconds()),
		)
		return
	}

	// 慢查询日志处理
	if duration > l.SlowThreshold {
		standard := "PURPLE"
		switch {
		case duration > l.SlowThreshold*8:
			standard = "RED"
		case duration > l.SlowThreshold*4:
			standard = "ORANGE"
		case duration > l.SlowThreshold*2:
			standard = "YELLOW"
		}

		// 对慢查询分级别记录 WARN 日志
		if duration > l.SlowThreshold*4 {
			Warnw("trace",
				zap.String("pos", funcName4Gorm()),
				zap.String("standard", standard),
				zap.String("sql", sql),
				zap.Int64("rows", rows),
				zap.Float64("timestamp", duration.Seconds()),
			)
			return
		}

		Warnw("trace slow",
			zap.String("pos", funcName4Gorm()),
			zap.String("standard", standard),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("timestamp", duration.Seconds()),
		)
		return
	}

	// 非慢查询，只有数据变更操作才记录日志
	if rows > 0 && (strings.Contains(sql, "INSERT") || strings.Contains(sql, "UPDATE")) {
		Infow("trace",
			zap.String("pos", funcName4Gorm()),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("timestamp", duration.Seconds()),
		)
		return
	}

	// 对更新但影响0行的日志重点记录
	if strings.Contains(sql, "UPDATE") {
		Warnw("trace",
			zap.String("pos", funcName4Gorm()),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("timestamp", duration.Seconds()),
		)
		return
	}

	// 其他情况记录普通信息日志
	Infow("trace",
		zap.String("pos", funcName4Gorm()),
		zap.String("sql", sql),
		zap.Float64("timestamp", duration.Seconds()),
	)
}

func funcName4Gorm() string {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		return "unknown:0"
	}

	file = filepath.ToSlash(file)

	return fmt.Sprintf("%s:%d", file, line)
}
