package use_case

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/domain/model"
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

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

func (u *userServiceImpl) FirstCome(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	user, err := u.userRepo.GetByTelegramID(update.Message.From.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		err = u.saveUser(update.Message.From)
		if err != nil {
			return nil, err
		}
	}

	if user.HaveChosenOffice() {

	}
	message := fmt.Sprintf("Привет, %s! Для начала давай выберем офис)", user.Name)
	msg.Text = message
	msg.ReplyMarkup = numericKeyboard
	//msg.ReplyToMessageID = update.Message.MessageID

	return &msg, nil
}

func (u *userServiceImpl) saveUser(TGUser *tgbotapi.User) error {
	userModel := model.User{
		Name:         TGUser.FirstName,
		TelegramID:   TGUser.ID,
		TelegramName: TGUser.UserName,
	}
	err := u.userRepo.Create(userModel)
	if err != nil {
		return err
	}
	return nil
}
