package main

import (
	"github.com/google/wire"
	"telegram-api/internal/infrastructure/handler"
)

var handlerSet = wire.NewSet(
	handler.NewStartHandle,
	handler.NewOfficeListHandle,
	handler.NewOfficeMenuHandle,
	handler.NewSeatListHandle,
	handler.NewOwnSeatMenuHandle,
	handler.NewFreeSeatMenuHandle,
)
