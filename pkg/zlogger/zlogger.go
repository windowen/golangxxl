package zlogger

import (
	"fmt"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"queueJob/pkg/common/config"
	"queueJob/pkg/safego"
)

var (
	ZLog         *zap.Logger
	ZLogInner    *zap.Logger
	SLogInner    *zap.SugaredLogger
	logFileName  string
	globalFields map[string]string
)

type logConfig struct {
	Level      string // 日志等级
	StackTrace string // 堆栈追踪等级
	Output     string // 输出位置
	Caller     bool   // 是否输出调用者
	prefix     string // 配置前缀
}

func newLogConfig() *logConfig {
	return &logConfig{
		Level:      config.Config.Log.LogLevel,
		StackTrace: "panic",
		Output:     "stdout",
		Caller:     true,
	}
}

// InitLogConfig 初始化日志配置
func InitLogConfig(logName string) {
	logFileName = logName
	logicConfig := newLogConfig()
	logicConfig.Output = logFileName

	if err := initLog(logicConfig); err != nil {
		fmt.Printf("init log error: %v\n", err)
	}
}

func getLogLevel(level string) zap.AtomicLevel {
	lv := zap.NewAtomicLevel()

	switch level {
	case "panic":
		lv.SetLevel(zap.PanicLevel)
	case "fatal":
		lv.SetLevel(zap.FatalLevel)
	case "error":
		lv.SetLevel(zap.ErrorLevel)
	case "warn":
		lv.SetLevel(zap.WarnLevel)
	case "info":
		lv.SetLevel(zap.InfoLevel)
	case "debug", "trace":
		lv.SetLevel(zap.DebugLevel)
	default:
		lv.SetLevel(zap.InfoLevel) // 默认info
	}

	return lv
}

func newRotateLogger(cfg *logConfig) (*zap.Logger, error) {
	if logFileName == "" {
		panic("log file name is empty")
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = getLogLevel(cfg.Level)
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapCfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	var opts []zap.Option
	if cfg.StackTrace != "" {
		opts = append(opts, zap.AddStacktrace(getLogLevel(cfg.StackTrace).Level()))
	}
	if cfg.Caller {
		opts = append(opts, zap.AddCaller())
	}

	linkName := fmt.Sprintf("%s.log", logFileName)
	rotateLogger, err := rotatelogs.New(
		fmt.Sprintf("%s-%s.log", logFileName, "%Y-%m-%d"), // 文件名按时间分割
		rotatelogs.WithLinkName(linkName),                 // 创建软链接
		rotatelogs.WithMaxAge(7*24*time.Hour),             // 保留7天
		rotatelogs.WithRotationTime(24*time.Hour),         // 每24小时切割
		rotatelogs.WithRotationSize(100*1024*1024),        // 每100MB切割（注意：rotatelogs版本不同可能不支持此选项）
	)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapCfg.EncoderConfig),
		zapcore.AddSync(rotateLogger),
		zapCfg.Level,
	)

	return zap.New(core, opts...), nil
}

func initLog(cfg *logConfig) error {
	fmt.Println("init logger")

	logger, err := newRotateLogger(cfg)
	if err != nil {
		return err
	}

	ZLog = logger
	SLogInner = logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
	ZLogInner = logger.WithOptions(zap.AddCallerSkip(1))

	safego.PanicCatchFunc = func(name string, p interface{}) {
		ZLog.Error("receive panic", zap.String("name", name), zap.Any("p", p))
	}

	if globalFields != nil {
		setGlobalFields(globalFields)
	}

	return nil
}

func setGlobalFields(fields map[string]string) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.String(k, v))
	}

	ZLog = ZLog.With(zapFields...)
	SLogInner = ZLog.WithOptions(zap.AddCallerSkip(1)).Sugar()
	ZLogInner = ZLog.WithOptions(zap.AddCallerSkip(1))
}

func Error(args ...interface{}) {
	SLogInner.Error(args...)
}

func Info(args ...interface{}) {
	SLogInner.Info(args...)
}

func Warn(args ...interface{}) {
	SLogInner.Warn(args...)
}

func Errorf(format string, args ...interface{}) {
	SLogInner.Errorf(format, args...)
}

func Infof(format string, args ...interface{}) {
	SLogInner.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	SLogInner.Warnf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	SLogInner.Debugf(format, args...)
}

func Debugw(msg string, fields ...zap.Field) {
	ZLogInner.Debug(msg, fields...)
}

func Errorw(msg string, fields ...zap.Field) {
	ZLogInner.Error(msg, fields...)
}

func Infow(msg string, fields ...zap.Field) {
	ZLogInner.Info(msg, fields...)
}

func Warnw(msg string, fields ...zap.Field) {
	ZLogInner.Warn(msg, fields...)
}
