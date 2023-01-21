package hundlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type CustomMessageHandler interface {
	Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type customMessageHandlerImpl struct {
	logger *zap.Logger
}

func NewCustomMessageHandler(logger *zap.Logger) CustomMessageHandler {
	return &inlineMessageHandlerImpl{
		logger: logger,
	}
}

func (s *customMessageHandlerImpl) Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {

	// And finally, send a message containing the data received.
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	return &msg, nil
}

//func (r *Router) EntryPoint(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
//	var msgConfig tgbotapi.MessageConfig
//	switch update.Message.Text {
//	case "/start":
//		msg, err := r.primaryHundler.Start(update)
//		if err != nil {
//			r.logger.Info("StartTelegramServer")
//		}
//		msgConfig = msg
//
//	case hundlers.Yakutsk203, hundlers.YakutskGluhoi, hundlers.Moscow, hundlers.Almaty:
//		msgConfig = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
//		msgConfig.ReplyToMessageID = update.Message.MessageID
//		//msg, err := r.officeHundler.SetOffice(update)
//		//if err != nil {
//		//	r.logger.Info("StartTelegramServer")
//		//}
//		//msgConfig = msg
//	case "Выбрать место":
//
//	case "Я пойду на обед":
//		msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
//	}
//	//msg.ReplyToMessageID = update.Message.MessageID
//
//	return msgConfig, nil
//}
