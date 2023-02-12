package router

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/handler/dto"
)

type Data struct {
	Command string
	Data    dto.CommandResponse
}

type Router interface {
	Route(ctx context.Context, data Data) (*tgbotapi.MessageConfig, error)
}

type routerImpl struct {
	logger *zap.Logger
}

func NewRouter(logger *zap.Logger) Router {
	return &routerImpl{
		logger: logger,
	}
}

func (r routerImpl) Route(ctx context.Context, data Data) (*tgbotapi.MessageConfig, error) {
	//TODO implement me
	panic("implement me")
}
