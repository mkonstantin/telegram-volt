package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type CustomMessageHandler interface {
	Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type customMessageHandlerImpl struct {
	logger *zap.Logger
}

func NewCustomMessageHandler(logger *zap.Logger) CustomMessageHandler {
	return &customMessageHandlerImpl{
		logger: logger,
	}
}

func (s *customMessageHandlerImpl) Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ничего не могу с этим сделать)")
	return &msg, nil
}
