package model

type User struct {
	name        string
	officeName  string
	placeName   string
	isGoToLunch bool
}

func (u *User) Name() string {
	return u.name
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) IsOfficeChoosed() bool {
	return u.officeName != ""
}

func (u *User) ChooseOffice(name string) {
	u.officeName = name
}

func (u *User) OfficeName() string {
	return u.officeName
}
