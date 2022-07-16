package main

import (
	"github.com/google/wire"
	"telegram-api/internal/repo"
)

var repositorySet = wire.NewSet(
	repo.NewOfficeRepository,
	repo.NewUserRepository,
	repo.NewPlaceRepository,
)
