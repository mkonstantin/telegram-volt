package menu

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
	"time"
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

func (s *seatListMenuImpl) Call(ctx context.Context, date time.Time) (*tgbotapi.MessageConfig, error) {

	currentUser := model.GetCurrentUser(ctx)

	seats, err := s.bookSeatRepo.GetAllByOfficeIDAndDate(currentUser.OfficeID, date.String())
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Выберите место:")

	data := form.SeatListFormData{
		BookSeats: seats,
		Message:   message,
	}

	return s.seatListForm.Build(ctx, data)
}
