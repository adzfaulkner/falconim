package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type wrapperLogger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Sync() error
}

type Handler struct {
	wrapLogger wrapperLogger
}

func (h *Handler) Info(msg string, fields ...zap.Field) {
	h.wrapLogger.Info(msg, fields...)
}

func (h *Handler) Error(msg string, fields ...zap.Field) {
	h.wrapLogger.Error(msg, fields...)
}

func (h *Handler) Sync() error {
	if err := h.wrapLogger.Sync(); err != nil {
		return fmt.Errorf("error whilst sync logger: %w", err)
	}

	return nil
}

func NewHandler() (*Handler, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, fmt.Errorf("could not create zap logger: %w", err)
	}

	return &Handler{
		wrapLogger: logger,
	}, nil
}
