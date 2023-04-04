package informer

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/informer/form"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type InformerService interface {
	SeatComeFree(ctx context.Context, id int64) error
	SendNotifiesWithMessage(office model.Office, message string) error
}

type informerServiceImpl struct {
	botAPI       *tgbotapi.BotAPI
	infoForm     form.InfoMenuForm
	userRepo     interfaces.UserRepository
	bookSeatRepo interfaces.BookSeatRepository
	logger       *zap.Logger
}

func NewInformer(botAPI *tgbotapi.BotAPI, infoForm form.InfoMenuForm, userRepo interfaces.UserRepository,
	bookSeatRepo interfaces.BookSeatRepository, logger *zap.Logger) InformerService {
	return &informerServiceImpl{
		botAPI:       botAPI,
		infoForm:     infoForm,
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

	todayUTC := helper.TodayZeroTimeUTC()
	currentTime, err := helper.CurrentTimeWithTimeZone(bookSeat.Office.TimeZone)
	eveningTime, err := helper.EveningTimeWithTimeZone(bookSeat.Office.TimeZone)

	if err != nil {
		i.logger.Error("Error TodayWithTimeZone", zap.Error(err))
		// Ошибку не возвращаем, show must go on
		//return err
	}

	if bookSeat.BookDate.Before(todayUTC) || (todayUTC == bookSeat.BookDate && currentTime.After(eveningTime)) {
		return nil
	}
	return i.chooseUsersAndSendNotifies(ctx, bookSeat)
}

func (i *informerServiceImpl) chooseUsersAndSendNotifies(ctx context.Context, bookSeat *model.BookSeat) error {
	formattedDate := bookSeat.BookDate.Format(helper.DateFormat)
	text := fmt.Sprintf("Освободилось место в офисе: %s на %s", bookSeat.Office.Name, formattedDate)

	currentUser := model.GetCurrentUser(ctx)

	users, err := i.userRepo.GetUsersToNotify(bookSeat.Office.ID)
	if err != nil {
		return err
	}

	seats, err := i.bookSeatRepo.GetUsersByOfficeIDAndDate(bookSeat.Office.ID, bookSeat.BookDate.String())
	var mapper = make(map[int64]int)
	mapper[currentUser.ID]++
	for _, seat := range seats {
		mapper[seat.User.ID]++
	}

	data := form.InfoFormData{
		Message: text,
		Office:  &bookSeat.Office,
		Date:    bookSeat.BookDate,
	}

	for _, user := range users {
		if mapper[user.ID] > 0 {
			continue
		}
		if currentUser.ID != user.ID {
			data.ChatID = user.ChatID
			i.sendInfoForm(ctx, data)
		}
	}

	return nil
}

func (i *informerServiceImpl) sendInfoForm(ctx context.Context, data form.InfoFormData) {
	build, err := i.infoForm.Build(ctx, data)
	if err != nil {
		i.logger.Error("Error Informer sendInfoForm", zap.Error(err))
		// Ошибку не возвращаем, show must go on
	}
	if _, err = i.botAPI.Send(build); err != nil {
		i.logger.Error("Error when try to send NOTIFY", zap.Error(err))
	}
}

func (i *informerServiceImpl) SendNotifiesWithMessage(office model.Office, message string) error {
	users, err := i.userRepo.GetUsersToNotify(office.ID)
	if err != nil {
		return err
	}

	for _, user := range users {
		i.sendMessage(user.ChatID, message)
	}

	return nil
}

func (i *informerServiceImpl) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := i.botAPI.Send(msg); err != nil {
		i.logger.Error("Error when try to send NOTIFY", zap.Error(err))
	}
}
