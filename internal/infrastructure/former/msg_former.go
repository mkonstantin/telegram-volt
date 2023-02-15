package former

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	usecasedto "telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
)

type MessageFormer interface {
	FormBookSeatResult(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormCancelBookResult(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormTextMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
}

type messageFormerImpl struct {
	logger *zap.Logger
}

func NewMessageFormer(logger *zap.Logger) MessageFormer {
	return &messageFormerImpl{
		logger: logger,
	}
}

func (s *messageFormerImpl) FormBookSeatResult(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {

	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = result.Message

	return &msg, nil
}

func (s *messageFormerImpl) FormCancelBookResult(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = result.Message

	return &msg, nil
}

func (s *messageFormerImpl) FormTextMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = result.Message

	return &msg, nil
}
