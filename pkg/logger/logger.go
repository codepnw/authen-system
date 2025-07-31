package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init() (*zap.Logger, error) {
	// Create logs directory if exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		_ = os.Mkdir("logs", os.ModePerm)
	}

	logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// encoder (JSON format or console)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	// Output File
	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, zap.InfoLevel),
		zapcore.NewCore(encoder, consoleWriter, zap.InfoLevel),
	)

	log = zap.New(core, zap.AddCaller())
	return log, nil
}

func Info(code, msg string, data any) {
	log.Info(msg, zap.String("code", code), zap.Any("data", data))
}

func InfoMiddleware(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Error(code, msg string, err error) {
	log.Error(msg, zap.String("code", code), zap.Error(err))
}

func Logger() *zap.Logger {
	return log
}
