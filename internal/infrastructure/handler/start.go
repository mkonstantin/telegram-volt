package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/domain/model"
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

	currentUser := model.GetCurrentUser(ctx)

	if currentUser.HaveChosenOffice() {
		result, err := s.userService.CallOfficeMenu(ctx)
		if err != nil {
			return nil, err
		}
		return s.msgFormer.FormOfficeMenuMsg(ctx, result)
	} else {
		result, err := s.userService.CallChooseOfficeMenu(ctx)
		if err != nil {
			return nil, err
		}
		return s.msgFormer.FormChooseOfficeMenuMsg(ctx, result)
	}
}
