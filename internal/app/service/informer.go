package service

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"log"
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

	if todayUTC == bookSeat.BookDate && currentTime.After(eveningTime) {
		return nil
	}
	return i.sendNotifies(ctx, bookSeat.Office.ID)

	//if todayUTC == bookSeat.BookDate {
	//	currentTime, err := service.CurrentTimeWithTimeZone(bookSeat.Office.TimeZone)
	//	eveningTime, err := service.EveningTimeWithTimeZone(bookSeat.Office.TimeZone)
	//	if err != nil {
	//		i.logger.Error("Error TodayWithTimeZone", zap.Error(err))
	//		// Ошибку не возвращаем, show must go on
	//		//return err
	//	}
	//
	//	if currentTime.Before(eveningTime) {
	//		i.sendNotifies(ctx)
	//	}
	//} else {
	//	i.sendNotifies(ctx)
	//}
	//return nil
}

func (i *informerServiceImpl) sendNotifies(ctx context.Context, officeID int64) error {
	users, err := i.userRepo.GetUsersToNotify(officeID)
	if err != nil {
		return err
	}

	for _, user := range users {
		msg := tgbotapi.NewMessage(user.ChatID, "Освободилось место")
		fmt.Println(msg)
		if _, err = i.botAPI.Send(msg); err != nil {
			log.Printf("Error when try to send message %d", err)
			continue
		}
	}

	return nil
}
