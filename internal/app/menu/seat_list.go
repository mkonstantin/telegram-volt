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

const dateFormat = "02 January 2006"

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

func (s *seatListMenuImpl) Call(ctx context.Context, date time.Time, officeID int64) (*tgbotapi.MessageConfig, error) {

	currentUser := model.GetCurrentUser(ctx)

	var callingOfficeID int64
	if officeID > 0 {
		callingOfficeID = officeID
	} else {
		callingOfficeID = currentUser.OfficeID
	}

	seats, err := s.bookSeatRepo.GetAllByOfficeIDAndDate(callingOfficeID, date.String())
	if err != nil {
		return nil, err
	}

	formattedDate := date.Format(dateFormat)
	message := fmt.Sprintf("Выберите место на %s:", formattedDate)

	data := form.SeatListFormData{
		BookSeats: seats,
		Message:   message,
	}

	return s.seatListForm.Build(ctx, data)
}
