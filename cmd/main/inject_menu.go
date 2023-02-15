package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/menu"
)

var menuSet = wire.NewSet(
	menu.NewOfficeMenu,
	menu.NewOfficeListMenu,
)
