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
		//ctx := r.setUserContext(update.Message.From.ID)
		if update.Message.IsCommand() {
			return r.commandHandler.Handle(update)
		} else {
			return r.customMessageHandler.Handle(update)
		}
	} else if update.CallbackQuery != nil {
		ctx := r.setUserContext(update.CallbackQuery.From.ID)
		return r.inlineHandler.Handle(ctx, update)
	}

	// TODO
	return nil, nil
}

func (r *Router) setUserContext(id int64) context.Context {
	ctx := context.Background()

	user, err := r.userRepo.GetByTelegramID(id)
	if err != nil || user == nil {
		return ctx
	}

	ctx = context.WithValue(ctx, model.ContextUserKey, *user)
	return ctx
}
