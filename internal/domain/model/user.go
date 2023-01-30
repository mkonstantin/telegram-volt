package model

import "context"

const ContextUserKey string = "ContextUserKey"

func GetCurrentUser(ctx context.Context) User {
	return ctx.Value(ContextUserKey).(User)
}

type User struct {
	ID           int64
	Name         string
	TelegramID   int64
	TelegramName string
	OfficeID     int64
}

func (u *User) HaveChosenOffice() bool {
	return u.OfficeID != 0
}
