package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
)

type Start interface {
	Handle(ctx context.Context) (*tgbotapi.MessageConfig, error)
}

type startImpl struct {
	officeMenu     interfaces.OfficeMenu
	officeListMenu interfaces.OfficeListMenu
	logger         *zap.Logger
}

func NewStartHandle(
	officeMenu interfaces.OfficeMenu,
	officeListMenu interfaces.OfficeListMenu,
	logger *zap.Logger) Start {

	return &startImpl{
		officeMenu:     officeMenu,
		officeListMenu: officeListMenu,
		logger:         logger,
	}
}

func (s *startImpl) Handle(ctx context.Context) (*tgbotapi.MessageConfig, error) {

	currentUser := model.GetCurrentUser(ctx)

	if currentUser.HaveChosenOffice() {
		return s.officeMenu.Call(ctx)
	} else {
		return s.officeListMenu.Call(ctx)
	}
}
