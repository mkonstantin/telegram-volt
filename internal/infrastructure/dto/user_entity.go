package dto

type UserEntity struct {
	ID          int64  `json:"id,omitempty" db:"id"`
	Name        string `json:"name,omitempty" db:"name"`
	OfficeName  string `json:"office_name,omitempty" db:"office_name"`
	PlaceName   string `json:"place_name,omitempty" db:"place_name"`
	IsGotoLunch bool   `json:"is_goto_lunch,omitempty" db:"is_goto_lunch"`
}
