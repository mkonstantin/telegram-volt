package router

import (
	"context"
	"fmt"
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
		fullName := fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)
		ctx := r.setUserContext(update.Message.From.ID,
			update.Message.Chat.ID, update.Message.MessageID,
			update.Message.From.UserName, fullName)

		if update.Message.IsCommand() {
			return r.commandHandler.Handle(ctx, update)
		} else {
			return r.customMessageHandler.Handle(ctx, update)
		}
	} else if update.CallbackQuery != nil {
		fullName := fmt.Sprintf("%s %s", update.CallbackQuery.Message.From.FirstName, update.CallbackQuery.Message.From.LastName)
		ctx := r.setUserContext(update.CallbackQuery.From.ID,
			update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID,
			update.CallbackQuery.From.UserName, fullName)

		return r.inlineHandler.Handle(ctx, update)
	}

	// TODO
	return nil, nil
}

func (r *Router) setUserContext(tgID, chatID int64, MessageID int, tgUserName, fullName string) context.Context {
	ctx := context.Background()

	ctx = context.WithValue(ctx, model.ContextChatIDKey, chatID)
	ctx = context.WithValue(ctx, model.ContextMessageIDKey, MessageID)

	user, err := r.userRepo.GetByTelegramID(tgID)
	if err != nil {
		return ctx
	}
	if user == nil {
		user, err = r.createUser(tgID, tgUserName, fullName)
	}
	if err != nil {
		return ctx
	}
	ctx = context.WithValue(ctx, model.ContextUserKey, *user)
	return ctx
}

func (r *Router) createUser(tgID int64, tgUserName, fullName string) (*model.User, error) {
	userModel := model.User{
		Name:         fullName,
		TelegramID:   tgID,
		TelegramName: tgUserName,
	}
	err := r.userRepo.Create(userModel)
	if err != nil {
		return nil, err
	}

	user, err := r.userRepo.GetByTelegramID(tgID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
