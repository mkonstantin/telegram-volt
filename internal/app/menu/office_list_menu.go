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
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type officeListMenuImpl struct {
	officeRepo interfaces.OfficeRepository
	form       form.OfficeListMenuForm
	logger     *zap.Logger
}

func NewOfficeListMenu(
	officeRepo interfaces.OfficeRepository,
	form form.OfficeListMenuForm,
	logger *zap.Logger) interfaces.OfficeListMenu {

	return &officeListMenuImpl{
		officeRepo: officeRepo,
		form:       form,
		logger:     logger,
	}
}

func (o *officeListMenuImpl) Call(ctx context.Context) (*tgbotapi.MessageConfig, error) {

	currentUser := model.GetCurrentUser(ctx)

	offices, err := o.officeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("%s, давай выберем офис:", currentUser.Name)
	result := &dto.UserResult{
		Key:     "",
		Office:  nil,
		Offices: offices,
		Message: message,
	}

	return o.form.Build(ctx, result)
}
