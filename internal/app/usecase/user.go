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
)

type UserService interface {
	FirstCome(ctx context.Context) (*dto.UserResult, error)
	CallChooseOfficeMenu(ctx context.Context) (*dto.UserResult, error)
	SetOfficeScript(ctx context.Context, officeID int64) (*dto.UserResult, error)
	CallSeatsMenu(ctx context.Context) (*dto.UserResult, error)
	BookSeatTap(ctx context.Context, bookSeatID int64) (*dto.UserResult, error)
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

func (u *userServiceImpl) BookSeatTap(ctx context.Context, bookSeatID int64) (*dto.UserResult, error) {

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
			message = "Место уже занято другим юзером"
		}
	} else {
		// место свободно
		answerType = SeatFree
		message = fmt.Sprintf("Чтобы занять место №%d, укажите время:", bookSeat.Seat.SeatNumber)
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
