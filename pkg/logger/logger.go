package logger

import (
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interfeysi
type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) Logger
	Sync() error
}

// loggerImpl implementatsiyasi
type loggerImpl struct {
	zap *zap.Logger
}

func (l *loggerImpl) Debug(msg string, fields ...zapcore.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *loggerImpl) Info(msg string, fields ...zapcore.Field) {
	l.zap.Info(msg, fields...)
}

func (l *loggerImpl) Warn(msg string, fields ...zapcore.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *loggerImpl) Error(msg string, fields ...zapcore.Field) {
	l.zap.Error(msg, fields...)
}

func (l *loggerImpl) Fatal(msg string, fields ...zapcore.Field) {
	l.zap.Fatal(msg, fields...)
}

func (l *loggerImpl) With(fields ...zapcore.Field) Logger {
	return &loggerImpl{zap: l.zap.With(fields...)}
}

func (l *loggerImpl) Sync() error {
	return l.zap.Sync()
}

// productionConfig funksiyasi
func productionConfig(file string) zap.Config {
	configZap := zap.NewProductionConfig()
	configZap.OutputPaths = []string{"stdout", file}
	configZap.DisableStacktrace = true
	return configZap
}

// developmentConfig funksiyasi
func developmentConfig(file string) zap.Config {
	configZap := zap.NewDevelopmentConfig()
	configZap.OutputPaths = []string{"stdout", file}
	configZap.ErrorOutputPaths = []string{"stderr"}
	return configZap
}

// New funksiyasi - Logger interfeysini yaratadi
func New(level, environment, file_name string) (Logger, error) {
	file := filepath.Join("./" + file_name)

	configZap := productionConfig(file)

	if environment == "development" {
		configZap = developmentConfig(file)
	}

	switch level {
	case "debug":
		configZap.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		configZap.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		configZap.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		configZap.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "dpanic":
		configZap.Level = zap.NewAtomicLevelAt(zap.DPanicLevel)
	case "panic":
		configZap.Level = zap.NewAtomicLevelAt(zap.PanicLevel)
	case "fatal":
		configZap.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		configZap.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	zapLogger, err := configZap.Build()
	if err != nil {
		return nil, err
	}

	return &loggerImpl{zap: zapLogger}, nil
}

// Error funksiyasi
func Error(err error) zapcore.Field {
	return zap.Error(err)
}
