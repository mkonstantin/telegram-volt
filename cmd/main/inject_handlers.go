package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/handler"
)

var handlerSet = wire.NewSet(
	handler.NewStartHandle,
	handler.NewOfficeListHandle,
	handler.NewOfficeMenuHandle,
	handler.NewSeatListHandle,
	handler.NewOwnSeatMenuHandle,
	handler.NewFreeSeatMenuHandle,
	handler.NewDateMenuHandle,
)
