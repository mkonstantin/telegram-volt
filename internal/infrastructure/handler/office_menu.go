package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/handler/dto"
)

type OfficeMenuHandle interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type officeMenuHandleImpl struct {
	userService usecase.UserService
	msgFormer   MessageFormer
	logger      *zap.Logger
}

func NewOfficeMenuHandle(userService usecase.UserService, logger *zap.Logger) OfficeMenuHandle {
	return &officeMenuHandleImpl{
		userService: userService,
		logger:      logger,
	}
}

func (o *officeMenuHandleImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.OfficeMenuFreeSeats:
		result, err := o.userService.CallDateMenu(ctx)
		if err != nil {
			return nil, err
		}
		return o.msgFormer.FormSeatListMsg(ctx, result)

	case dto.OfficeMenuSubscribe:
		result, err := o.userService.SubscribeWork(ctx)
		if err != nil {
			return nil, err
		}
		return o.msgFormer.FormTextMsg(ctx, result)

	case dto.OfficeMenuChooseAnotherOffice:
		result, err := o.userService.CallChooseOfficeMenu(ctx)
		if err != nil {
			return nil, err
		}
		return o.msgFormer.FormChooseOfficeMenuMsg(ctx, result)
	}

	return nil, nil
}
