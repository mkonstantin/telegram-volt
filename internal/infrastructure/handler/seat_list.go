package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/handler/dto"
)

type SeatListHandle interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type seatListHandleImpl struct {
	userService usecase.UserService
	msgFormer   MessageFormer
	logger      *zap.Logger
}

func NewSeatListHandle(userService usecase.UserService, logger *zap.Logger) SeatListHandle {
	return &seatListHandleImpl{
		userService: userService,
		logger:      logger,
	}
}

func (s *seatListHandleImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	result, err := s.userService.SeatListTap(ctx, command.BookSeatID)
	if err != nil {
		return nil, err
	}

	return s.msgFormer.FormBookSeatMsg(ctx, result)
}
