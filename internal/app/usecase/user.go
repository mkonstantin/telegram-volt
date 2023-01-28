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
)

type UserService interface {
	FirstCome(data dto.FirstStartDTO) (*dto.UserResult, error)
	OfficeChosenScenery(data dto.OfficeChosenDTO) (*dto.UserResult, error)
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
		return u.callChooseOfficeMenu(data.User.Name, data.ChatID, data.MessageID)
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
	message := fmt.Sprintf("Офис: %s, выберите действие:", office.Name)
	return &dto.UserResult{
		Key:       OfficeMenu,
		Office:    office,
		Offices:   nil,
		Message:   message,
		ChatID:    chatID,
		MessageID: MessageID,
	}, nil
}

func (u *userServiceImpl) callChooseOfficeMenu(name string, chatID int64, messageID int) (*dto.UserResult, error) {

	offices, err := u.officeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("Привет, %s! Давай выберем офис)", name)
	return &dto.UserResult{
		Key:       ChooseOfficeMenu,
		Office:    nil,
		Offices:   offices,
		Message:   message,
		ChatID:    chatID,
		MessageID: messageID,
	}, nil
}

//========= Office выбран, теперь надо выбрать место

func (u *userServiceImpl) OfficeChosenScenery(data dto.OfficeChosenDTO) (*dto.UserResult, error) {
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
