package router

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/handler"
)

type Router struct {
	customMessageHandler handler.CustomMessageHandler
	commandHandler       handler.CommandHandler
	inlineHandler        handler.InlineMessageHandler
	logger               *zap.Logger
}

func NewRouter(customMessageHandler handler.CustomMessageHandler,
	commandHandler handler.CommandHandler,
	inlineHandler handler.InlineMessageHandler,
	logger *zap.Logger) Router {

	return Router{
		customMessageHandler: customMessageHandler,
		commandHandler:       commandHandler,
		inlineHandler:        inlineHandler,
		logger:               logger,
	}
}

func (r *Router) MainEntryPoint(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	if update.Message != nil {
		if update.Message.IsCommand() {
			return r.commandHandler.Handle(update)
		} else {
			return r.customMessageHandler.Handle(update)
		}
	} else if update.CallbackQuery != nil {
		return r.inlineHandler.Handle(update)
	}

	// TODO
	return nil, nil
}
