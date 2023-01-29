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
	case usecase.ChooseOfficeMenu:
		return s.chooseOfficeMenuTap(update.CallbackQuery.From.ID, command.OfficeID, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	case usecase.OfficeMenu:
		return s.officeMenuTapScript(command, update)
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

func (s *inlineMessageHandlerImpl) chooseOfficeMenuTap(telegramID, officeID, chatID int64, messageID int) (*tgbotapi.MessageConfig, error) {

	data := dto2.OfficeChosenDTO{
		TelegramID: telegramID,
		OfficeID:   officeID,
		ChatID:     chatID,
		MessageID:  messageID,
	}

	result, err := s.userService.SetOfficeScript(data)
	if err != nil {
		return nil, err
	}

	return s.msgFormer.FormOfficeMenuMsg(result)
}

func (s *inlineMessageHandlerImpl) officeMenuTapScript(command *dto.CommandResponse, update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.OfficeMenuFreeSeats:
		dtoO := dto2.BookSeatDTO{
			TelegramID: update.CallbackQuery.From.ID,
			OfficeID:   command.OfficeID,
			MessageID:  update.CallbackQuery.Message.MessageID,
			ChatID:     update.CallbackQuery.Message.Chat.ID,
		}

		result, err := s.userService.CallSeatsMenu(dtoO)
		if err != nil {
			return nil, err
		}
		return s.msgFormer.FormSeatListMsg(result)

	case dto.OfficeMenuSubscribe:

	case dto.OfficeMenuChooseAnotherOffice:
		startDTO := dto2.FirstStartDTO{
			User: model.User{
				Name:         update.CallbackQuery.From.FirstName,
				TelegramID:   update.CallbackQuery.From.ID,
				TelegramName: update.CallbackQuery.From.UserName,
			},
			MessageID: update.CallbackQuery.Message.MessageID,
			ChatID:    update.CallbackQuery.Message.Chat.ID,
		}

		result, err := s.userService.CallChooseOfficeMenu(startDTO)
		if err != nil {
			return nil, err
		}
		return s.msgFormer.FormChooseOfficeMenuMsg(result)
	}

	return nil, nil
}
