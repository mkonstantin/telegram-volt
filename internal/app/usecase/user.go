package usecase

import (
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

const (
	ChooseOffice  = "choose_office"
	ConfirmOffice = "confirm_office"
)

type UserLogicData struct {
	User      model.User
	MessageID int
	ChatID    int64
}

type UserLogicResult struct {
	Key       string
	Office    *model.Office
	Offices   []*model.Office
	Message   string
	User      model.User
	MessageID int
	ChatID    int64
}

type UserService interface {
	FirstCome(data UserLogicData) (*UserLogicResult, error)
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

func (u *userServiceImpl) FirstCome(data UserLogicData) (*UserLogicResult, error) {

	user, err := u.userRepo.GetByTelegramID(data.User.TelegramID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		err := u.userRepo.Create(data.User)
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

func (u *userServiceImpl) confirmAlreadyChosenOffice(data UserLogicData) (*UserLogicResult, error) {

	office, err := u.officeRepo.FindByID(data.User.OfficeID)
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("%s, хотите занять место в: %s?", data.User.Name, office.Name)
	return &UserLogicResult{
		Key:       ConfirmOffice,
		Office:    office,
		Offices:   nil,
		Message:   message,
		User:      data.User,
		ChatID:    data.ChatID,
		MessageID: data.MessageID,
	}, nil
}

func (u *userServiceImpl) chooseOffice(data UserLogicData) (*UserLogicResult, error) {

	offices, err := u.officeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("Привет, %s! Давай выберем офис)", data.User.Name)
	return &UserLogicResult{
		Key:       ChooseOffice,
		Office:    nil,
		Offices:   offices,
		Message:   message,
		User:      data.User,
		ChatID:    data.ChatID,
		MessageID: data.MessageID,
	}, nil
}
