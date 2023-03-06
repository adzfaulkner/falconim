package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type logLine struct {
	level  string
	msg    string
	fields []zap.Field
}

type mockLogger struct {
	syncd bool
	line  logLine
}

func (l *mockLogger) Info(msg string, fields ...zap.Field) {
	line := logLine{
		level:  "info",
		msg:    msg,
		fields: fields,
	}

	l.line = line
}

func (l *mockLogger) Error(msg string, fields ...zap.Field) {
	line := logLine{
		level:  "error",
		msg:    msg,
		fields: fields,
	}

	l.line = line
}

func (l *mockLogger) Sync() error {
	l.syncd = true

	return nil
}

func Test_NewHandler(t *testing.T) {
	t.Parallel()

	logger := mockLogger{
		syncd: false,
	}

	sut := NewHandler(&logger)

	sut.Info(
		"info called",
		zap.String("test1", "value1"),
		zap.Int("test2", 2),
	)

	assert.Equal(t, logger.line, logLine{
		level: "info",
		msg:   "info called",
		fields: []zap.Field{
			zap.String("test1", "value1"),
			zap.Int("test2", 2),
		},
	})

	sut.Error(
		"error called",
		zap.String("test3", "value3"),
		zap.Int("test4", 4),
	)

	assert.Equal(t, logger.line, logLine{
		level: "error",
		msg:   "error called",
		fields: []zap.Field{
			zap.String("test3", "value3"),
			zap.Int("test4", 4),
		},
	})

	_ = sut.Sync()

	assert.True(t, logger.syncd)
}
