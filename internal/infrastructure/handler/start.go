package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/domain/model"
)

type StartHandle interface {
	Handle(ctx context.Context) (*tgbotapi.MessageConfig, error)
}

type startHandleImpl struct {
	userService usecase.UserService
	msgFormer   MessageFormer
	logger      *zap.Logger
}

func NewStartHandle(userService usecase.UserService, logger *zap.Logger) StartHandle {
	return &startHandleImpl{
		userService: userService,
		logger:      logger,
	}
}

func (s *startHandleImpl) Handle(ctx context.Context) (*tgbotapi.MessageConfig, error) {

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

	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "Неизвестный вызов")
	return &msg, nil
}
