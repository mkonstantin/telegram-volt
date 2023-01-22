package use_case

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type UserService interface {
	FirstCome(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
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

var confirmOfficeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "yes"),
		tgbotapi.NewInlineKeyboardButtonData("Нет", "no"),
	),
)

func (u *userServiceImpl) FirstCome(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {

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
		return u.confirmAlreadyChosenOffice(update.Message.MessageID, update.Message.Chat.ID, user)
	} else {
		return u.chooseOffice(update.Message.MessageID, update.Message.Chat.ID, user)
	}
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

func (u *userServiceImpl) confirmAlreadyChosenOffice(messageID int, chatID int64, user *model.User) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(chatID, "")

	office, err := u.officeRepo.FindByID(user.OfficeID)
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("%s, хотите занять место в: %s?", user.Name, office.Name)
	msg.Text = message
	msg.ReplyMarkup = confirmOfficeKeyboard
	msg.ReplyToMessageID = messageID
	return &msg, nil
}

func (u *userServiceImpl) chooseOffice(messageID int, chatID int64, user *model.User) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(chatID, "")

	offices, err := u.officeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, office := range offices {
		button := tgbotapi.NewInlineKeyboardButtonData(office.Name, strconv.FormatInt(office.ID, 10))
		row := tgbotapi.NewInlineKeyboardRow(button)
		rows = append(rows, row)
	}

	var chooseOfficeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	)

	fmt.Println(offices)
	message := fmt.Sprintf("Привет, %s! Давай выберем офис)", user.Name)
	msg.Text = message
	msg.ReplyMarkup = chooseOfficeKeyboard
	msg.ReplyToMessageID = messageID

	return &msg, nil
}
