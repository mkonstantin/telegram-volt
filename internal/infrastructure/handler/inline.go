package handler

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	dto2 "telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/handler/dto"
)

type InlineMessageHandler interface {
	Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type inlineMessageHandlerImpl struct {
	msgFormer   MessageFormer
	userService usecase.UserService
	logger      *zap.Logger
}

func NewInlineMessageHandler(msgFormer MessageFormer, userService usecase.UserService, logger *zap.Logger) InlineMessageHandler {
	return &inlineMessageHandlerImpl{
		msgFormer:   msgFormer,
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
	case usecase.OfficeMenu:
		return s.officeMenuTapScript(command, update)
	case usecase.ChooseOfficeMenu:
		return s.chooseOfficeMenuTap(update.CallbackQuery.From.ID, command.OfficeID, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
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

func (s *inlineMessageHandlerImpl) officeMenuTapScript(command *dto.CommandResponse, update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.OfficeMenuFreeSeats:

	case dto.OfficeMenuSubscribe:

	case dto.OfficeMenuChooseAnotherOffice:
		startDTO := dto2.FirstStartDTO{
			User:      model.User{},
			MessageID: 0,
			ChatID:    0,
		}
		result, err := s.userService.CallChooseOfficeMenu(startDTO)
		if err != nil {
			return nil, err
		}
		return s.msgFormer.FormChooseOfficeMenuMsg(result)
	}

	return nil, nil
}

func (s *inlineMessageHandlerImpl) chooseOfficeMenuTap(telegramID, officeID, chatID int64, messageID int) (*tgbotapi.MessageConfig, error) {

	data := dto2.OfficeChosenDTO{
		TelegramID: telegramID,
		OfficeID:   officeID,
		ChatID:     chatID,
		MessageID:  messageID,
	}

	result, err := s.userService.OfficeChosenScenery(data)
	if err != nil {
		return nil, err
	}

	return s.msgFormer.FormOfficeMenuMsg(result)
}
