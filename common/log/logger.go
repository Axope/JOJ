package log

import (
	"github.com/Axope/JOJ/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
)

var Any = zap.Any

var Logger *zap.Logger
var LoggerSugar *zap.SugaredLogger

func InitLogger() {
	LogConfig := configs.GetLogConfig()
	// 日志分割
	hook := lumberjack.Logger{
		Filename:   LogConfig.Path,
		MaxSize:    LogConfig.MaxSize,
		MaxBackups: LogConfig.MaxBackups,
		MaxAge:     LogConfig.MaxAge,
		Compress:   LogConfig.Compress,
	}
	write := zapcore.AddSync(&hook)

	// 日志等级，默认INFO
	logLevel := strings.ToUpper(LogConfig.Level)
	var level zapcore.Level
	switch logLevel {
	case "DEBUG":
		level = zap.DebugLevel
	case "INFO":
		level = zap.InfoLevel
	case "WARN":
		level = zap.WarnLevel
	case "ERROR":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		write,
		level,
	)
	caller := zap.AddCaller()
	development := zap.Development()

	Logger = zap.New(core, caller, development)
	LoggerSugar = Logger.Sugar()
	Logger.Info("Log module init success")
}
