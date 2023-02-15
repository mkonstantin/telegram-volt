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
)

type officeMenuImpl struct {
	officeRepo repo.OfficeRepository
	form       form.OfficeMenuForm
	logger     *zap.Logger
}

func NewOfficeMenu(
	officeRepo repo.OfficeRepository,
	form form.OfficeMenuForm,
	logger *zap.Logger) interfaces.OfficeMenu {

	return &officeMenuImpl{
		officeRepo: officeRepo,
		form:       form,
		logger:     logger,
	}
}

func (o *officeMenuImpl) Call(ctx context.Context) (*tgbotapi.MessageConfig, error) {
	currentUser := model.GetCurrentUser(ctx)

	office, err := o.officeRepo.FindByID(currentUser.OfficeID)
	if err != nil {
		return nil, err
	}

	var buttonText string

	if currentUser.OfficeID == currentUser.NotifyOfficeID {
		buttonText = "Отписаться"
	} else {
		buttonText = "Подписаться на свободные места"
	}

	message := fmt.Sprintf("Офис: %s, действия:", office.Name)

	result := &dto.UserResult{
		Key:                 "",
		Office:              office,
		Offices:             nil,
		Message:             message,
		SubscribeButtonText: buttonText,
	}
	return o.form.Build(ctx, result)
}
