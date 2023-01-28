package handler

import (
	"encoding/json"
	"fmt"
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
		return s.officeChosenScenery(update.CallbackQuery.From.ID, command.ChooseOffice.OfficeID)
	case usecase.ConfirmOffice:
		return s.officeConfirmScenery(command, update)
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

func (s *inlineMessageHandlerImpl) officeConfirmScenery(command *dto.CommandResponse,
	update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	if command.ConfirmOffice.IsConfirm {
		return s.officeChosenScenery(update.CallbackQuery.From.ID, command.ChooseOffice.OfficeID)
	} else {

	}
	return nil, nil
}

func (s *inlineMessageHandlerImpl) officeChosenScenery(TelegramID, OfficeID int64) (*tgbotapi.MessageConfig, error) {

	data := dto2.OfficeChosenDTO{
		TelegramID: TelegramID,
		OfficeID:   OfficeID,
	}

	result, err := s.userService.OfficeChosenScenery(data)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)
	return nil, nil
}
