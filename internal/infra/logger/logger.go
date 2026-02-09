package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.Logger
}

func New() (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "@timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	z, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{log: z}, nil
}

func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()

	if !sc.IsValid() {
		return l.log
	}

	return l.log.With(
		zap.String("trace_id", sc.TraceID().String()),
		zap.String("span_id", sc.SpanID().String()),
	)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	l.WithContext(ctx).Error(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Warn(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	l.WithContext(ctx).Fatal(msg, fields...)
}
