package log

import (
	"fmt"
	"go-test/internal/config"
	"go-test/internal/utils"

	"os"
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	defaultLogger Logger
	defaultConfig = &config.LogConfig{
		LogPath:    "./log",
		LogLevel:   0,
		MaxSize:    64,
		MaxAge:     31,
		MaxBackups: 10,
	}
	zapLogger *zap.Logger
)

// Logger 日志操作接口
type Logger interface {
	Debugw(msg string, keyvals ...interface{})
	Infow(msg string, keyvals ...interface{})
	Warnw(msg string, keyvals ...interface{})
	Errorw(msg string, keyvals ...interface{})
	Panicw(msg string, keyvals ...interface{})
	Fatalw(msg string, keyvals ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

// InitLogger 初始化日志
func InitLogger(config *config.LogConfig) {
	if config != nil {
		defaultConfig = config
	}
	appName := utils.GetProcessName()
	zapLogger = NewZapLogger(appName, defaultConfig)
	defaultLogger = zapLogger.Sugar()
}

// NewZapLogger 创建ZAP日志对象
func NewZapLogger(logName string, logConfig *config.LogConfig) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = CustomTimeEncoder
	encoderConfig.TimeKey = "time"

	var writeSyncer zapcore.WriteSyncer
	if logConfig.OutputFile {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s.log", logConfig.LogPath, logName),
			MaxSize:    logConfig.MaxSize,
			MaxAge:     logConfig.MaxAge,
			MaxBackups: logConfig.MaxBackups,
			Compress:   logConfig.Compress,
		}
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger))
	} else {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}

	zapOptions := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(1)}
	if logConfig.StackTrace {
		zapOptions = append(zapOptions, zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		zapOptions = append(zapOptions, zap.AddStacktrace(zapcore.FatalLevel))
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	zapCore := zapcore.NewCore(encoder, writeSyncer, ZapLevelWithLogLevel(logConfig.LogLevel))
	return zap.New(zapCore, zapOptions...)
}

// GetLogLevel 获取当前日志等级
func GetLogLevel() zapcore.Level {
	return ZapLevelWithLogLevel(defaultConfig.LogLevel)
}

// ZapLevelWithLogLevel 将日志等级映射到ZAP的日志等级
func ZapLevelWithLogLevel(level int8) zapcore.Level {
	switch level {
	case 0:
		return zapcore.DebugLevel
	case 1:
		return zapcore.InfoLevel
	case 2:
		return zapcore.WarnLevel
	case 3:
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}

// GetZapLogger 获取日志对象
func GetZapLogger(level zapcore.Level) *zap.Logger {
	checkLevel := GetLogLevel()
	if level < checkLevel {
		level = checkLevel
	}
	return zapLogger.WithOptions(zap.IncreaseLevel(level))
}

// ZapInterceptor Zap拦截器(只输出警告级及以上的日志)
func ZapInterceptor() *zap.Logger {
	level := zap.WarnLevel
	checkLevel := GetLogLevel()
	if level < checkLevel {
		level = checkLevel
	}
	logger := zapLogger.WithOptions(zap.IncreaseLevel(level))
	grpc_zap.ReplaceGrpcLoggerV2(logger)
	return logger
}

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// Debug 输出调式级日志
func Debug(msg string, keyvals ...interface{}) {
	defaultLogger.Debugw(msg, keyvals...)
}

// Info 输出提示级日志
func Info(msg string, keyvals ...interface{}) {
	defaultLogger.Infow(msg, keyvals...)
}

// Warn 输出警告级日志
func Warn(msg string, keyvals ...interface{}) {
	defaultLogger.Warnw(msg, keyvals...)
}

// Error 输出错误级日志
func Error(msg string, keyvals ...interface{}) {
	defaultLogger.Errorw(msg, keyvals...)
}

// Fatal 输出致命错误级日志
func Fatal(msg string, keyvals ...interface{}) {
	defaultLogger.Fatalw(msg, keyvals...)
}

// Debugf 输出调式级日志
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Infof 输出提示级日志
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warnf 输出警告级日志
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Errorf 输出错误级日志
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Fatalf 输出致命错误级日志
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}
