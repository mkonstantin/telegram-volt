package main

import (
	"github.com/google/wire"
	handler2 "telegram-api/internal/app/handler"
)

var handlerSet = wire.NewSet(
	handler2.NewStartHandle,
	handler2.NewOfficeListHandle,
	handler2.NewOfficeMenuHandle,
	handler2.NewSeatListHandle,
	handler2.NewOwnSeatMenuHandle,
	handler2.NewFreeSeatMenuHandle,
)
