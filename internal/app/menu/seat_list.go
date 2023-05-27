package menu

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/domain/model"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
	"time"
)

type seatListMenuImpl struct {
	userService  usecase.UserService
	bookSeatRepo repo.BookSeatRepository
	seatListForm form.SeatListForm
	logger       *zap.Logger
}

const dateFormat = "02 January 2006"

func NewSeatListMenu(
	userService usecase.UserService,
	bookSeatRepo repo.BookSeatRepository,
	seatListForm form.SeatListForm,
	logger *zap.Logger) interfaces.SeatListMenu {

	return &seatListMenuImpl{
		userService:  userService,
		bookSeatRepo: bookSeatRepo,
		seatListForm: seatListForm,
		logger:       logger,
	}
}

func (s *seatListMenuImpl) Call(ctx context.Context, date time.Time, officeID int64) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	currentUser := model.GetCurrentUser(ctx)

	var callingOfficeID int64
	if officeID > 0 {
		callingOfficeID = officeID
	} else {
		callingOfficeID = currentUser.OfficeID
	}

	ctx, err := s.userService.SetOfficeScript(ctx, callingOfficeID)
	if err != nil {
		return nil, err
	}

	seats, err := s.bookSeatRepo.GetAllByOfficeIDAndDate(ctx, callingOfficeID, date.String())
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
