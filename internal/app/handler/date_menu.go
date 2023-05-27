package handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"io/ioutil"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
	"telegram-api/pkg/tracing"
)

type DateMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (tgbotapi.Chattable, error)
}

type dateMenuImpl struct {
	seatList   interfaces.SeatListMenu
	officeMenu interfaces.OfficeMenu
	logger     *zap.Logger
}

func NewDateMenuHandle(
	seatList interfaces.SeatListMenu,
	officeMenu interfaces.OfficeMenu,
	logger *zap.Logger) DateMenu {

	return &dateMenuImpl{
		seatList:   seatList,
		officeMenu: officeMenu,
		logger:     logger,
	}
}

func (o *dateMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (tgbotapi.Chattable, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	switch command.Action {
	case dto.DateListShowMap:
		currentUser := model.GetCurrentUser(ctx)
		path := fmt.Sprintf("./picture/%d.jpg", currentUser.OfficeID)

		photoBytes, err := ioutil.ReadFile(path)
		if err != nil {
			o.logger.Warn("error while trying send photo", zap.Error(err))

			msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "")
			msg.Text = "К сожалению, пока нет карты расположения мест для этого офиса"
			return msg, nil
		}

		photoFileBytes := tgbotapi.FileBytes{
			Name:  fmt.Sprintf("Map of office: %d", currentUser.OfficeID),
			Bytes: photoBytes,
		}
		msgPhoto := tgbotapi.NewPhoto(model.GetCurrentChatID(ctx), photoFileBytes)
		return msgPhoto, err

	case dto.Back:
		return o.officeMenu.Call(ctx, "", 0)
	}

	if command.BookDate == nil {
		return o.officeMenu.Call(ctx, "", 0)
	}

	return o.seatList.Call(ctx, *command.BookDate, 0)
}
