package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type InformerService interface {
	SeatComeFree(ctx context.Context, id int64) error
}

type informerServiceImpl struct {
	bookSeatRepo interfaces.BookSeatRepository
	logger       *zap.Logger
}

func NewInformer(bookSeatRepo interfaces.BookSeatRepository, logger *zap.Logger) InformerService {
	return &informerServiceImpl{
		bookSeatRepo: bookSeatRepo,
		logger:       logger,
	}
}

func (i *informerServiceImpl) SeatComeFree(ctx context.Context, id int64) error {
	//currentUser := model.GetCurrentUser(ctx)

	seat, err := i.bookSeatRepo.FindByID(id)
	if err != nil {
		return err
	}

	fmt.Println(seat.BookDate)

	return nil
}
