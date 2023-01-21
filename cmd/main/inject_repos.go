package main

import (
	"github.com/google/wire"
	repo2 "telegram-api/internal/infrastructure_layer/repo"
)

var repositorySet = wire.NewSet(
	//repo2.NewOfficeRepository,
	repo2.NewUserRepository,
	//repo2.NewPlaceRepository,
)
