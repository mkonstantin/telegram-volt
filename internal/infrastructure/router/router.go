package router

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"telegram-api/config"
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
	Route(ctx context.Context, data Data) (tgbotapi.Chattable, error)
}

type routerImpl struct {
	cfg          config.AppConfig
	start        handler2.Start
	officeList   handler2.OfficeList
	officeMenu   handler2.OfficeMenu
	seatList     handler2.SeatList
	ownSeatMenu  handler2.OwnSeatMenu
	freeSeatMenu handler2.FreeSeatMenu
	dateMenu     handler2.DateMenu
	infoMenu     handler2.InfoMenu
	holdSeatMenu handler2.HoldSeatMenu
	logger       *zap.Logger
}

func NewRouter(
	cfg config.AppConfig,
	startHandle handler2.Start,
	officeListHandle handler2.OfficeList,
	officeMenuHandle handler2.OfficeMenu,
	seatListHandle handler2.SeatList,
	ownSeatMenuHandle handler2.OwnSeatMenu,
	freeSeatMenuHandle handler2.FreeSeatMenu,
	dateMenu handler2.DateMenu,
	infoMenu handler2.InfoMenu,
	holdSeatMenu handler2.HoldSeatMenu,
	logger *zap.Logger) Router {

	return &routerImpl{
		cfg:          cfg,
		start:        startHandle,
		officeList:   officeListHandle,
		officeMenu:   officeMenuHandle,
		seatList:     seatListHandle,
		ownSeatMenu:  ownSeatMenuHandle,
		freeSeatMenu: freeSeatMenuHandle,
		dateMenu:     dateMenu,
		infoMenu:     infoMenu,
		holdSeatMenu: holdSeatMenu,
		logger:       logger,
	}
}

func (r *routerImpl) Route(ctx context.Context, data Data) (tgbotapi.Chattable, error) {

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

func (r *routerImpl) command(ctx context.Context, command string) (tgbotapi.Chattable, error) {
	switch command {
	case "start":
		return r.start.Handle(ctx)
	case "how_works":
		text := "Бот умеет отображать свободные места в офисе, а также напоминать о том, что бронирование открылось. " +
			"\nДля того чтобы не пропускать уведомления от бота об открытии записи необходимо оформить подписку на интересующий офис (без регистрации и смс, шутка)) " +
			"\nВ этом случае ежедневно в 14:00 вам будет направлено сообщение о доступности бронирования. " +
			"\nБронирование будет автоматически открываться в 14:00, к сожалению, бот пока не умеет определять выходные и праздничные дни. " +
			"И бронирование на понедельник будет открываться в ВОСКРЕСЕНЬЕ!" +
			"\nВ день бронирования с 09:00 - 10:00 утра вам будет необходимо подтвердить свое намерение прийти в офис, в противном случае бронь отменится. " +
			"\nЕсли вы не успели записаться и все места разобрали, подпишитесь на нужный офис и вы получите уведомление если кто-то решит удалить свою запись."
		msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), text)
		return &msg, nil
	case "about":
		text := fmt.Sprintf("developed by @km505603"+
			"\nversion: %s", r.cfg.Version)
		msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), text)
		return &msg, nil
	}

	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "Неизвестная команда =(")
	return &msg, nil
}

func (r *routerImpl) inline(ctx context.Context, request dto.InlineRequest) (tgbotapi.Chattable, error) {

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
	case constants.HoldSeatMenuTap:
		return r.holdSeatMenu.Handle(ctx, request)
	}

	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "Неизвестная команда =(")
	return &msg, nil
}
