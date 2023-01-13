package rcalc

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.SugaredLogger

func InitDevLogger(filePath string) {

	configEncoder := zap.NewProductionEncoderConfig()
	configEncoder.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewConsoleEncoder(configEncoder)
	logFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	baseLogger := zap.New(zapcore.NewCore(fileEncoder, writer, defaultLogLevel), zap.AddStacktrace(zapcore.ErrorLevel))
	logger = baseLogger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return logger
}
