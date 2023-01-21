package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/use_case"
)

var servicesSet = wire.NewSet(
	//repo2.NewOfficeRepository,
	use_case.NewUserService,
	//repo2.NewPlaceRepository,
)
