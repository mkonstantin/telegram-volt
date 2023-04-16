package menu

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/config"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
)

type freeSeatMenuImpl struct {
	bookSeatRepo repo.BookSeatRepository
	freeSeatForm form.FreeSeatForm
	cfg          config.AppConfig
	logger       *zap.Logger
}

func NewFreeSeatMenu(
	bookSeatRepo repo.BookSeatRepository,
	freeSeatForm form.FreeSeatForm,
	cfg config.AppConfig,
	logger *zap.Logger) interfaces.FreeSeatMenu {

	return &freeSeatMenuImpl{
		bookSeatRepo: bookSeatRepo,
		freeSeatForm: freeSeatForm,
		cfg:          cfg,
		logger:       logger,
	}
}

func (f *freeSeatMenuImpl) Call(ctx context.Context, bookSeatID int64) (*tgbotapi.MessageConfig, error) {
	bookSeat, err := f.bookSeatRepo.FindByID(bookSeatID)
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("Занять место №%s?", bookSeat.Seat.SeatSign)

	formData := form.FreeSeatFormData{
		Message:    message,
		BookSeatID: bookSeatID,
	}
	currentUser := model.GetCurrentUser(ctx)
	return f.freeSeatForm.Build(ctx, formData, f.cfg.IsAdmin(currentUser.TelegramName))
}
