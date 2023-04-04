package router

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	handler2 "telegram-api/internal/app/handler"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/router/constants"
)

type Data struct {
	Command string
	Request dto.InlineRequest
}

type Router interface {
	Route(ctx context.Context, data Data) (*tgbotapi.MessageConfig, error)
}

type routerImpl struct {
	start        handler2.Start
	officeList   handler2.OfficeList
	officeMenu   handler2.OfficeMenu
	seatList     handler2.SeatList
	ownSeatMenu  handler2.OwnSeatMenu
	freeSeatMenu handler2.FreeSeatMenu
	dateMenu     handler2.DateMenu
	infoMenu     handler2.InfoMenu
	logger       *zap.Logger
}

func NewRouter(
	startHandle handler2.Start,
	officeListHandle handler2.OfficeList,
	officeMenuHandle handler2.OfficeMenu,
	seatListHandle handler2.SeatList,
	ownSeatMenuHandle handler2.OwnSeatMenu,
	freeSeatMenuHandle handler2.FreeSeatMenu,
	dateMenu handler2.DateMenu,
	infoMenu handler2.InfoMenu,
	logger *zap.Logger) Router {

	return &routerImpl{
		start:        startHandle,
		officeList:   officeListHandle,
		officeMenu:   officeMenuHandle,
		seatList:     seatListHandle,
		ownSeatMenu:  ownSeatMenuHandle,
		freeSeatMenu: freeSeatMenuHandle,
		dateMenu:     dateMenu,
		infoMenu:     infoMenu,
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
	case constants.OfficeListTap:
		return r.officeList.Handle(ctx, request)
	case constants.OfficeMenuTap:
		return r.officeMenu.Handle(ctx, request)
	case constants.SeatListTap:
		return r.seatList.Handle(ctx, request)
	case constants.OwnSeatMenuTap:
		return r.ownSeatMenu.Handle(ctx, request)
	case constants.FreeSeatMenuTap:
		return r.freeSeatMenu.Handle(ctx, request)
	case constants.DateMenuTap:
		return r.dateMenu.Handle(ctx, request)
	case constants.InformerTap:
		return r.infoMenu.Handle(ctx, request)
	}

	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "Неизвестная команда =(")
	return &msg, nil
}
