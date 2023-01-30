package router

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/handler"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type Router struct {
	userRepo             interfaces.UserRepository
	customMessageHandler handler.CustomMessageHandler
	commandHandler       handler.CommandHandler
	inlineHandler        handler.InlineMessageHandler
	logger               *zap.Logger
}

func NewRouter(userRepo interfaces.UserRepository, customMessageHandler handler.CustomMessageHandler,
	commandHandler handler.CommandHandler,
	inlineHandler handler.InlineMessageHandler,
	logger *zap.Logger) Router {

	return Router{
		userRepo:             userRepo,
		customMessageHandler: customMessageHandler,
		commandHandler:       commandHandler,
		inlineHandler:        inlineHandler,
		logger:               logger,
	}
}

func (r *Router) MainEntryPoint(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {

	if update.Message != nil {
		ctx := r.setUserContext(update.Message.From.ID,
			update.Message.Chat.ID, update.Message.MessageID)
		if update.Message.IsCommand() {
			return r.commandHandler.Handle(ctx, update)
		} else {
			return r.customMessageHandler.Handle(ctx, update)
		}
	} else if update.CallbackQuery != nil {
		ctx := r.setUserContext(update.CallbackQuery.From.ID,
			update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
		return r.inlineHandler.Handle(ctx, update)
	}

	// TODO
	return nil, nil
}

func (r *Router) setUserContext(telegramID, chatID int64, MessageID int) context.Context {
	ctx := context.Background()

	ctx = context.WithValue(ctx, model.ContextChatIDKey, chatID)
	ctx = context.WithValue(ctx, model.ContextMessageIDKey, MessageID)

	user, err := r.userRepo.GetByTelegramID(telegramID)
	if err != nil || user == nil {
		return ctx
	}
	ctx = context.WithValue(ctx, model.ContextUserKey, *user)
	return ctx
}
