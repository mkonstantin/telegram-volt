package job

import (
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"time"
)

const (
	daysToSeatsFill = 10
)

type seatJobImpl struct {
	officeRepo   interfaces.OfficeRepository
	dateRepo     interfaces.WorkDateRepository
	bookSeatRepo interfaces.BookSeatRepository
	seatRepo     interfaces.SeatRepository
	logger       *zap.Logger
}

type SeatJob interface {
	SetSeats() error
}

func NewSeatsJob(officeRepo interfaces.OfficeRepository,
	dateRepo interfaces.WorkDateRepository,
	bookSeatRepo interfaces.BookSeatRepository,
	seatRepo interfaces.SeatRepository,
	logger *zap.Logger) SeatJob {

	return &seatJobImpl{
		officeRepo:   officeRepo,
		dateRepo:     dateRepo,
		bookSeatRepo: bookSeatRepo,
		seatRepo:     seatRepo,
		logger:       logger,
	}
}

func (w *seatJobImpl) SetSeats() error {

	today := helper.TodayZeroTimeUTC()

	dates, err := w.dateRepo.FindByDateAndStatus(today.String(), model.StatusWait, daysToSeatsFill)
	if err != nil {
		return err
	}

	officeIDs, err := w.getOfficeIDs()
	if err != nil {
		return err
	}

	for _, day := range dates {
		err = w.fillByOffices(officeIDs, day)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *seatJobImpl) fillByOffices(officeIDs []int64, workDate model.WorkDate) error {
	for _, officeID := range officeIDs {
		err := w.insertSeatsTo(officeID, workDate.Date)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *seatJobImpl) getOfficeIDs() ([]int64, error) {
	w.logger.Info("Starting Enable Book scheduled job")

	offices, err := w.officeRepo.GetAll()
	if err != nil {
		w.logger.Error("Scheduler jobs get all offices error", zap.Error(err))
		return []int64{}, err
	}

	var ids []int64
	for _, office := range offices {
		ids = append(ids, office.ID)
	}
	return ids, nil
}

func (w *seatJobImpl) insertSeatsTo(officeID int64, date time.Time) error {
	seats, err := w.seatRepo.GetAllByOfficeID(officeID)
	if err != nil {
		return err
	}

	for _, seat := range seats {
		err = w.bookSeatRepo.InsertSeat(officeID, seat.ID, date)
		if err != nil {
			w.logger.Error("InsertSeat", zap.Error(err))
			return err
		}
	}
	w.logger.Info(fmt.Sprintf("Insert seats for office with ID %d, seats amount: %d, date: %s", officeID, len(seats), date.String()))
	return nil
}
