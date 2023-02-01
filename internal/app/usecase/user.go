package usecase

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

const (
	ChooseOfficeMenu = "choose_office"
	OfficeMenu       = "office_menu"
	ChooseSeatsMenu  = "choose_seats_menu"
	SeatOwn          = "seat_own"
	SeatBusy         = "seat_busy"
	SeatFree         = "seat_free"
	BookSeat         = "book_seat"
)

type UserService interface {
	FirstCome(ctx context.Context) (*dto.UserResult, error)
	CallChooseOfficeMenu(ctx context.Context) (*dto.UserResult, error)
	SetOfficeScript(ctx context.Context, officeID int64) (*dto.UserResult, error)
	CallSeatsMenu(ctx context.Context) (*dto.UserResult, error)
	SeatListTap(ctx context.Context, bookSeatID int64) (*dto.UserResult, error)
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

func (u *userServiceImpl) FirstCome(ctx context.Context) (*dto.UserResult, error) {

	currentUser := model.GetCurrentUser(ctx)

	if currentUser.HaveChosenOffice() {
		return u.callOfficeMenu(ctx)
	} else {
		return u.CallChooseOfficeMenu(ctx)
	}
}

func (u *userServiceImpl) callOfficeMenu(ctx context.Context) (*dto.UserResult, error) {

	currentUser := model.GetCurrentUser(ctx)

	office, err := u.officeRepo.FindByID(currentUser.OfficeID)
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("Офис: %s, действия:", office.Name)
	return &dto.UserResult{
		Key:     OfficeMenu,
		Office:  office,
		Offices: nil,
		Message: message,
	}, nil
}

func (u *userServiceImpl) CallChooseOfficeMenu(ctx context.Context) (*dto.UserResult, error) {

	currentUser := model.GetCurrentUser(ctx)

	offices, err := u.officeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("%s, давай выберем офис:", currentUser.Name)
	return &dto.UserResult{
		Key:     ChooseOfficeMenu,
		Office:  nil,
		Offices: offices,
		Message: message,
	}, nil
}

//========= Выбрали офис и вызываем его меню

func (u *userServiceImpl) SetOfficeScript(ctx context.Context, officeID int64) (*dto.UserResult, error) {

	currentUser := model.GetCurrentUser(ctx)

	err := u.userRepo.SetOffice(officeID, currentUser.TelegramID)
	if err != nil {
		return nil, err
	}

	currentUser.OfficeID = officeID
	ctx = context.WithValue(ctx, model.ContextUserKey, currentUser)

	return u.callOfficeMenu(ctx)
}

//========= Места в офисе

func (u *userServiceImpl) CallSeatsMenu(ctx context.Context) (*dto.UserResult, error) {

	currentUser := model.GetCurrentUser(ctx)

	seats, err := u.bookSeatRepo.GetAllByOfficeID(currentUser.OfficeID)
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Выберите место:")

	return &dto.UserResult{
		Key:       ChooseOfficeMenu,
		Office:    nil,
		Offices:   nil,
		BookSeats: seats,
		Message:   message,
	}, nil
}

// ========== Выбрали место в списке

func (u *userServiceImpl) SeatListTap(ctx context.Context, bookSeatID int64) (*dto.UserResult, error) {

	currentUser := model.GetCurrentUser(ctx)

	bookSeat, err := u.bookSeatRepo.FindByID(bookSeatID)
	if err != nil {
		return nil, err
	}

	var answerType string
	var message string
	if bookSeat.User != nil {
		if bookSeat.User.TelegramID == currentUser.TelegramID {
			// место уже занято самим же юзером
			answerType = SeatOwn
			message = "Вы уже заняли это место, хотите его освободить?"
		} else {
			// место занято другим юзером
			answerType = SeatBusy
			message = fmt.Sprintf("Место №%d уже занято %s aka @%s",
				bookSeat.Seat.SeatNumber, bookSeat.User.Name, bookSeat.User.TelegramName)
		}
	} else {
		// место свободно
		answerType = SeatFree
		message = fmt.Sprintf("Занять место №%d?", bookSeat.Seat.SeatNumber)
	}

	return &dto.UserResult{
		Key:        answerType,
		Office:     nil,
		Offices:    nil,
		BookSeats:  nil,
		BookSeatID: bookSeat.ID,
		Message:    message,
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
	message = fmt.Sprintf("Все в порядке,т место №%d в офисе: %s освобождено. Спасибо!", userBookSeat.Seat.SeatNumber, userBookSeat.Office.Name)

	return &dto.UserResult{
		Key:        BookSeat,
		Office:     nil,
		Offices:    nil,
		BookSeats:  nil,
		BookSeatID: bookSeatID,
		Message:    message,
	}, nil
}
