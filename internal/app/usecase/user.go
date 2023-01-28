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
	ChooseOffice  = "choose_office"
	ConfirmOffice = "confirm_office"
)

type UserService interface {
	FirstCome(data dto.FirstStartDTO) (*dto.FirstStartResult, error)
	OfficeChosenScenery(data dto.OfficeChosenDTO) error
}

type userServiceImpl struct {
	userRepo   interfaces.UserRepository
	officeRepo interfaces.OfficeRepository
	seatRepo   interfaces.SeatRepository
	logger     *zap.Logger
}

func NewUserService(userRepo interfaces.UserRepository,
	officeRepo interfaces.OfficeRepository, seatRepo interfaces.SeatRepository,
	logger *zap.Logger) UserService {
	return &userServiceImpl{
		userRepo:   userRepo,
		officeRepo: officeRepo,
		seatRepo:   seatRepo,
		logger:     logger,
	}
}

func (u *userServiceImpl) FirstCome(data dto.FirstStartDTO) (*dto.FirstStartResult, error) {

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
		return u.confirmAlreadyChosenOffice(data)
	} else {
		return u.chooseOffice(data)
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

func (u *userServiceImpl) confirmAlreadyChosenOffice(data dto.FirstStartDTO) (*dto.FirstStartResult, error) {

	office, err := u.officeRepo.FindByID(data.User.OfficeID)
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("%s, хотите занять место в: %s?", data.User.Name, office.Name)
	return &dto.FirstStartResult{
		Key:       ConfirmOffice,
		Office:    office,
		Offices:   nil,
		Message:   message,
		User:      data.User,
		ChatID:    data.ChatID,
		MessageID: data.MessageID,
	}, nil
}

func (u *userServiceImpl) chooseOffice(data dto.FirstStartDTO) (*dto.FirstStartResult, error) {

	offices, err := u.officeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("Привет, %s! Давай выберем офис)", data.User.Name)
	return &dto.FirstStartResult{
		Key:       ChooseOffice,
		Office:    nil,
		Offices:   offices,
		Message:   message,
		User:      data.User,
		ChatID:    data.ChatID,
		MessageID: data.MessageID,
	}, nil
}

// Office выбран, теперь надо выбрать место

func (u *userServiceImpl) OfficeChosenScenery(data dto.OfficeChosenDTO) error {
	user, err := u.userRepo.GetByTelegramID(data.TelegramID)
	if err != nil {
		return err
	}
	if user == nil {
		return common.ErrUserNotFound
	}

	user.OfficeID = data.OfficeID

	err = u.userRepo.SetOffice(user)
	if err != nil {
		return err
	}

	return nil
}
