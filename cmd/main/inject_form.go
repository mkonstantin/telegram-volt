package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/form"
)

var formSet = wire.NewSet(
	form.NewOfficeMenuForm,
	form.NewOfficeListMenuForm,
)
