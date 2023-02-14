package handler

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
)

type Start interface {
	Handle(ctx context.Context) (*tgbotapi.MessageConfig, error)
}

type startImpl struct {
	userService usecase.UserService
	msgFormer   MessageFormer
	logger      *zap.Logger
}

func NewStartHandle(
	userService usecase.UserService,
	msgFormer MessageFormer,
	logger *zap.Logger) Start {

	return &startImpl{
		userService: userService,
		msgFormer:   msgFormer,
		logger:      logger,
	}
}

func (s *startImpl) Handle(ctx context.Context) (*tgbotapi.MessageConfig, error) {

	result, err := s.userService.FirstCome(ctx)
	if err != nil || result == nil {
		return nil, err
	}

	switch result.Key {
	case usecase.OfficeMenu:
		return s.msgFormer.FormOfficeMenuMsg(ctx, result)
	case usecase.ChooseOfficeMenu:
		return s.msgFormer.FormChooseOfficeMenuMsg(ctx, result)
	}

	s.logger.Error("Route. Unknown data", zap.Error(errors.New("unknown data")))
	return nil, nil
}
