package rcalc

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func InitDevLogger(filePath string) {

	configEncoder := zap.NewProductionEncoderConfig()
	configEncoder.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewConsoleEncoder(configEncoder)

	var writer zapcore.WriteSyncer
	if filePath == "-" {
		writer = zapcore.AddSync(os.Stdout)
	} else {
		logFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		writer = zapcore.AddSync(logFile)
	}
	defaultLogLevel := zapcore.DebugLevel
	baseLogger := zap.New(zapcore.NewCore(fileEncoder, writer, defaultLogLevel), zap.AddStacktrace(zapcore.ErrorLevel))
	logger = baseLogger.Sugar()
}

func InitProdLogger(filePath string) {

	configEncoder := zap.NewProductionEncoderConfig()
	configEncoder.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewConsoleEncoder(configEncoder)
	logFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.InfoLevel
	baseLogger := zap.New(zapcore.NewCore(fileEncoder, writer, defaultLogLevel), zap.AddStacktrace(zapcore.ErrorLevel))
	logger = baseLogger.Sugar()
}

func InitProdStdOutLogger() {
	configEncoder := zap.NewProductionEncoderConfig()
	configEncoder.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewConsoleEncoder(configEncoder)

	defaultLogLevel := zapcore.InfoLevel
	baseLogger := zap.New(zapcore.NewCore(fileEncoder, os.Stdout, defaultLogLevel), zap.AddStacktrace(zapcore.ErrorLevel))
	logger = baseLogger.Sugar()

}

func GetLogger() *zap.SugaredLogger {
	return logger
}
