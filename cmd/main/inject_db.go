package main

import (
	"context"
	"github.com/google/wire"
	"go.uber.org/zap"
	repository "telegram-api/pkg"
)

var dbSet = wire.NewSet(
	provideDBConnection,
	context.Background,
	repository.NewDB,
)

func provideDBConnection(ctx context.Context, l *zap.Logger) (repository.Connection, func()) {
	con, cleanup, err := repository.InitConnection(ctx, l)
	if err != nil {
		l.Error("Can't Init connection", zap.Error(err))
	}
	return con, cleanup
}
