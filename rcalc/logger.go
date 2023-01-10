package rcalc

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.Logger

func InitDevLogger(filePath string) {

	configEncoder := zap.NewProductionEncoderConfig()
	configEncoder.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewConsoleEncoder(configEncoder)
	logFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	logger = zap.New(zapcore.NewCore(fileEncoder, writer, defaultLogLevel), zap.AddStacktrace(zapcore.ErrorLevel))
}

func GetLogger() *zap.Logger {
	return logger
}
