package informer

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
)

type InformerService interface {
	SendNotifySeatBecomeFree(ctx context.Context, id int64) error
	SendNotifyTomorrowBookingOpen(ctx context.Context, office model.Office, message string) error
	SendNotifiesToConfirm(ctx context.Context, office *model.Office) error
	SendNotifyToBookDeletedBySystem(ctx context.Context, bookSeats []*model.BookSeat, officeName string) error
}

type informerServiceImpl struct {
	botAPI       *tgbotapi.BotAPI
	infoForm     form.InfoMenuForm
	userRepo     interfaces.UserRepository
	bookSeatRepo interfaces.BookSeatRepository
	sender       Sender
	logger       *zap.Logger
}

func NewInformer(botAPI *tgbotapi.BotAPI, infoForm form.InfoMenuForm, userRepo interfaces.UserRepository,
	bookSeatRepo interfaces.BookSeatRepository, sender Sender, logger *zap.Logger) InformerService {
	return &informerServiceImpl{
		botAPI:       botAPI,
		infoForm:     infoForm,
		userRepo:     userRepo,
		bookSeatRepo: bookSeatRepo,
		sender:       sender,
		logger:       logger,
	}
}

// SendNotifySeatBecomeFree Сообщение подписавшимся кроме тех кто уже занимает место, что место освободилось

func (i *informerServiceImpl) SendNotifySeatBecomeFree(ctx context.Context, bookSeatID int64) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	bookSeat, err := i.bookSeatRepo.FindByID(ctx, bookSeatID)
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
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	formattedDate := bookSeat.BookDate.Format(helper.DateFormat)
	text := fmt.Sprintf("Освободилось место в хотдеске: %s на %s", bookSeat.Office.Name, formattedDate)

	currentUser := model.GetCurrentUser(ctx)

	users, err := i.userRepo.GetUsersToNotify(ctx, bookSeat.Office.ID)
	if err != nil {
		return err
	}

	seats, err := i.bookSeatRepo.FindByOfficeIDAndDate(ctx, bookSeat.Office.ID, bookSeat.BookDate.String())
	var mapper = make(map[int64]int)
	mapper[currentUser.ID]++
	for _, seat := range seats {
		mapper[seat.User.ID]++
	}

	var arrayToSend []form.InfoFormData

	for _, user := range users {
		if mapper[user.ID] > 0 {
			continue
		}
		if currentUser.ID != user.ID {
			data := form.InfoFormData{
				Action:     dto.ActionShowSeatList,
				Message:    text,
				BookSeatID: bookSeat.ID,
				ChatID:     user.ChatID,
			}
			arrayToSend = append(arrayToSend, data)
		}
	}

	if len(arrayToSend) > 0 {
		i.sender.AddToQueue(arrayToSend)
	}
	return nil
}

// SendNotifiesToConfirm рассылка уведомлений на подтверждение брони

func (i *informerServiceImpl) SendNotifiesToConfirm(ctx context.Context, office *model.Office) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	today := helper.TodayZeroTimeUTC()
	bookSeats, err := i.bookSeatRepo.FindNotConfirmedByOfficeIDAndDate(ctx, office.ID, today.String())
	if err != nil {
		return err
	}

	var arrayToSend []form.InfoFormData

	for _, bookSeat := range bookSeats {
		message := fmt.Sprintf("Подтвердите или отмените свое бронирование на сегодня до 10:00")
		data := form.InfoFormData{
			Action:     dto.ActionShowOfficeMenu,
			Message:    message,
			BookSeatID: bookSeat.ID,
			ChatID:     bookSeat.User.ChatID,
		}
		arrayToSend = append(arrayToSend, data)
	}

	if len(arrayToSend) > 0 {
		i.sender.AddToQueue(arrayToSend)
	}
	return nil
}

// SendNotifyTomorrowBookingOpen Сообщение подписавшимся, что открыта запись на завтра

func (i *informerServiceImpl) SendNotifyTomorrowBookingOpen(ctx context.Context, office model.Office, message string) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	users, err := i.userRepo.GetUsersToNotify(ctx, office.ID)
	if err != nil {
		return err
	}

	var arrayToSend []form.InfoFormData

	for _, user := range users {
		data := form.InfoFormData{
			Action:     dto.ActionShowOfficeMenu,
			Message:    message,
			BookSeatID: 0,
			ChatID:     user.ChatID,
		}
		arrayToSend = append(arrayToSend, data)
	}

	if len(arrayToSend) > 0 {
		i.sender.AddToQueue(arrayToSend)
	}
	return nil
}

func (i *informerServiceImpl) SendNotifyToBookDeletedBySystem(ctx context.Context, bookSeats []*model.BookSeat, officeName string) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	var arrayToSend []form.InfoFormData

	for _, bookSeat := range bookSeats {
		formattedDate := bookSeat.BookDate.Format(helper.DateFormat)
		message := fmt.Sprintf("Мы удалили вашю бронь в хотдеске %s на %s, так как вы ее не подтвердили", officeName, formattedDate)

		data := form.InfoFormData{
			Action:     dto.ActionShowOfficeMenu,
			Message:    message,
			BookSeatID: bookSeat.ID,
			ChatID:     bookSeat.User.ChatID,
		}
		arrayToSend = append(arrayToSend, data)
	}

	if len(arrayToSend) > 0 {
		i.sender.AddToQueue(arrayToSend)
	}

	return nil
}
