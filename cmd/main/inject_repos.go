package main

import (
	"github.com/google/wire"
	"telegram-api/internal/infrastructure/repo"
)

var repositorySet = wire.NewSet(
	repo.NewOfficeRepository,
	repo.NewUserRepository,
	repo.NewSeatRepository,
	repo.NewBookSeatRepository,
	repo.NewWorkDateRepository,
)
