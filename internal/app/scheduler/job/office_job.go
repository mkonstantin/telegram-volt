package job

import (
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"time"
)

type officeJobsImpl struct {
	bookSeatRepo interfaces.BookSeatRepository
	seatRepo     interfaces.SeatRepository
	logger       *zap.Logger
}

type OfficeJob interface {
	SetSeatsWeek(officeID int64, year, week int) error
}

func NewOfficeJob(bookSeatRepo interfaces.BookSeatRepository, seatRepo interfaces.SeatRepository, logger *zap.Logger) OfficeJob {
	return &officeJobsImpl{
		bookSeatRepo: bookSeatRepo,
		seatRepo:     seatRepo,
		logger:       logger,
	}
}

func (o *officeJobsImpl) SetSeatsWeek(officeID int64, year, week int) error {
	weekDays := helper.WeekRange(year, week)

	for _, day := range weekDays {
		result, err := o.isSeatsExists(officeID, day)
		if err != nil {
			o.logger.Error("SetNewSeatList", zap.Error(err))
			return err
		}
		if !result {
			return o.insertSeatsTo(officeID, day)
		}
	}

	return nil
}

func (o *officeJobsImpl) isSeatsExists(officeID int64, bookDate time.Time) (bool, error) {

	bookedSeats, err := o.bookSeatRepo.GetAllByOfficeIDAndDate(officeID, bookDate.String())
	if err != nil {
		return false, err
	}

	if len(bookedSeats) > 0 {
		return true, nil
	}
	return false, nil
}

func (o *officeJobsImpl) insertSeatsTo(officeID int64, date time.Time) error {
	seats, err := o.seatRepo.GetAllByOfficeID(officeID)
	if err != nil {
		return err
	}

	for _, seat := range seats {
		err = o.bookSeatRepo.InsertSeat(officeID, seat.ID, date)
		if err != nil {
			o.logger.Error("InsertSeat", zap.Error(err))
			return err
		}
	}
	return nil
}
