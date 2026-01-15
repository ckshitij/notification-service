package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	log *zap.Logger
}

func NewZapLogger(env string, level int) (Logger, error) {
	var cfg zap.Config

	if env == "local" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(level))
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	zl, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return nil, err
	}

	return &zapLogger{log: zl}, nil
}

func (z *zapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	z.log.Debug(msg, z.toZapFields(ctx, fields)...)
}

func (z *zapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	z.log.Info(msg, z.toZapFields(ctx, fields)...)
}

func (z *zapLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	z.log.Warn(msg, z.toZapFields(ctx, fields)...)
}

func (z *zapLogger) Error(ctx context.Context, msg string, fields ...Field) {
	z.log.Error(msg, z.toZapFields(ctx, fields)...)
}

func (z *zapLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	z.log.Fatal(msg, z.toZapFields(ctx, fields)...)
	os.Exit(1)
}

func (z *zapLogger) toZapFields(ctx context.Context, fields []Field) []zap.Field {
	all := extractContextFields(ctx)

	all = append(all, fields...)

	zapFields := make([]zap.Field, 0, len(all))
	for _, f := range all {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
