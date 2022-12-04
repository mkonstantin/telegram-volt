package router

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure_layer/hundlers"
)

type Router struct {
	primaryHundler hundlers.UserHundler
	officeHundler  hundlers.OfficeHundler
	logger         *zap.Logger
}

func NewRouter(officeHundler hundlers.OfficeHundler,
	primaryHundler hundlers.UserHundler, logger *zap.Logger) Router {
	return Router{
		officeHundler:  officeHundler,
		primaryHundler: primaryHundler,
		logger:         logger,
	}
}

func (r *Router) EntryPoint(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msgConfig tgbotapi.MessageConfig
	switch update.Message.Text {
	case "/start":
		msg, err := r.primaryHundler.Start(update)
		if err != nil {
			r.logger.Info("StartTelegramServer")
		}
		msgConfig = msg

	case hundlers.Yakutsk203, hundlers.YakutskGluhoi, hundlers.Moscow, hundlers.Almaty:
		//msg, err := r.officeHundler.SetOffice(update)
		//if err != nil {
		//	r.logger.Info("StartTelegramServer")
		//}
		//msgConfig = msg
	case "Выбрать место":

	case "Я пойду на обед":
		msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}
	//msg.ReplyToMessageID = update.Message.MessageID

	return msgConfig, nil
}
