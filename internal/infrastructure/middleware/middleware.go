package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/handler/dto"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/internal/infrastructure/router"
)

type UserMW struct {
	userRepo interfaces.UserRepository
	router   router.Router
	logger   *zap.Logger
}

func NewUserMW(userRepo interfaces.UserRepository, router router.Router,
	logger *zap.Logger) UserMW {

	return UserMW{
		userRepo: userRepo,
		router:   router,
		logger:   logger,
	}
}

func (r *UserMW) EntryPoint(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {

	var ctx context.Context
	var data router.Data

	if update.Message != nil {
		fullName := fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)
		ctx = r.setUserContext(
			update.Message.From.ID,
			update.Message.Chat.ID,
			update.Message.MessageID,
			update.Message.From.UserName,
			fullName)

		if update.Message.IsCommand() {
			data.Command = update.Message.Command()
		}
	} else if update.CallbackQuery != nil {
		fullName := fmt.Sprintf("%s %s", update.CallbackQuery.Message.From.FirstName, update.CallbackQuery.Message.From.LastName)
		ctx = r.setUserContext(
			update.CallbackQuery.From.ID,
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			update.CallbackQuery.From.UserName,
			fullName)

		cData, err := extractData(update.CallbackQuery.Data)
		if err != nil {
			return nil, err
		}
		data.Data = cData
	}

	return r.router.Route(ctx, data)
}

func extractData(callbackData string) (dto.InlineRequest, error) {
	command := dto.InlineRequest{}

	err := json.Unmarshal([]byte(callbackData), &command)
	if err != nil {
		return command, err
	}
	return command, nil
}

func (r *UserMW) setUserContext(tgID, chatID int64, MessageID int, tgUserName, fullName string) context.Context {
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

func (r *UserMW) createUser(tgID int64, tgUserName, fullName string) (*model.User, error) {
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
