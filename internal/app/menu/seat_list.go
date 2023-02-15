package menu

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/internal/infrastructure/service"
)

type seatListMenuImpl struct {
	bookSeatRepo repo.BookSeatRepository
	seatListForm form.SeatListForm
	logger       *zap.Logger
}

func NewSeatListMenu(
	bookSeatRepo repo.BookSeatRepository,
	seatListForm form.SeatListForm,
	logger *zap.Logger) interfaces.SeatListMenu {

	return &seatListMenuImpl{
		bookSeatRepo: bookSeatRepo,
		seatListForm: seatListForm,
		logger:       logger,
	}
}

func (s *seatListMenuImpl) Call(ctx context.Context) (*tgbotapi.MessageConfig, error) {

	currentUser := model.GetCurrentUser(ctx)

	date := service.TodayZeroTimeUTC()

	seats, err := s.bookSeatRepo.GetAllByOfficeIDAndDate(currentUser.OfficeID, date.String())
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Выберите место:")

	result := &dto.UserResult{
		Key:       "",
		Office:    nil,
		Offices:   nil,
		BookSeats: seats,
		Message:   message,
	}

	return s.seatListForm.Build(ctx, result)
}
