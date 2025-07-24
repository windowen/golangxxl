package zlogger

import (
	"fmt"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"liveJob/pkg/common/config"
	"liveJob/pkg/safego"
)

var (
	ZLog *zap.Logger

	// 逻辑相关 内部日志接口
	ZLogInner *zap.Logger
	SLogInner *zap.SugaredLogger

	// 外部追加属性
	AppendOptions []zap.Option
)

type logConfig struct {
	// 日志等级. debug/info/warn/error
	Level string
	// 日志追踪等级
	StackTrace string
	// 输出
	Output string
	// 是否添加caller
	Caller bool
	// 配置读取前缀
	prefix string
}

func newLogConfig() *logConfig {
	cfg := &logConfig{
		Level:      config.Config.Log.LogLevel,
		StackTrace: "panic",
		Output:     "stdout",
		Caller:     true,
	}
	return cfg
}

var logFileName string
var alertCore zapcore.Core
var globalFields map[string]string

func InitLogConfig(logName string) {
	// 导入落地功能
	// lumberjack.InitLumberjackLogger()
	// 默认配置设置
	logFileName = logName
	// 日志配置
	var logicConfig = newLogConfig()
	logicConfig.Output = logFileName

	if err := initLog(logicConfig); err != nil {
		return
	}
}

func getLogLevel(lvl string) zap.AtomicLevel {
	lv := zap.NewAtomicLevel()
	switch lvl {
	case "panic":
		lv.SetLevel(zap.PanicLevel)
	case "fatal":
		lv.SetLevel(zap.FatalLevel)
	case "error":
		lv.SetLevel(zap.ErrorLevel)
	case "info":
		lv.SetLevel(zap.InfoLevel)
	case "debug", "trace":
		lv.SetLevel(zap.DebugLevel)
	case "warn":
		lv.SetLevel(zap.WarnLevel)
	}
	return lv
}

func SetEmptyLogger() {
	ZLogInner = zap.NewNop()
	SLogInner = ZLogInner.Sugar()
	ZLog = ZLogInner
}

// 新建日志接口
func newLogger(dev bool, logCfg *logConfig) (*zap.Logger, error) {
	var cfg zap.Config
	// if dev {
	// 	cfg = zap.NewDevelopmentConfig()
	// } else {
	cfg = zap.NewProductionConfig()
	// }
	cfg.Level = getLogLevel(logCfg.Level)
	cfg.DisableStacktrace = true
	var opts []zap.Option
	if logCfg.StackTrace != "" {
		opts = append(opts, zap.AddStacktrace(getLogLevel(logCfg.StackTrace).Level()))
	}
	if logCfg.Caller {
		opts = append(opts, zap.AddCaller())
	}
	if logFileName != "" {
		logCfg.Output = logFileName
	}
	cfg.OutputPaths = []string{logCfg.Output, "stdout"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05]")
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	if alertCore != nil {
		opts = append(opts, zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, alertCore)
		}))
	}
	// 追加属性
	if len(AppendOptions) > 0 {
		opts = append(opts, AppendOptions...)
	}

	return cfg.Build(opts...)
}

func newRotateLogger(logCfg *logConfig) (*zap.Logger, error) {
	if logFileName == "" {
		panic("log file name is empty")
	}

	var cfg zap.Config
	cfg = zap.NewProductionConfig()
	cfg.Level = getLogLevel(logCfg.Level)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	var opts []zap.Option
	if logCfg.StackTrace != "" {
		opts = append(opts, zap.AddStacktrace(getLogLevel(logCfg.StackTrace).Level()))
	}
	if logCfg.Caller {
		opts = append(opts, zap.AddCaller())
	}

	linkName := fmt.Sprintf("%s.log", logFileName)
	rotateLogger, err := rotatelogs.New(
		fmt.Sprintf("%s-%s.log", logFileName, "%Y-%m-%d"), // 日志文件名带时间格式
		rotatelogs.WithLinkName(linkName),                 // 创建符号链接
		rotatelogs.WithMaxAge(7*24*time.Hour),             // 保留一周的日志
		rotatelogs.WithRotationTime(24*time.Hour),         // 每 24 小时切分
		rotatelogs.WithRotationSize(100*1024*1024),        // 100M大小切分文件
	)
	if err != nil {
		return nil, err
	}
	fileSyncer := zapcore.AddSync(rotateLogger)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig), // 可以用 NewConsoleEncoder 切换到普通文本格式
		fileSyncer,
		cfg.Level,
	)
	return zap.New(core, opts...), nil
}

// 从配置文件加载Log配置，并初始化Log Hook
func initLog(logicLogCfg *logConfig) error {
	fmt.Println("init logger")

	// 生成日志接口
	logic, err := newRotateLogger(logicLogCfg)
	if err != nil {
		return err
	}

	// 本地日志接口
	ZLog = logic
	SLogInner = logic.WithOptions(zap.AddCallerSkip(1)).Sugar()
	ZLogInner = logic.WithOptions(zap.AddCallerSkip(1))

	safego.PanicCatchFunc = func(name string, p interface{}) {
		ZLog.Error("receive panic", zap.String("name", name), zap.Any("p", p))
	}

	if globalFields != nil {
		setGlobalFields(globalFields)
	}

	return nil
}

func SetGlobalFields(fields map[string]string) {
	globalFields = fields
}

func setGlobalFields(fields map[string]string) {
	nf := make(map[string]interface{})
	for k, v := range fields {
		nf[k] = v
	}
	zf := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zf = append(zf, zap.String(k, v))
	}

	ZLog = ZLog.With(zf...)
	SLogInner = ZLog.WithOptions(zap.AddCallerSkip(1)).Sugar()
	ZLogInner = ZLog.WithOptions(zap.AddCallerSkip(1))
}

func Fatal(format ...interface{}) {
	SLogInner.Fatal(format...)
}
func Debug(args ...interface{}) {
	SLogInner.Debug(args...)
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
func Panic(args ...interface{}) {
	SLogInner.Panic(args...)
}
func Warnf(format string, args ...interface{}) {
	SLogInner.Warnf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	SLogInner.Panicf(format, args...)
}
func Infof(format string, args ...interface{}) {
	SLogInner.Infof(format, args...)
}
func Errorf(format string, args ...interface{}) {
	SLogInner.Errorf(format, args...)
}
func Debugf(format string, args ...interface{}) {
	SLogInner.Debugf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	SLogInner.Fatalf(format, args...)
}
func Fatalln(args ...interface{}) {
	SLogInner.Fatal(args...)
}
func Debugln(args ...interface{}) {
	SLogInner.Debug(args...)
}
func Errorln(args ...interface{}) {
	SLogInner.Error(args...)
}
func Infoln(args ...interface{}) {
	SLogInner.Info(args...)
}
func Warnln(args ...interface{}) {
	SLogInner.Warn(args...)
}
func Panicln(args ...interface{}) {
	SLogInner.Panic(args...)
}

func Fatalw(msg string, fields ...zap.Field) {
	ZLogInner.Fatal(msg, fields...)
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
func Panicw(msg string, fields ...zap.Field) {
	ZLogInner.Panic(msg, fields...)
}
