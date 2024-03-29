package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/menu"
)

var menuSet = wire.NewSet(
	menu.NewOfficeListMenu,
	menu.NewOfficeMenu,
	menu.NewSeatListMenu,
	menu.NewFreeSeatMenu,
	menu.NewOwnSeatMenu,
	menu.NewDateMenu,
	menu.NewHoldSeatMenu,
)
