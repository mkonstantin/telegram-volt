package handler

import (
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/repo/interfaces"
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

	seats, err := o.seatRepo.GetAllByOfficeID(officeID)
	if err != nil {
		return err
	}

	for _, seat := range seats {
		currentTime := time.Now()
		date := currentTime.AddDate(0, 0, 1)
		bookDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, officeLocation)

		err = o.bookSeatRepo.InsertSeat(officeID, seat.ID, bookDate)
		if err != nil {
			return err
		}
	}
	return nil
}
