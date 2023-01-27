package handler

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	dto2 "telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/infrastructure/handler/dto"
)

type InlineMessageHandler interface {
	Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type inlineMessageHandlerImpl struct {
	userService usecase.UserService
	logger      *zap.Logger
}

func NewInlineMessageHandler(userService usecase.UserService, logger *zap.Logger) InlineMessageHandler {
	return &inlineMessageHandlerImpl{
		userService: userService,
		logger:      logger,
	}
}

func (s *inlineMessageHandlerImpl) Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	if update.CallbackQuery.Data == "" {
		// TODO
		return nil, nil
	}

	command, err := getCommand(update)
	if err != nil {
		return nil, err
	}

	switch command.Type {
	case usecase.ChooseOffice:
		return s.setOffice(update.CallbackQuery.From.ID, command.ChooseOffice.OfficeID)
	case usecase.ConfirmOffice:

	}

	// TODO
	return nil, nil
}

func getCommand(update tgbotapi.Update) (*dto.CommandResponse, error) {
	callbackData := update.CallbackQuery.Data
	command := dto.CommandResponse{}

	err := json.Unmarshal([]byte(callbackData), &command)
	if err != nil {
		return nil, err
	}
	return &command, nil
}

func (s *inlineMessageHandlerImpl) setOffice(telegramID, officeID int64) (*tgbotapi.MessageConfig, error) {
	data := dto2.SetOfficeDTO{
		TelegramID: telegramID,
		OfficeID:   officeID,
	}

	err := s.userService.SetOffice(data)
	if err != nil {
		return nil, err
	}

	// TODO
	return nil, nil
	//msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
}
