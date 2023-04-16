package usecase

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type UserService interface {
	SetOfficeScript(ctx context.Context, officeID int64) (context.Context, error)
	SubscribeWork(ctx context.Context) (context.Context, string, error)
	BookSeat(ctx context.Context, bookSeatID int64) (string, error)
	CancelBookSeat(ctx context.Context, bookSeatID int64) (string, bool, error)
	ConfirmBookSeat(ctx context.Context, bookSeatID int64) (string, error)
}

type userServiceImpl struct {
	userRepo     interfaces.UserRepository
	officeRepo   interfaces.OfficeRepository
	bookSeatRepo interfaces.BookSeatRepository
	logger       *zap.Logger
}

func NewUserService(userRepo interfaces.UserRepository,
	officeRepo interfaces.OfficeRepository, bookSeatRepo interfaces.BookSeatRepository,
	logger *zap.Logger) UserService {
	return &userServiceImpl{
		userRepo:     userRepo,
		officeRepo:   officeRepo,
		bookSeatRepo: bookSeatRepo,
		logger:       logger,
	}
}

//========= Выбрали офис и вызываем его меню

func (u *userServiceImpl) SetOfficeScript(ctx context.Context, officeID int64) (context.Context, error) {
	currentUser := model.GetCurrentUser(ctx)

	err := u.userRepo.SetOffice(officeID, currentUser.TelegramID)
	if err != nil {
		return ctx, err
	}

	currentUser.OfficeID = officeID
	ctx = context.WithValue(ctx, model.ContextUserKey, currentUser)

	return ctx, nil
}

// ========== Забронировали место

func (u *userServiceImpl) BookSeat(ctx context.Context, bookSeatID int64) (string, error) {
	var message string
	currentUser := model.GetCurrentUser(ctx)

	bookSeat, err := u.bookSeatRepo.FindByID(bookSeatID)
	if err != nil {
		return "", err
	}

	userBookSeat, err := u.bookSeatRepo.FindByUserIDAndDate(currentUser.ID, bookSeat.BookDate.String())
	if err != nil {
		return "", err
	}
	if userBookSeat != nil {
		message = "У вас уже есть бронь в этом офисе на эту дату"
		return message, nil
	}

	if bookSeat.User != nil {
		message = fmt.Sprintf("Место №%s уже занято", bookSeat.Seat.SeatSign)
	} else {
		today := helper.TodayZeroTimeUTC()
		isToday := false
		if bookSeat.BookDate.Equal(today) {
			isToday = true
		}
		err = u.bookSeatRepo.BookSeatWithID(currentUser.ID, bookSeatID, isToday)
		if err != nil {
			return "", err
		}

		if isToday {
			message = fmt.Sprintf("Отлично! Ваше место №%s в офисе %s забронировано!",
				bookSeat.Seat.SeatSign, bookSeat.Office.Name)
		} else {
			message = fmt.Sprintf("Ваше место №%s в офисе %s забронировано. "+
				"Завтра в 9:00 откроется возможность подтверждения бронирования, "+
				"если вы не подтвердите его до 10:00, бронь будет аннулирована",
				bookSeat.Seat.SeatSign, bookSeat.Office.Name)
		}
	}

	return message, nil
}

// ========== Отменили бронирование места

func (u *userServiceImpl) CancelBookSeat(ctx context.Context, bookSeatID int64) (string, bool, error) {
	var message string
	currentUser := model.GetCurrentUser(ctx)

	bookSeat, err := u.bookSeatRepo.FindByID(bookSeatID)
	if err != nil {
		return "", false, err
	}

	userBookSeat, err := u.bookSeatRepo.FindByUserIDAndDate(currentUser.ID, bookSeat.BookDate.String())
	if err != nil {
		return "", false, err
	}
	if userBookSeat == nil || (userBookSeat != nil && userBookSeat.User.ID != currentUser.ID) {
		message = "У вас нет брони на эту дату"
		return message, false, nil
	}

	err = u.bookSeatRepo.CancelBookSeatWithID(bookSeatID)
	if err != nil {
		return "", false, err
	}
	message = fmt.Sprintf("Место №%s в офисе: %s освобождено. Спасибо!", userBookSeat.Seat.SeatSign, userBookSeat.Office.Name)

	return message, true, nil
}

// ========== Подписка/отписка на свободные места

func (u *userServiceImpl) SubscribeWork(ctx context.Context) (context.Context, string, error) {
	var message string
	currentUser := model.GetCurrentUser(ctx)

	if currentUser.OfficeID == 0 {
		return ctx, "Произошла ошибка: необходимо выбрать офис", nil
	}

	office, err := u.officeRepo.FindByID(currentUser.OfficeID)
	if err != nil {
		return ctx, "", err
	}

	if currentUser.NotifyOfficeID == currentUser.OfficeID {
		err = u.userRepo.Unsubscribe(currentUser.TelegramID)
		if err != nil {
			return ctx, "", err
		}

		currentUser.NotifyOfficeID = 0
		message = fmt.Sprintf("Вы отменили подписку на свободные места в офисе: %s", office.Name)
	} else {
		err = u.userRepo.Subscribe(currentUser.TelegramID, currentUser.OfficeID)
		if err != nil {
			return ctx, "", err
		}

		currentUser.NotifyOfficeID = currentUser.OfficeID
		message = fmt.Sprintf("Вы подписались на свободные места в офисе: %s", office.Name)
	}

	ctx = context.WithValue(ctx, model.ContextUserKey, currentUser)
	return ctx, message, nil
}

func (u *userServiceImpl) ConfirmBookSeat(ctx context.Context, bookSeatID int64) (string, error) {
	var message string

	bookSeat, err := u.bookSeatRepo.FindByID(bookSeatID)
	if err != nil {
		return "", err
	}

	err = u.bookSeatRepo.ConfirmBookSeat(bookSeatID)
	if err != nil {
		return "", err
	}

	formattedDate := bookSeat.BookDate.Format(helper.DateFormat)
	message = fmt.Sprintf("Отлично! Вы подтвердили, что придете сегодня: %s. "+
		"Ваше место №%s в офисе: %s", formattedDate, bookSeat.Seat.SeatSign, bookSeat.Office.Name)

	return message, nil
}
