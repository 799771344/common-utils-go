package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"strings"
	"time"
)

func NewLogger(level string, filePath string) {
	// 创建文本编码器
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		ConsoleSeparator: "	",
	}
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"))
	}
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	file, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	sink := zapcore.AddSync(file)

	zapLevel := zap.DebugLevel
	switch level {
	case "DEBUG":
		zapLevel = zap.DebugLevel
	case "INFO":
		zapLevel = zap.InfoLevel
	case "ERROR":
		zapLevel = zap.ErrorLevel
	case "WARN":
		zapLevel = zap.WarnLevel
	case "DPANIC":
		zapLevel = zap.DPanicLevel
	case "PANIC":
		zapLevel = zap.PanicLevel
	case "FATAL":
		zapLevel = zap.FatalLevel
	}

	Logger = new(zap.Logger)

	// 创建Logger
	Logger = zap.New(zapcore.NewCore(consoleEncoder, zapcore.Lock(sink), zap.NewAtomicLevelAt(zapLevel)))
	defer Logger.Sync()
	Logger = Logger.WithOptions(zap.AddCaller())

}

func Debug(msg string, arr ...interface{}) {
	msgNew := fmt.Sprintf(msg, arr...)
	location := caller()
	msgNew = fmt.Sprintf("%s	%s", location, msg)
	Logger.Debug(msgNew)
}
func Info(msg string, arr ...interface{}) {
	msgNew := fmt.Sprintf(msg, arr...)
	location := caller()
	msgNew = fmt.Sprintf("%s	%s", location, msg)
	Logger.Info(msgNew)
}
func Error(msg string, arr ...interface{}) {
	msgNew := fmt.Sprintf(msg, arr...)
	location := caller()
	msgNew = fmt.Sprintf("%s	%s", location, msg)
	Logger.Error(msgNew)
}
func Warn(msg string, arr ...interface{}) {
	msgNew := fmt.Sprintf(msg, arr...)
	location := caller()
	msgNew = fmt.Sprintf("%s	%s", location, msg)
	Logger.Warn(msgNew)
}
func Panic(msg string, arr ...interface{}) {
	msgNew := fmt.Sprintf(msg, arr...)
	location := caller()
	msgNew = fmt.Sprintf("%s	%s", location, msg)
	Logger.Panic(msgNew)
}
func DPanic(msg string, arr ...interface{}) {
	msgNew := fmt.Sprintf(msg, arr...)
	location := caller()
	msgNew = fmt.Sprintf("%s	%s", location, msg)
	Logger.DPanic(msgNew)
}
func Fatal(msg string, arr ...interface{}) {
	msgNew := fmt.Sprintf(msg, arr...)
	location := caller()
	msgNew = fmt.Sprintf("%s	%s", location, msg)
	Logger.Fatal(msgNew)
}
func caller() string {
	_, file, line, _ := runtime.Caller(2)
	fileP := strings.Split(file, "/")
	fileP1 := fileP[len(fileP)-2:]
	result := fmt.Sprintf("%s:%d", strings.Join(fileP1, "/"), line)
	return result
}
