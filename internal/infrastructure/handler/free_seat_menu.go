package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/handler/dto"
)

type FreeSeatMenuHandle interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type freeSeatMenuHandleImpl struct {
	userService usecase.UserService
	msgFormer   MessageFormer
	logger      *zap.Logger
}

func NewFreeSeatMenuHandle(userService usecase.UserService, logger *zap.Logger) FreeSeatMenuHandle {
	return &freeSeatMenuHandleImpl{
		userService: userService,
		logger:      logger,
	}
}

func (f *freeSeatMenuHandleImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.ActionBookYes:
		result, err := f.userService.BookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		return f.msgFormer.FormBookSeatResult(ctx, result)

	case dto.ActionBookNo:
		fallthrough
	default:
		return f.callSeatsMenu(ctx)
	}
}

func (f *freeSeatMenuHandleImpl) callSeatsMenu(ctx context.Context) (*tgbotapi.MessageConfig, error) {
	result, err := f.userService.CallSeatsMenu(ctx)
	if err != nil {
		return nil, err
	}
	return f.msgFormer.FormSeatListMsg(ctx, result)
}
