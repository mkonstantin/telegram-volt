package handler

import "go.uber.org/zap"

type officeJobsImpl struct {
	logger *zap.Logger
}

type OfficeJob interface {
	BeginJob(officeID int64) error
}

func NewOfficeJob(logger *zap.Logger) OfficeJob {
	return &officeJobsImpl{
		logger: logger,
	}
}

func (o officeJobsImpl) BeginJob(officeID int64) error {
	//TODO implement me
	panic("implement me")
}
