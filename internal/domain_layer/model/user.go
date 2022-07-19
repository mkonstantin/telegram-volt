package model

import "time"

type User struct {
	ID            int64
	Name          string
	officeID      int64
	placeID       int64
	wishGoToLunch bool
	LunchTime     time.Time
}

func (u *User) IsOfficeChoosed() bool {
	return u.officeID > 0
}

func (u *User) ChooseOffice(id int64) {
	u.officeID = id
}

func (u *User) WantLunch() {
	u.wishGoToLunch = true
	u.LunchTime = time.Now()
}

func (u *User) AbortLunch() {
	u.wishGoToLunch = false
}

func (u *User) ChoosePlace(id int64) {
	u.placeID = id
}
