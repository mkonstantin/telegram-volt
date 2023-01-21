package use_case

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type UserService interface {
	FirstCome(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type userServiceImpl struct {
	userRepo interfaces.UserRepository
	logger   *zap.Logger
}

func NewUserService(userRepo interfaces.UserRepository, logger *zap.Logger) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u userServiceImpl) FirstCome(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	//TODO implement me
	panic("implement me")
}
