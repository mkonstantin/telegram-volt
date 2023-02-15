package usecase

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/internal/infrastructure/service"
)

const (
	DateMenu = "date_menu"
	BookSeat = "book_seat"
)

type UserService interface {
	SetOfficeScript(ctx context.Context, officeID int64) (context.Context, error)
	SubscribeWork(ctx context.Context) (string, error)

	CallDateMenu(ctx context.Context) (*dto.UserResult, error)
	CallSeatsMenu(ctx context.Context) (*dto.UserResult, error)
	BookSeat(ctx context.Context, bookSeatID int64) (*dto.UserResult, error)
	CancelBookSeat(ctx context.Context, bookSeatID int64) (*dto.UserResult, error)
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

//=========  Выбираем дату:

func (u *userServiceImpl) CallDateMenu(ctx context.Context) (*dto.UserResult, error) {
	currentUser := model.GetCurrentUser(ctx)

	office, err := u.officeRepo.FindByID(currentUser.OfficeID)
	if err != nil {
		return nil, err
	}

	today := service.TodayZeroTimeUTC()
	todaySeats, err := u.bookSeatRepo.GetAllByOfficeIDAndDate(currentUser.OfficeID, today.String())
	if err != nil {
		return nil, err
	}

	tomorrow := service.TomorrowZeroTimeUTC()
	tomorrowSeats, err := u.bookSeatRepo.GetAllByOfficeIDAndDate(currentUser.OfficeID, tomorrow.String())
	if err != nil {
		return nil, err
	}

	todayD := dto.DaySeat{
		Date:        today.String(),
		SeatsNumber: len(todaySeats),
	}

	tomorrowD := dto.DaySeat{
		Date:        tomorrow.String(),
		SeatsNumber: len(tomorrowSeats),
	}

	var seatByDates []dto.DaySeat
	seatByDates = append(seatByDates, todayD)
	seatByDates = append(seatByDates, tomorrowD)

	message := fmt.Sprintf("Выберите дату:")

	return &dto.UserResult{
		Key:         DateMenu,
		Message:     message,
		Office:      office,
		SeatByDates: seatByDates,
	}, nil
}

//========= Места в офисе

func (u *userServiceImpl) CallSeatsMenu(ctx context.Context) (*dto.UserResult, error) {

	currentUser := model.GetCurrentUser(ctx)

	date := service.TodayZeroTimeUTC()

	seats, err := u.bookSeatRepo.GetAllByOfficeIDAndDate(currentUser.OfficeID, date.String())
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Выберите место:")

	return &dto.UserResult{
		Key:       "",
		Office:    nil,
		Offices:   nil,
		BookSeats: seats,
		Message:   message,
	}, nil
}

// ========== Забронировали место

func (u *userServiceImpl) BookSeat(ctx context.Context, bookSeatID int64) (*dto.UserResult, error) {

	var message string
	currentUser := model.GetCurrentUser(ctx)

	userBookSeat, err := u.bookSeatRepo.FindByUserID(currentUser.ID)
	if err != nil {
		return nil, err
	}
	if userBookSeat != nil {
		message = "У вас уже есть бронь в этом офисе на сегодня"
		return &dto.UserResult{
			Key:        BookSeat,
			Office:     nil,
			Offices:    nil,
			BookSeats:  nil,
			BookSeatID: bookSeatID,
			Message:    message,
		}, nil
	}

	bookSeat, err := u.bookSeatRepo.FindByID(bookSeatID)
	if err != nil {
		return nil, err
	}

	if bookSeat.User != nil {
		message = fmt.Sprintf("Место №%d уже занято", bookSeat.Seat.SeatNumber)
	} else {
		err = u.bookSeatRepo.BookSeatWithID(currentUser.ID, bookSeatID)
		if err != nil {
			return nil, err
		}

		message = fmt.Sprintf("Отлично! Вы заняли место №%d в офисе: %s", bookSeat.Seat.SeatNumber, bookSeat.Office.Name)
	}

	return &dto.UserResult{
		Key:        BookSeat,
		Office:     nil,
		Offices:    nil,
		BookSeats:  nil,
		BookSeatID: bookSeatID,
		Message:    message,
	}, nil
}

// ========== Отменили бронирование места

func (u *userServiceImpl) CancelBookSeat(ctx context.Context, bookSeatID int64) (*dto.UserResult, error) {
	var message string
	currentUser := model.GetCurrentUser(ctx)

	userBookSeat, err := u.bookSeatRepo.FindByUserID(currentUser.ID)
	if err != nil {
		return nil, err
	}
	if userBookSeat == nil || (userBookSeat != nil && userBookSeat.User.ID != currentUser.ID) {
		message = "У вас нет брони на сегодня"
		return &dto.UserResult{
			Key:        BookSeat,
			Office:     nil,
			Offices:    nil,
			BookSeats:  nil,
			BookSeatID: bookSeatID,
			Message:    message,
		}, nil
	}

	err = u.bookSeatRepo.CancelBookSeatWithID(bookSeatID)
	if err != nil {
		return nil, err
	}
	message = fmt.Sprintf("Место №%d в офисе: %s освобождено. Спасибо!", userBookSeat.Seat.SeatNumber, userBookSeat.Office.Name)

	return &dto.UserResult{
		Key:        BookSeat,
		Office:     nil,
		Offices:    nil,
		BookSeats:  nil,
		BookSeatID: bookSeatID,
		Message:    message,
	}, nil
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
