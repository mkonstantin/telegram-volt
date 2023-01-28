package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type MessageFormer interface {
}

type messageFormerImpl struct {
	logger *zap.Logger
}

func NewMessageFormer(logger *zap.Logger) MessageFormer {
	return &messageFormerImpl{
		logger: logger,
	}
}

func (s *messageFormerImpl) Start(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {

	return tgbotapi.MessageConfig{}, nil
}
