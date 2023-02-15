package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/former"
	"telegram-api/internal/infrastructure/handler/dto"
)

type SeatList interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type seatListImpl struct {
	userService usecase.UserService
	msgFormer   former.MessageFormer
	logger      *zap.Logger
}

func NewSeatListHandle(
	userService usecase.UserService,
	msgFormer former.MessageFormer,
	logger *zap.Logger) SeatList {

	return &seatListImpl{
		userService: userService,
		msgFormer:   msgFormer,
		logger:      logger,
	}
}

func (s *seatListImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	result, err := s.userService.SeatListTap(ctx, command.BookSeatID)
	if err != nil {
		return nil, err
	}

	return s.msgFormer.FormBookSeatMsg(ctx, result)
}
