package model

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
