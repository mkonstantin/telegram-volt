package service

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/internal/infrastructure/service"
)

type InformerService interface {
	SeatComeFree(ctx context.Context, id int64) error
}

type informerServiceImpl struct {
	botAPI       *tgbotapi.BotAPI
	officeRepo   int
	userRepo     interfaces.UserRepository
	bookSeatRepo interfaces.BookSeatRepository
	logger       *zap.Logger
}

func NewInformer(botAPI *tgbotapi.BotAPI, userRepo interfaces.UserRepository,
	bookSeatRepo interfaces.BookSeatRepository, logger *zap.Logger) InformerService {
	return &informerServiceImpl{
		botAPI:       botAPI,
		userRepo:     userRepo,
		bookSeatRepo: bookSeatRepo,
		logger:       logger,
	}
}

func (i *informerServiceImpl) SeatComeFree(ctx context.Context, bookSeatID int64) error {

	bookSeat, err := i.bookSeatRepo.FindByID(bookSeatID)
	if err != nil {
		return err
	}

	todayUTC := service.TodayZeroTimeUTC()
	currentTime, err := service.CurrentTimeWithTimeZone(bookSeat.Office.TimeZone)
	eveningTime, err := service.EveningTimeWithTimeZone(bookSeat.Office.TimeZone)

	if err != nil {
		i.logger.Error("Error TodayWithTimeZone", zap.Error(err))
		// Ошибку не возвращаем, show must go on
		//return err
	}

	if bookSeat.BookDate.Before(todayUTC) || (todayUTC == bookSeat.BookDate && currentTime.After(eveningTime)) {
		return nil
	}
	return i.sendNotifies(ctx, bookSeat.Office)
}

func (i *informerServiceImpl) sendNotifies(ctx context.Context, office model.Office) error {
	text := fmt.Sprintf("Освободилось место в офисе: %s", office.Name)

	currentUserChatID := model.GetCurrentChatID(ctx)

	users, err := i.userRepo.GetUsersToNotify(office.ID)
	if err != nil {
		return err
	}

	for _, user := range users {
		if currentUserChatID != user.ChatID {
			i.sendMessage(user.ChatID, text)
		}
	}

	return nil
}

func (i *informerServiceImpl) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := i.botAPI.Send(msg); err != nil {
		i.logger.Error("Error when try to send NOTIFY", zap.Error(err))
	}
}
