package router

import (
	"go.uber.org/zap"
)

type RouterData struct {
}

type Router interface {
}

type routerImpl struct {
	logger *zap.Logger
}

func NewRouter(logger *zap.Logger) Router {
	return &routerImpl{
		logger: logger,
	}
}
