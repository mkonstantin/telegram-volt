package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/usecase"
)

type DateMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type dateMenuImpl struct {
	userService usecase.UserService
	officeMenu  interfaces.OfficeMenu
	logger      *zap.Logger
}

func NewDateMenuHandle(
	userService usecase.UserService,
	officeMenu interfaces.OfficeMenu,
	logger *zap.Logger) DateMenu {

	return &dateMenuImpl{
		userService: userService,
		officeMenu:  officeMenu,
		logger:      logger,
	}
}

func (o *dateMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	if command.BookDate == nil {
		return o.officeMenu.Call(ctx)
	}

	return nil, nil
}
