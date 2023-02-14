package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/handler/dto"
)

type OfficeListHandle interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type officeListHandleImpl struct {
	userService usecase.UserService
	msgFormer   MessageFormer
	logger      *zap.Logger
}

func NewOfficeListHandle(userService usecase.UserService, logger *zap.Logger) OfficeListHandle {
	return &officeListHandleImpl{
		userService: userService,
		logger:      logger,
	}
}

func (o *officeListHandleImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	result, err := o.userService.SetOfficeScript(ctx, command.OfficeID)
	if err != nil {
		return nil, err
	}

	return o.msgFormer.FormOfficeMenuMsg(ctx, result)
}
