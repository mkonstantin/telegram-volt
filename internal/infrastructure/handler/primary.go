package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type UserHundler struct {
	userRepo interfaces.UserRepository
	logger   *zap.Logger
}

func NewPrimaryHundler(repo interfaces.UserRepository, logger *zap.Logger) UserHundler {
	return UserHundler{
		userRepo: repo,
		logger:   logger,
	}
}

const (
	Yakutsk203    = "Якутск, 203"
	YakutskGluhoi = "Якутск, пер. Глухой"
	Moscow        = "Москва"
	Almaty        = "Алматы"
)

var firstOfficeChoose = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Almaty),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Moscow),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Yakutsk203),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(YakutskGluhoi),
	),
)

var choosePlaceOrChangeOffice = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Я пойду на обед"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Поменять офис"),
	),
)

func (s *UserHundler) Start(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	user, err := s.userRepo.GetByTelegramID(update.Message.From.ID)

	if err != nil {
		s.logger.Warn("UserHundler cant get telegram id", zap.Error(err))
		return tgbotapi.MessageConfig{}, err
	}
	var msgConfig tgbotapi.MessageConfig

	if user.IsOfficeChoosed() {
		officeName := user.OfficeName()
		strName := fmt.Sprintf("Ваш офис: %s", officeName)
		msgConfig = tgbotapi.NewMessage(update.Message.Chat.ID, strName)
		msgConfig.ReplyMarkup = choosePlaceOrChangeOffice
		msgConfig.ReplyToMessageID = update.Message.MessageID
	} else {
		name := update.Message.From.FirstName
		strName := fmt.Sprintf("Привет, %s! Для начала давай выберем офис)", name)
		msgConfig = tgbotapi.NewMessage(update.Message.Chat.ID, strName)
		msgConfig.ReplyMarkup = firstOfficeChoose
		msgConfig.ReplyToMessageID = update.Message.MessageID
	}
	return msgConfig, nil
}
