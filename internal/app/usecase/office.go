package usecase

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type OfficeService interface {
	ChooseOffice(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type officeServiceImpl struct {
	userRepo   interfaces.UserRepository
	officeRepo interfaces.OfficeRepository
	logger     *zap.Logger
}

func NewOfficeService(userRepo interfaces.UserRepository,
	officeRepo interfaces.OfficeRepository,
	logger *zap.Logger) OfficeService {
	return &officeServiceImpl{
		userRepo:   userRepo,
		officeRepo: officeRepo,
		logger:     logger,
	}
}

func (o *officeServiceImpl) ChooseOffice(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	//user, err := o.userRepo.GetByTelegramID(update.Message.From.ID)
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}
