package handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"os"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
)

type DateMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (tgbotapi.Chattable, error)
}

type dateMenuImpl struct {
	seatList   interfaces.SeatListMenu
	officeMenu interfaces.OfficeMenu
	officeRepo repo.OfficeRepository
	logger     *zap.Logger
}

func NewDateMenuHandle(
	seatList interfaces.SeatListMenu,
	officeMenu interfaces.OfficeMenu,
	officeRepo repo.OfficeRepository,
	logger *zap.Logger) DateMenu {

	return &dateMenuImpl{
		seatList:   seatList,
		officeMenu: officeMenu,
		officeRepo: officeRepo,
		logger:     logger,
	}
}

func (o *dateMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (tgbotapi.Chattable, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	switch command.Action {
	case dto.DateListShowMap:
		currentUser := model.GetCurrentUser(ctx)

		office, err := o.officeRepo.FindByID(ctx, currentUser.OfficeID)
		if err != nil {
			o.logger.Warn("error while get office images", zap.Error(err))

			return sendNoImageMessage(ctx)
		}

		image := office.Image
		if image == "" {
			return sendNoImageMessage(ctx)
		}

		path := fmt.Sprintf("./picture/%s.jpg", image)

		photoBytes, err := os.ReadFile(path)
		if err != nil {
			o.logger.Warn("error while trying send photo", zap.Error(err))

			return sendNoImageMessage(ctx)
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

func sendNoImageMessage(ctx context.Context) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(model.GetCurrentChatID(ctx), "")
	msg.Text = "К сожалению, пока нет карты расположения мест для этого хотдеска"
	return msg, nil
}
