package model

type User struct {
	ID             int64
	Name           string
	TelegramID     int64
	TelegramName   string
	OfficeID       int64
	NotifyOfficeID int
}

func (u *User) HaveChosenOffice() bool {
	return u.OfficeID != 0
}

func (u *User) isNeedNotify() bool {
	return u.NotifyOfficeID > 0
}
