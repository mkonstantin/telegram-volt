package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/form"
	form2 "telegram-api/internal/app/informer/form"
)

var formSet = wire.NewSet(
	form.NewOfficeListForm,
	form.NewOfficeMenuForm,
	form.NewSeatListForm,
	form.NewFreeSeatForm,
	form.NewOwnSeatForm,
	form.NewDateMenutForm,
	form2.NewInfoMenuForm,
)
