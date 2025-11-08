package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New cria um logger configurado para ambientes de API.
func New(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if level != "" {
		if err := cfg.Level.UnmarshalText([]byte(level)); err != nil {
			return nil, err
		}
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return cfg.Build()
}
