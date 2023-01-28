package usecase

import (
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

const (
	ChooseOffice  = "choose_office"
	ConfirmOffice = "confirm_office"
)

type UserService interface {
	FirstCome(data dto.UserLogicData) (*dto.UserLogicResult, error)
	SetOffice(data dto.SetOfficeDTO) error
}

type userServiceImpl struct {
	userRepo   interfaces.UserRepository
	officeRepo interfaces.OfficeRepository
	logger     *zap.Logger
}

func NewUserService(userRepo interfaces.UserRepository,
	officeRepo interfaces.OfficeRepository,
	logger *zap.Logger) UserService {
	return &userServiceImpl{
		userRepo:   userRepo,
		officeRepo: officeRepo,
		logger:     logger,
	}
}

func (u *userServiceImpl) FirstCome(data dto.UserLogicData) (*dto.UserLogicResult, error) {

	user, err := u.userRepo.GetByTelegramID(data.User.TelegramID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		err = u.userRepo.Create(data.User)
		if err != nil {
			return nil, err
		}
		user, err = u.userRepo.GetByTelegramID(data.User.TelegramID)
		if err != nil {
			return nil, err
		}
	}

	data.User = *user
	if user.HaveChosenOffice() {
		return u.confirmAlreadyChosenOffice(data)
	} else {
		return u.chooseOffice(data)
	}
}

func (u *userServiceImpl) confirmAlreadyChosenOffice(data dto.UserLogicData) (*dto.UserLogicResult, error) {

	office, err := u.officeRepo.FindByID(data.User.OfficeID)
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("%s, хотите занять место в: %s?", data.User.Name, office.Name)
	return &dto.UserLogicResult{
		Key:       ConfirmOffice,
		Office:    office,
		Offices:   nil,
		Message:   message,
		User:      data.User,
		ChatID:    data.ChatID,
		MessageID: data.MessageID,
	}, nil
}

func (u *userServiceImpl) chooseOffice(data dto.UserLogicData) (*dto.UserLogicResult, error) {

	offices, err := u.officeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("Привет, %s! Давай выберем офис)", data.User.Name)
	return &dto.UserLogicResult{
		Key:       ChooseOffice,
		Office:    nil,
		Offices:   offices,
		Message:   message,
		User:      data.User,
		ChatID:    data.ChatID,
		MessageID: data.MessageID,
	}, nil
}

func (u *userServiceImpl) SetOffice(data dto.SetOfficeDTO) error {
	user, err := u.userRepo.GetByTelegramID(data.TelegramID)
	if err != nil {
		return err
	}
	if user == nil {
		// TODO добавить ошибку NotFoundUser
		return nil
	}
	user.OfficeID = data.OfficeID
	err = u.userRepo.SetOffice(user)
	if err != nil {
		return err
	}
	return nil
}
