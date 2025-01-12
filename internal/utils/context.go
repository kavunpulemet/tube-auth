package utils

import (
	"context"

	"go.uber.org/zap"
)

type MyContext struct {
	Ctx    context.Context
	Logger *zap.SugaredLogger
}

func NewMyContext(ctx context.Context, logger *zap.SugaredLogger) MyContext {
	return MyContext{
		Ctx:    ctx,
		Logger: logger,
	}
}
