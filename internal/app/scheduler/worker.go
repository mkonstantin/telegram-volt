package scheduler

import (
	"go.uber.org/zap"
)

type workerImpl struct {
	logger *zap.Logger
}

type Worker interface {
	CleanTables() error
}

func NewWorker(logger *zap.Logger) Worker {
	return &workerImpl{
		logger: logger,
	}
}

func (w *workerImpl) CleanTables() error {

	return nil
}
