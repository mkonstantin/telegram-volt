package informer

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/form"
	"time"
)

type senderImpl struct {
	botAPI   *tgbotapi.BotAPI
	infoForm form.InfoMenuForm
	logger   *zap.Logger
}

type Sender interface {
	StartSenderJob()
	AddToQueue(array []form.InfoFormData)
}

func NewSender(botAPI *tgbotapi.BotAPI,
	infoForm form.InfoMenuForm,
	logger *zap.Logger) Sender {
	return &senderImpl{
		botAPI:   botAPI,
		infoForm: infoForm,
		logger:   logger,
	}
}

var messageQueue []form.InfoFormData

func (s *senderImpl) StartSenderJob() {

	s.logger.Info(fmt.Sprintf("SenderJob startedAt: %s", time.Now().String()))

	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(1).
		Second().
		Do(func() {
			s.logger.Info("gocron start SenderJob")

			err := s.getFromQueueAndSend()
			if err != nil {
				s.logger.Error("gocron execution SenderJob error", zap.Error(err))
			}
		})
	if err != nil {
		s.logger.Error("gocron create SenderJob error", zap.Error(err))
	}

	scheduler.StartAsync()
	s.logger.Info("Successfully started scheduled job: SenderJob")
}

func (s *senderImpl) AddToQueue(array []form.InfoFormData) {
	messageQueue = append(messageQueue, array...)
}

func (s *senderImpl) getFromQueueAndSend() error {
	var toSendSlice []form.InfoFormData

	if len(messageQueue) > 20 {
		toSendSlice = messageQueue[:20]
		messageQueue = messageQueue[20:]
	} else {
		toSendSlice = messageQueue
		messageQueue = nil
	}

	for _, data := range toSendSlice {
		s.send(data)
	}
	return nil
}

func (s *senderImpl) send(data form.InfoFormData) {
	build, err := s.infoForm.Build(context.Background(), data)

	if err != nil {
		s.logger.Error("Error Sender send", zap.Error(err)) // Ошибку не возвращаем, show must go on
	}

	if _, err = s.botAPI.Send(build); err != nil {
		s.logger.Error("Error Sender botAPI.Send(build)", zap.Error(err))
	}
}
