package router

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/handler"
	"telegram-api/internal/infrastructure/handler/dto"
)

type Data struct {
	Command string
	Request dto.InlineRequest
}

type Router interface {
	Route(ctx context.Context, data Data) (*tgbotapi.MessageConfig, error)
}

type routerImpl struct {
	startHandle        handler.StartHandle
	officeListHandle   handler.OfficeListHandle
	officeMenuHandle   handler.OfficeMenuHandle
	seatListHandle     handler.SeatListHandle
	ownSeatMenuHandle  handler.OwnSeatMenuHandle
	freeSeatMenuHandle handler.FreeSeatMenuHandle
	logger             *zap.Logger
}

func NewRouter(
	startHandle handler.StartHandle,
	officeListHandle handler.OfficeListHandle,
	officeMenuHandle handler.OfficeMenuHandle,
	seatListHandle handler.SeatListHandle,
	ownSeatMenuHandle handler.OwnSeatMenuHandle,
	freeSeatMenuHandle handler.FreeSeatMenuHandle,
	logger *zap.Logger) Router {

	return &routerImpl{
		startHandle:        startHandle,
		officeListHandle:   officeListHandle,
		officeMenuHandle:   officeMenuHandle,
		seatListHandle:     seatListHandle,
		ownSeatMenuHandle:  ownSeatMenuHandle,
		freeSeatMenuHandle: freeSeatMenuHandle,
		logger:             logger,
	}
}

func (r *routerImpl) Route(ctx context.Context, data Data) (*tgbotapi.MessageConfig, error) {

	if data.Command != "" {
		return r.command(ctx, data.Command)
	}

	if data.Request.Type != "" {
		return r.inline(ctx, data.Request)
	}

	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "Ничего не могу с этим сделать)")
	return &msg, nil
}

func (r *routerImpl) command(ctx context.Context, command string) (*tgbotapi.MessageConfig, error) {
	switch command {
	case "start":
		return r.startHandle.Handle(ctx)
	}

	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "Неизвестная команда =(")
	return &msg, nil
}

func (r *routerImpl) inline(ctx context.Context, request dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch request.Type {
	case usecase.ChooseOfficeMenu:
		return r.officeListHandle.Handle(ctx, request)
	case usecase.OfficeMenu:
		return r.officeMenuHandle.Handle(ctx, request)
	case usecase.ChooseSeatsMenu:
		return r.seatListHandle.Handle(ctx, request)
	case usecase.SeatOwn:
		return r.ownSeatMenuHandle.Handle(ctx, request)
	case usecase.SeatFree:
		return r.freeSeatMenuHandle.Handle(ctx, request)
	}

	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "Неизвестная команда =(")
	return &msg, nil
}
