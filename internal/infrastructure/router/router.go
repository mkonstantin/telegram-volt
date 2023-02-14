package router

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
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
	start        handler.Start
	officeList   handler.OfficeList
	officeMenu   handler.OfficeMenu
	seatList     handler.SeatList
	ownSeatMenu  handler.OwnSeatMenu
	freeSeatMenu handler.FreeSeatMenu
	logger       *zap.Logger
}

func NewRouter(
	startHandle handler.Start,
	officeListHandle handler.OfficeList,
	officeMenuHandle handler.OfficeMenu,
	seatListHandle handler.SeatList,
	ownSeatMenuHandle handler.OwnSeatMenu,
	freeSeatMenuHandle handler.FreeSeatMenu,
	logger *zap.Logger) Router {

	return &routerImpl{
		start:        startHandle,
		officeList:   officeListHandle,
		officeMenu:   officeMenuHandle,
		seatList:     seatListHandle,
		ownSeatMenu:  ownSeatMenuHandle,
		freeSeatMenu: freeSeatMenuHandle,
		logger:       logger,
	}
}

func (r *routerImpl) Route(ctx context.Context, data Data) (*tgbotapi.MessageConfig, error) {

	if data.Command != "" {
		return r.command(ctx, data.Command)
	}

	if data.Request.Type != "" {
		return r.inline(ctx, data.Request)
	}

	r.logger.Error("Route. Unknown data", zap.Error(errors.New("Unknown data")))
	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "Ничего не могу с этим сделать)")
	return &msg, nil
}

func (r *routerImpl) command(ctx context.Context, command string) (*tgbotapi.MessageConfig, error) {
	switch command {
	case "start":
		return r.start.Handle(ctx)
	}

	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "Неизвестная команда =(")
	return &msg, nil
}

func (r *routerImpl) inline(ctx context.Context, request dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch request.Type {
	case usecase.ChooseOfficeMenu:
		return r.officeList.Handle(ctx, request)
	case usecase.CallOfficeMenu:
		return r.officeMenu.Handle(ctx, request)
	case usecase.ChooseSeatsMenu:
		return r.seatList.Handle(ctx, request)
	case usecase.SeatOwn:
		return r.ownSeatMenu.Handle(ctx, request)
	case usecase.SeatFree:
		return r.freeSeatMenu.Handle(ctx, request)
	}

	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "Неизвестная команда =(")
	return &msg, nil
}
