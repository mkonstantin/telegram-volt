package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/former"
	"telegram-api/internal/infrastructure/handler/dto"
)

type OfficeList interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type officeListImpl struct {
	userService usecase.UserService
	msgFormer   former.MessageFormer
	logger      *zap.Logger
}

func NewOfficeListHandle(
	userService usecase.UserService,
	msgFormer former.MessageFormer,
	logger *zap.Logger) OfficeList {

	return &officeListImpl{
		userService: userService,
		msgFormer:   msgFormer,
		logger:      logger,
	}
}

func (o *officeListImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	result, err := o.userService.SetOfficeScript(ctx, command.OfficeID)
	if err != nil {
		return nil, err
	}

	return o.msgFormer.FormOfficeMenuMsg(ctx, result)
}
