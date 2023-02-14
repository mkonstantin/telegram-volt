package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/handler/dto"
)

type OwnSeatMenuHandle interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type ownSeatMenuHandleImpl struct {
	userService usecase.UserService
	msgFormer   MessageFormer
	logger      *zap.Logger
}

func NewOwnSeatMenuHandle(userService usecase.UserService, logger *zap.Logger) OwnSeatMenuHandle {
	return &ownSeatMenuHandleImpl{
		userService: userService,
		logger:      logger,
	}
}

func (o *ownSeatMenuHandleImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.ActionCancelBookYes:
		result, err := o.userService.CancelBookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		return o.msgFormer.FormCancelBookResult(ctx, result)
	case dto.ActionCancelBookNo:
		fallthrough
	default:
		return o.callSeatsMenu(ctx)
	}
}

func (o *ownSeatMenuHandleImpl) callSeatsMenu(ctx context.Context) (*tgbotapi.MessageConfig, error) {
	result, err := o.userService.CallSeatsMenu(ctx)
	if err != nil {
		return nil, err
	}
	return o.msgFormer.FormSeatListMsg(ctx, result)
}
