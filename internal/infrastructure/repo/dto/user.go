package dto

import (
	"telegram-api/internal/domain/model"
	"time"
)

type User struct {
	ID             int64     `db:"id,omitempty"`
	Name           string    `db:"name,omitempty"`
	TelegramID     int64     `db:"telegram_id,omitempty"`
	TelegramName   string    `db:"telegram_name,omitempty"`
	OfficeID       *int64    `db:"office_id,omitempty"`
	NotifyOfficeID int       `db:"notify_office_id,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

func (u *User) ToModel() *model.User {
	user := model.User{
		ID:             u.ID,
		Name:           u.Name,
		TelegramID:     u.TelegramID,
		TelegramName:   u.TelegramName,
		NotifyOfficeID: u.NotifyOfficeID,
	}
	if u.OfficeID == nil {
		user.OfficeID = 0
	} else {
		user.OfficeID = *u.OfficeID
	}
	return &user
}
