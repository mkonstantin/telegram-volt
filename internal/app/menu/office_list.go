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
	"telegram-api/pkg/tracing"
)

type officeListMenuImpl struct {
	officeRepo     repo.OfficeRepository
	officeListForm form.OfficeListForm
	logger         *zap.Logger
}

func NewOfficeListMenu(
	officeRepo repo.OfficeRepository,
	officeListForm form.OfficeListForm,
	logger *zap.Logger) interfaces.OfficeListMenu {

	return &officeListMenuImpl{
		officeRepo:     officeRepo,
		officeListForm: officeListForm,
		logger:         logger,
	}
}

func (o *officeListMenuImpl) Call(ctx context.Context) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	currentUser := model.GetCurrentUser(ctx)

	offices, err := o.officeRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("%s, давай выберем хотдеск:", currentUser.Name)
	data := form.OfficeListFormData{
		Offices: offices,
		Message: message,
	}

	return o.officeListForm.Build(ctx, data)
}
