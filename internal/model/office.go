package model

type Office struct {
	ID   int64  `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}
