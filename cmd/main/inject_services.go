package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/usecase"
)

var servicesSet = wire.NewSet(
	usecase.NewOfficeService,
	usecase.NewUserService,
	//repo2.NewPlaceRepository,
)
