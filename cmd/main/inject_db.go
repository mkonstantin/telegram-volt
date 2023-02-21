package main

import (
	"context"
	"github.com/google/wire"
	"go.uber.org/zap"
	"telegram-api/config"
	repository "telegram-api/pkg"
)

var dbSet = wire.NewSet(
	provideDBConnection,
	context.Background,
	repository.NewDB,
)

func provideDBConnection(ctx context.Context, cfg config.AppConfig, l *zap.Logger) (repository.Connection, func()) {
	con, cleanup, err := repository.InitConnection(ctx, cfg, l)
	if err != nil {
		l.Error("Can't Init connection", zap.Error(err))
	}
	return con, cleanup
}
