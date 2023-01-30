package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	usecasedto "telegram-api/internal/app/usecase/dto"
)

type CommandHandler interface {
	Handle(ctx context.Context, update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type commandHandlerImpl struct {
	msgFormer   MessageFormer
	userService usecase.UserService
	logger      *zap.Logger
}

func NewCommandHandler(msgFormer MessageFormer, userService usecase.UserService, logger *zap.Logger) CommandHandler {
	return &commandHandlerImpl{
		msgFormer:   msgFormer,
		userService: userService,
		logger:      logger,
	}
}

func (s *commandHandlerImpl) Handle(ctx context.Context, update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {

	switch update.Message.Command() {
	case "start":
		return s.handleStartCommand(ctx, update)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Я незнаю этой команды"
	return &msg, nil
}

func (s *commandHandlerImpl) handleStartCommand(ctx context.Context, update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {

	result, err := s.userService.FirstCome(ctx)
	if err != nil || result == nil {
		return nil, err
	}

	switch result.Key {
	case usecase.OfficeMenu:
		return s.sendOfficeMenu(ctx, result)
	case usecase.ChooseOfficeMenu:
		return s.sendChooseOfficeMenu(ctx, result)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестный вызов")
	return &msg, nil
}

func (s *commandHandlerImpl) sendOfficeMenu(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	return s.msgFormer.FormOfficeMenuMsg(ctx, result)
}

func (s *commandHandlerImpl) sendChooseOfficeMenu(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	return s.msgFormer.FormChooseOfficeMenuMsg(ctx, result)
}
