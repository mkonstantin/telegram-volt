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
	"telegram-api/internal/infrastructure/helper"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
)

type officeMenuImpl struct {
	userService    usecase.UserService
	dateRepo       repo.WorkDateRepository
	bookSeatRepo   repo.BookSeatRepository
	officeRepo     repo.OfficeRepository
	officeMenuForm form.OfficeMenuForm
	logger         *zap.Logger
}

func NewOfficeMenu(
	userService usecase.UserService,
	dateRepo repo.WorkDateRepository,
	bookSeatRepo repo.BookSeatRepository,
	officeRepo repo.OfficeRepository,
	officeMenuForm form.OfficeMenuForm,
	logger *zap.Logger) interfaces.OfficeMenu {

	return &officeMenuImpl{
		userService:    userService,
		dateRepo:       dateRepo,
		bookSeatRepo:   bookSeatRepo,
		officeRepo:     officeRepo,
		officeMenuForm: officeMenuForm,
		logger:         logger,
	}
}

func (o *officeMenuImpl) Call(ctx context.Context, title string, officeID int64) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	currentUser := model.GetCurrentUser(ctx)

	var callingOfficeID int64
	if officeID > 0 {
		callingOfficeID = officeID
	} else {
		callingOfficeID = currentUser.OfficeID
	}

	ctx, err := o.userService.SetOfficeScript(ctx, callingOfficeID)
	if err != nil {
		return nil, err
	}

	office, err := o.officeRepo.FindByID(ctx, callingOfficeID)
	if err != nil {
		return nil, err
	}

	dates, err := o.dateRepo.FindByStatus(ctx, model.StatusAccept)
	if err != nil {
		return nil, err
	}

	today := helper.TodayZeroTimeUTC()

	var needConfirmBookSeat *model.BookSeat
	var todayBook *model.BookSeat

	var bookSeats []*model.BookSeat
	for _, date := range dates {
		bookSeat, err := o.bookSeatRepo.FindByUserIDAndDate(ctx, currentUser.ID, date.Date.String())
		if err != nil {
			return nil, err
		}
		if bookSeat != nil {
			if bookSeat.Office.ID == callingOfficeID {
				bookSeats = append(bookSeats, bookSeat)
			}

			if bookSeat.BookDate == today && bookSeat.Office.ID == callingOfficeID {
				todayBook = bookSeat
			}

			// Проверяем нужно ли подтверждение места на сегодня
			if bookSeat.BookDate == today && !bookSeat.Confirm {
				currentTime, err := helper.CurrentTimeWithTimeZone(bookSeat.Office.TimeZone)
				morningTime, err := helper.TimeWithTimeZone(helper.Morning, bookSeat.Office.TimeZone)
				if err != nil {
					return nil, err
				}
				if currentTime.After(morningTime) || currentTime.Equal(morningTime) {
					if bookSeat.Office.ID == callingOfficeID {
						needConfirmBookSeat = bookSeat
					}
				}
			}
		}
	}

	var buttonText string
	if callingOfficeID == currentUser.NotifyOfficeID {
		buttonText = "🔕 Отписаться от уведомлений"
	} else {
		buttonText = "🔔 Подписаться на свободные места"
	}

	var message string
	if title != "" {
		message = fmt.Sprintf("%s\nОфис: %s, действия:", title, office.Name)
	} else {
		if todayBook != nil {
			if needConfirmBookSeat != nil {
				message = fmt.Sprintf("Офис: %s\nНа сегодня вы забронировали место №%s. Ждем подтверждения, что вы придете",
					office.Name, needConfirmBookSeat.Seat.SeatSign)
			} else {
				message = fmt.Sprintf("Офис: %s\nНа сегодня вы забронировали место №%s",
					office.Name, todayBook.Seat.SeatSign)
			}
		} else {
			if office.ID == 4 {
				message = fmt.Sprintf("Офис: %s\n‼️Бронирование в этом разделе возможно только для сотрудников департамента Marketing‼️", office.Name)
			} else {
				message = fmt.Sprintf("Офис: %s, действия:", office.Name)
			}
		}
	}

	data := form.OfficeMenuFormData{
		Office:              office,
		Message:             message,
		SubscribeButtonText: buttonText,
		BookSeats:           bookSeats,
		NeedConfirmBookSeat: needConfirmBookSeat,
	}
	return o.officeMenuForm.Build(ctx, data)
}
