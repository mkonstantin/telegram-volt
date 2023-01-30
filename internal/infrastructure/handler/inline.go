package handler

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/handler/dto"
)

type InlineMessageHandler interface {
	Handle(ctx context.Context, update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
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

func (s *inlineMessageHandlerImpl) Handle(ctx context.Context, update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
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
		return s.chooseOfficeMenuTap(ctx, command)
	case usecase.OfficeMenu:
		return s.officeMenuTapScript(ctx, command)
	case usecase.ChooseSeatsMenu:
		return s.chooseSeatsMenuTap(ctx, command)
	case usecase.SeatOwn:
		return s.chooseSeatOwnTap(ctx, command)
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

func (s *inlineMessageHandlerImpl) chooseOfficeMenuTap(ctx context.Context, command *dto.CommandResponse) (*tgbotapi.MessageConfig, error) {
	result, err := s.userService.SetOfficeScript(ctx, command.OfficeID)
	if err != nil {
		return nil, err
	}

	return s.msgFormer.FormOfficeMenuMsg(ctx, result)
}

func (s *inlineMessageHandlerImpl) officeMenuTapScript(ctx context.Context, command *dto.CommandResponse) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.OfficeMenuFreeSeats:
		return s.callSeatsMenu(ctx)

	case dto.OfficeMenuSubscribe:

	case dto.OfficeMenuChooseAnotherOffice:
		result, err := s.userService.CallChooseOfficeMenu(ctx)
		if err != nil {
			return nil, err
		}
		return s.msgFormer.FormChooseOfficeMenuMsg(ctx, result)
	}

	return nil, nil
}

func (s *inlineMessageHandlerImpl) chooseSeatsMenuTap(ctx context.Context, command *dto.CommandResponse) (*tgbotapi.MessageConfig, error) {

	result, err := s.userService.BookSeatTap(ctx, command.BookSeatID)
	if err != nil {
		return nil, err
	}

	return s.msgFormer.FormBookSeatMsg(ctx, result)
}

func (s *inlineMessageHandlerImpl) chooseSeatOwnTap(ctx context.Context, command *dto.CommandResponse) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.ActionCancelBookYes:

	case dto.ActionCancelBookNo:
		return s.callSeatsMenu(ctx)
	}
	// TODO
	return nil, nil
}

func (s *inlineMessageHandlerImpl) callSeatsMenu(ctx context.Context) (*tgbotapi.MessageConfig, error) {
	result, err := s.userService.CallSeatsMenu(ctx)
	if err != nil {
		return nil, err
	}
	return s.msgFormer.FormSeatListMsg(ctx, result)
}
