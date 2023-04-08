package usecase

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type UserService interface {
	SetOfficeScript(ctx context.Context, officeID int64) (context.Context, error)
	SubscribeWork(ctx context.Context) (string, error)
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
		message = fmt.Sprintf("Место №%d уже занято", bookSeat.Seat.SeatNumber)
	} else {
		err = u.bookSeatRepo.BookSeatWithID(currentUser.ID, bookSeatID)
		if err != nil {
			return "", err
		}

		message = fmt.Sprintf("Отлично! Вы заняли место №%d в офисе: %s", bookSeat.Seat.SeatNumber, bookSeat.Office.Name)
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
	message = fmt.Sprintf("Место №%d в офисе: %s освобождено. Спасибо!", userBookSeat.Seat.SeatNumber, userBookSeat.Office.Name)

	return message, true, nil
}

// ========== Подписка/отписка на свободные места

func (u *userServiceImpl) SubscribeWork(ctx context.Context) (string, error) {
	var message string
	currentUser := model.GetCurrentUser(ctx)

	if currentUser.OfficeID == 0 {
		return "Произошла ошибка: необходимо выбрать офис", nil
	}

	if currentUser.NotifyOfficeID == currentUser.OfficeID {
		err := u.userRepo.Unsubscribe(currentUser.TelegramID)
		if err != nil {
			return "", err
		}

		office, err := u.officeRepo.FindByID(currentUser.NotifyOfficeID)
		if err != nil {
			return "", err
		}

		message = fmt.Sprintf("Вы отменили подписку на свободные места в офисе: %s", office.Name)
	} else {
		err := u.userRepo.Subscribe(currentUser.TelegramID, currentUser.OfficeID)
		if err != nil {
			return "", err
		}

		office, err := u.officeRepo.FindByID(currentUser.OfficeID)
		if err != nil {
			return "", err
		}

		message = fmt.Sprintf("Вы подписались на свободные места в офисе: %s", office.Name)
	}

	return message, nil
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

	message = fmt.Sprintf("Отлично! Вы Подтвердили, что придете сегодня. "+
		"Ваше место №%d в офисе: %s", bookSeat.Seat.SeatNumber, bookSeat.Office.Name)

	return message, nil
}
