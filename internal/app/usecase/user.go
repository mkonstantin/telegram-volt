package usecase

import (
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/common"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

const (
	ChooseOfficeMenu = "choose_office"
	OfficeMenu       = "office_menu"
	ChooseSeatsMenu  = "choose_seats_menu"
)

type UserService interface {
	FirstCome(data dto.FirstStartDTO) (*dto.UserResult, error)
	CallChooseOfficeMenu(data dto.FirstStartDTO) (*dto.UserResult, error)
	SetOfficeScript(data dto.OfficeChosenDTO) (*dto.UserResult, error)
	CallSeatsMenu(data dto.BookSeatDTO) (*dto.UserResult, error)
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

func (u *userServiceImpl) FirstCome(data dto.FirstStartDTO) (*dto.UserResult, error) {

	user, err := u.userRepo.GetByTelegramID(data.User.TelegramID)

	if err == common.ErrUserNotFound {
		user, err = u.createUser(data)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	data.User = *user
	if user.HaveChosenOffice() {
		return u.callOfficeMenu(data.User.OfficeID, data.ChatID, data.MessageID)
	} else {
		return u.CallChooseOfficeMenu(data)
	}
}

func (u *userServiceImpl) createUser(data dto.FirstStartDTO) (*model.User, error) {
	err := u.userRepo.Create(data.User)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetByTelegramID(data.User.TelegramID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userServiceImpl) callOfficeMenu(officeID, chatID int64, MessageID int) (*dto.UserResult, error) {

	office, err := u.officeRepo.FindByID(officeID)
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("Офис: %s, действия:", office.Name)
	return &dto.UserResult{
		Key:       OfficeMenu,
		Office:    office,
		Offices:   nil,
		Message:   message,
		ChatID:    chatID,
		MessageID: MessageID,
	}, nil
}

func (u *userServiceImpl) CallChooseOfficeMenu(data dto.FirstStartDTO) (*dto.UserResult, error) {

	offices, err := u.officeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("%s, давай выберем офис:", data.User.Name)
	return &dto.UserResult{
		Key:       ChooseOfficeMenu,
		Office:    nil,
		Offices:   offices,
		Message:   message,
		ChatID:    data.ChatID,
		MessageID: data.MessageID,
	}, nil
}

//========= Выбрали офис и вызываем его меню

func (u *userServiceImpl) SetOfficeScript(data dto.OfficeChosenDTO) (*dto.UserResult, error) {

	user, err := u.userRepo.GetByTelegramID(data.TelegramID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	user.OfficeID = data.OfficeID

	err = u.userRepo.SetOffice(user)
	if err != nil {
		return nil, err
	}

	return u.callOfficeMenu(user.OfficeID, data.ChatID, data.MessageID)
}

//========= Места в офисе

func (u *userServiceImpl) CallSeatsMenu(data dto.BookSeatDTO) (*dto.UserResult, error) {

	seats, err := u.bookSeatRepo.GetAllByOfficeID(data.OfficeID)
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
		ChatID:    data.ChatID,
		MessageID: data.MessageID,
	}, nil
}
