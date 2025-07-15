package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface with only essential methods
type Logger interface {
	Debug(msg string, tenantID string)
	Info(msg string, tenantID string)
	Warn(msg string, tenantID string)
	Error(msg string, tenantID string)
	Fatal(msg string, tenantID string)
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(env string) (*ZapLogger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		MessageKey:     "msg",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	var core zapcore.Core
	if env == "production" {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.InfoLevel,
		)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		)
	}

	return &ZapLogger{logger: zap.New(core)}, nil
}

// Debug logs a debug message
func (l *ZapLogger) Debug(msg string, tenantID string) {
	l.logger.Debug(msg, zap.String("tenant-id", tenantID))
}

// Info logs an info message
func (l *ZapLogger) Info(msg string, tenantID string) {
	l.logger.Info(msg, zap.String("tenant-id", tenantID))
}

// Warn logs a warning message
func (l *ZapLogger) Warn(msg string, tenantID string) {
	l.logger.Warn(msg, zap.String("tenant-id", tenantID))
}

// Error logs an error message
func (l *ZapLogger) Error(msg string, tenantID string) {
	l.logger.Error(msg, zap.String("tenant-id", tenantID))
}

// Fatal logs a fatal message and exits
func (l *ZapLogger) Fatal(msg string, tenantID string) {
	l.logger.Fatal(msg, zap.String("tenant-id", tenantID))
}
