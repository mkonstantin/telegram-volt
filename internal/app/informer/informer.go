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
)

type InformerService interface {
	SendNotifySeatBecomeFree(ctx context.Context, id int64) error
	SendNotifyTomorrowBookingOpen(office model.Office, message string) error
	SendNotifiesToConfirm(office *model.Office) error
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

// SendNotifySeatBecomeFree Сообщение подписавшимся кроме тех кто уже занимает место, что место освободилось

func (i *informerServiceImpl) SendNotifySeatBecomeFree(ctx context.Context, bookSeatID int64) error {

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

	seats, err := i.bookSeatRepo.FindByOfficeIDAndDate(bookSeat.Office.ID, bookSeat.BookDate.String())
	var mapper = make(map[int64]int)
	mapper[currentUser.ID]++
	for _, seat := range seats {
		mapper[seat.User.ID]++
	}

	data := form.InfoFormData{
		Action:     dto.ActionShowSeatList,
		Message:    text,
		BookSeatID: bookSeat.ID,
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

// SendNotifiesToConfirm рассылка уведомлений на подтверждение брони

func (i *informerServiceImpl) SendNotifiesToConfirm(office *model.Office) error {

	today := helper.TodayZeroTimeUTC()
	bookSeats, err := i.bookSeatRepo.FindNotConfirmedByOfficeIDAndDate(office.ID, today.String())
	if err != nil {
		return err
	}

	for _, bookSeat := range bookSeats {
		message := fmt.Sprintf("Подтвердите свою бронь на сегодня, иначе мы УДАЛИМ ее через час")
		data := form.InfoFormData{
			Action:     dto.ActionShowOfficeMenu,
			Message:    message,
			BookSeatID: bookSeat.ID,
			ChatID:     bookSeat.User.ChatID,
		}

		i.sendInfoForm(context.Background(), data)
	}

	return nil
}

// SendNotifyTomorrowBookingOpen Сообщение подписавшимся, что открыта запись на завтра

func (i *informerServiceImpl) SendNotifyTomorrowBookingOpen(office model.Office, message string) error {
	users, err := i.userRepo.GetUsersToNotify(office.ID)
	if err != nil {
		return err
	}

	for _, user := range users {
		data := form.InfoFormData{
			Action:     dto.ActionShowOfficeMenu,
			Message:    message,
			BookSeatID: 0,
			ChatID:     user.ChatID,
		}
		i.sendInfoForm(context.Background(), data)
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
