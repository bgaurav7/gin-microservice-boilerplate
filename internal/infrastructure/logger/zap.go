package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper around zap.Logger
type Logger struct {
	*zap.Logger
}

// Config holds the logger configuration
type Config struct {
	Level string `mapstructure:"level"`
}

// NewLogger creates a new JSON-only logger
func NewLogger(config *Config) (*Logger, error) {
	level := zap.InfoLevel
	if config != nil {
		switch config.Level {
		case "debug":
			level = zap.DebugLevel
		case "info":
			level = zap.InfoLevel
		case "warn":
			level = zap.WarnLevel
		case "error":
			level = zap.ErrorLevel
		}
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Always use JSON encoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Write to stderr
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stderr),
		level,
	)

	// Create the logger
	zapLogger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return &Logger{
		Logger: zapLogger,
	}, nil
}

// With adds structured context to the logger
func (l *Logger) With(fields ...interface{}) *Logger {
	if len(fields)%2 != 0 {
		l.Warn("Logger.With called with odd number of parameters")
		return l
	}

	zapFields := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			l.Warn("Logger.With called with non-string key", zap.Any("key", fields[i]))
			continue
		}
		zapFields = append(zapFields, zap.Any(key, fields[i+1]))
	}

	return &Logger{
		Logger: l.Logger.With(zapFields...),
	}
}

// Debug logs a debug message with optional fields
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.Logger.Debug(msg, l.toZapFields(fields...)...)
}

// Info logs an info message with optional fields
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.Logger.Info(msg, l.toZapFields(fields...)...)
}

// Warn logs a warning message with optional fields
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.Logger.Warn(msg, l.toZapFields(fields...)...)
}

// Error logs an error message with optional fields
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.Logger.Error(msg, l.toZapFields(fields...)...)
}

// Fatal logs a fatal message with optional fields and then calls os.Exit(1)
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.Logger.Fatal(msg, l.toZapFields(fields...)...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// toZapFields converts a list of interface{} parameters to zap.Field
func (l *Logger) toZapFields(fields ...interface{}) []zap.Field {
	if len(fields)%2 != 0 {
		l.Logger.Warn("Logger called with odd number of fields", zap.Int("count", len(fields)))
		return nil
	}

	zapFields := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			l.Logger.Warn("Logger called with non-string key", zap.Any("key", fields[i]))
			continue
		}
		zapFields = append(zapFields, zap.Any(key, fields[i+1]))
	}

	return zapFields
}
