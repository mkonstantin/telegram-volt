package hundlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type InlineMessageHandler interface {
	Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type inlineMessageHandlerImpl struct {
	logger *zap.Logger
}

func NewInlineMessageHandler(logger *zap.Logger) InlineMessageHandler {
	return &inlineMessageHandlerImpl{
		logger: logger,
	}
}

func (s *inlineMessageHandlerImpl) Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	// Respond to the callback query, telling Telegram to show the user
	// a message with the data received.
	//callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	//if _, err := s.botAPI.Request(callback); err != nil {
	//	return nil, err
	//}

	// And finally, send a message containing the data received.
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	return &msg, nil
}
