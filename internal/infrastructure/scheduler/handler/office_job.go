package handler

import (
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/internal/infrastructure/service"
	"time"
)

type officeJobsImpl struct {
	bookSeatRepo interfaces.BookSeatRepository
	seatRepo     interfaces.SeatRepository
	logger       *zap.Logger
}

type OfficeJob interface {
	SetNewSeatList(officeID int64, officeLocation *time.Location) error
}

func NewOfficeJob(bookSeatRepo interfaces.BookSeatRepository, seatRepo interfaces.SeatRepository, logger *zap.Logger) OfficeJob {
	return &officeJobsImpl{
		bookSeatRepo: bookSeatRepo,
		seatRepo:     seatRepo,
		logger:       logger,
	}
}

func (o *officeJobsImpl) SetNewSeatList(officeID int64, officeLocation *time.Location) error {

	bookDate := service.TomorrowZeroTimeUTC()
	result, err := o.isExistSeats(officeID, bookDate)
	if err != nil {
		o.logger.Error("SetNewSeatList", zap.Error(err))
		return err
	}
	if result {
		o.logger.Info("SetNewSeatList seats already set")
		return nil
	}

	seats, err := o.seatRepo.GetAllByOfficeID(officeID)
	if err != nil {
		return err
	}

	for _, seat := range seats {
		err = o.bookSeatRepo.InsertSeat(officeID, seat.ID, bookDate)
		if err != nil {
			o.logger.Error("InsertSeat", zap.Error(err))
			return err
		}
	}
	return nil
}

func (o *officeJobsImpl) isExistSeats(officeID int64, bookDate time.Time) (bool, error) {

	bookedSeats, err := o.bookSeatRepo.GetAllByOfficeID(officeID, bookDate.String())
	if err != nil {
		return false, err
	}

	if len(bookedSeats) > 0 {
		return true, nil
	}
	return false, nil
}
