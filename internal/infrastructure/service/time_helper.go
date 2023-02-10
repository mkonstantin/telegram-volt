package service

import (
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/common"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"time"
)

type TimeHelper interface {
	GetTodayZeroTimeWithOfficeID(officeID int64) (*time.Time, error)
	GetTomorrowZeroTimeWithOfficeID(officeID int64) (*time.Time, error)
}

type timeHelperImpl struct {
	officeRepo interfaces.OfficeRepository
	logger     *zap.Logger
}

func NewTimeHelper(officeRepo interfaces.OfficeRepository, logger *zap.Logger) TimeHelper {
	return &timeHelperImpl{
		officeRepo: officeRepo,
		logger:     logger,
	}
}

func (t *timeHelperImpl) GetTodayZeroTimeWithOfficeID(officeID int64) (*time.Time, error) {
	office, err := t.officeRepo.FindByID(officeID)
	if err != nil {
		t.logger.Error("GetTodayZeroTimeWithOfficeID officeRepo", zap.Error(err))
		return nil, err
	}
	if office == nil {
		return nil, common.ErrOfficeNotFound
	}

	location, err := time.LoadLocation(office.TimeZone)
	if err != nil {
		t.logger.Error("GetTodayZeroTimeWithOfficeID LoadLocation", zap.Error(err))
		return nil, err
	}

	date := time.Now()
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, location)
	return &bookDate, nil
}

func (t *timeHelperImpl) GetTomorrowZeroTimeWithOfficeID(officeID int64) (*time.Time, error) {
	office, err := t.officeRepo.FindByID(officeID)
	if err != nil {
		t.logger.Error("GetTodayZeroTimeWithOfficeID officeRepo", zap.Error(err))
		return nil, err
	}
	if office == nil {
		return nil, common.ErrOfficeNotFound
	}

	location, err := time.LoadLocation(office.TimeZone)
	if err != nil {
		t.logger.Error("GetTodayZeroTimeWithOfficeID LoadLocation", zap.Error(err))
		return nil, err
	}

	currentTime := time.Now()
	date := currentTime.AddDate(0, 0, 1)
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, location)
	return &bookDate, nil
}

func TodayZeroTime(location *time.Location) time.Time {
	date := time.Now()
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, location)
	return bookDate
}

func TomorrowZeroTime(location *time.Location) time.Time {
	currentTime := time.Now()
	date := currentTime.AddDate(0, 0, 1)
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, location)
	return bookDate
}

func TodayZeroTimeUTC() time.Time {
	date := time.Now()
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	return bookDate
}

func TomorrowZeroTimeUTC() time.Time {
	currentTime := time.Now()
	date := currentTime.AddDate(0, 0, 1)
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	return bookDate
}
