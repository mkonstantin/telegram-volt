package job

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
	"time"
)

const (
	plusDays = 10
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
	ctx := context.Background()

	startDate := helper.TodayZeroTimeUTC()
	endDate := helper.PlusDaysUTC(startDate, plusDays)

	dates, err := w.dateRepo.FindByDatesAndStatus(ctx, startDate.String(), endDate.String(), model.StatusWait)
	if err != nil {
		w.logger.Error("Scheduler Seat_jobs: w.dateRepo.FindByDatesAndStatus", zap.Error(err))
		return err
	}

	officeIDs, err := w.getOfficeIDs(ctx)
	if err != nil {
		w.logger.Error("Scheduler Seat_jobs: w.getOfficeIDs", zap.Error(err))
		return err
	}

	for _, day := range dates {
		err = w.fillByOffices(ctx, officeIDs, day)
		if err != nil {
			w.logger.Error("Scheduler Seat_jobs: w.fillByOffices", zap.Error(err))
			return err
		}
		err = w.dateRepo.UpdateStatusByID(ctx, day.ID, model.StatusSetBookSeats)
		if err != nil {
			w.logger.Error("Scheduler Seat_jobs: w.dateRepo.UpdateStatusByID", zap.Error(err))
			return err
		}
	}

	return nil
}

func (w *seatJobImpl) fillByOffices(ctx context.Context, officeIDs []int64, workDate model.WorkDate) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	for _, officeID := range officeIDs {
		err := w.insertSeatsTo(ctx, officeID, workDate.Date)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *seatJobImpl) getOfficeIDs(ctx context.Context) ([]int64, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	offices, err := w.officeRepo.GetAll(ctx)
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

func (w *seatJobImpl) insertSeatsTo(ctx context.Context, officeID int64, date time.Time) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	seats, err := w.seatRepo.GetAllByOfficeID(ctx, officeID)
	if err != nil {
		return err
	}

	for _, seat := range seats {
		err = w.bookSeatRepo.InsertSeat(ctx, officeID, seat.ID, date)
		if err != nil {
			w.logger.Error("InsertSeat", zap.Error(err))
			return err
		}
	}
	w.logger.Info(fmt.Sprintf("Insert seats for office with ID %d, seats amount: %d, date: %s", officeID, len(seats), date.String()))
	return nil
}
