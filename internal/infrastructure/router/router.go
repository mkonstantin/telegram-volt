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
	Route(ctx context.Context, data Data) (*tgbotapi.MessageConfig, error)
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
	case "how_works":
		text := "Бот поможет забронировать рабочее место в одном из наших офисов." +
			"\nДля того чтобы забронировать рабочее место, выбери нужный офис, и номер рабочего стола" +
			"\nВозможность забронировать открывается каждый будний день в 14:00 автоматически." +
			"\nЧтобы не пропускать уведомление, подпишитесь на интересующий вас офис, и вы сможете получать напоминание об открытии бронирования." +
			"\nВ день бронирования, в 9:00 откроется возможность подтвердить бронь, если этого не сделать до 10:00, бронь будет аннулирована"
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
