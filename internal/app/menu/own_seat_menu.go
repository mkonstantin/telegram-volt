package menu

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/pkg/tracing"
)

type ownSeatMenuImpl struct {
	ownSeatForm form.OwnSeatForm
	logger      *zap.Logger
}

func NewOwnSeatMenu(
	ownSeatForm form.OwnSeatForm,
	logger *zap.Logger) interfaces.OwnSeatMenu {

	return &ownSeatMenuImpl{
		ownSeatForm: ownSeatForm,
		logger:      logger,
	}
}

func (f *ownSeatMenuImpl) Call(ctx context.Context, bookSeatID int64) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	message := "Вы уже заняли это место, хотите его освободить?"

	formData := form.OwnSeatFormData{
		Message:    message,
		BookSeatID: bookSeatID,
	}
	return f.ownSeatForm.Build(ctx, formData)
}
