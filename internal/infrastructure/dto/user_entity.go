package dto

import (
	"telegram-api/internal/domain/model"
)

type UserEntity struct {
	ID          int64  `json:"id,omitempty" db:"id"`
	Name        string `json:"name,omitempty" db:"name"`
	OfficeName  string `json:"office_name,omitempty" db:"office_name"`
	PlaceName   string `json:"place_name,omitempty" db:"place_name"`
	IsGotoLunch bool   `json:"is_goto_lunch,omitempty" db:"is_goto_lunch"`
}

func (u *UserEntity) ToModel() model.User {
	user := model.User{}
	user.SetName(u.Name)
	user.ChooseOffice(u.OfficeName)
	user.ChoosePlace(u.PlaceName)
	if u.IsGotoLunch {
		user.WantLunch()
	} else {
		user.AbortLunch()
	}
	return user
}

func GetUserDTO(user *model.User) *UserEntity {
	if user == nil {
		return nil
	}
	uDTO := UserEntity{
		ID:          0,
		Name:        "",
		OfficeName:  "",
		PlaceName:   "",
		IsGotoLunch: false,
	}
	bDTO := BodyType{
		ID:          bodyType.ID,
		Name:        bodyType.NameLocalized,
		SimpleIcon:  bodyType.SimpleIcon,
		ColoredIcon: bodyType.ColoredIcon,
		Description: bodyType.DescriptionClientLocalized,
		Dimensions:  bodyType.DimensionClientLocalized,
	}
	if user != nil && user.IsDriver() {
		bDTO.Description = bodyType.DescriptionDriverLocalized
		bDTO.Dimensions = bodyType.DimensionDriverLocalized
	}
	return &bDTO
}
