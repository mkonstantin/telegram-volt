package job

import (
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"time"
)

type dateJobsImpl struct {
	bookSeatRepo interfaces.BookSeatRepository
	seatRepo     interfaces.SeatRepository
	logger       *zap.Logger
}

type DateJob interface {
	CheckAndSetDates() error
}

func NewDateJob(bookSeatRepo interfaces.BookSeatRepository, seatRepo interfaces.SeatRepository, logger *zap.Logger) DateJob {
	return &dateJobsImpl{
		bookSeatRepo: bookSeatRepo,
		seatRepo:     seatRepo,
		logger:       logger,
	}
}

func (o *dateJobsImpl) CheckAndSetDates() error {
	weekDays := helper.WeekRange(year, week)

	for _, day := range weekDays {
		result, err := o.canSetSeats(officeID, day)
		if err != nil {
			o.logger.Error("SetNewSeatList", zap.Error(err))
			return err
		}
		if result {
			return o.insertSeatsTo(officeID, day)
		}
	}

	return nil
}

// В этом методе определяются условия при которых мы НЕ можем засетать места на сегодня.
// Условия: дни суббота, воскресенье, дни раньше текущего дня, день уже заполнен

func (o *dateJobsImpl) canSetSeats(officeID int64, bookDate time.Time) (bool, error) {

	bookedSeats, err := o.bookSeatRepo.GetAllByOfficeIDAndDate(officeID, bookDate.String())
	if err != nil {
		return false, err
	}

	today := helper.TodayZeroTimeUTC()

	// Условия не позволяющие засетать места:
	switch {
	case len(bookedSeats) > 0:
		fallthrough
	case bookDate.Before(today):
		fallthrough
	case bookDate.Weekday() == time.Saturday:
		fallthrough
	case bookDate.Weekday() == time.Sunday:
		return false, nil
	}
	return true, nil
}

func (o *dateJobsImpl) insertSeatsTo(officeID int64, date time.Time) error {
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
	o.logger.Info(fmt.Sprintf("Insert seats for office with ID %d, seats amount: %d, date: %s", officeID, len(seats), date.String()))
	return nil
}
