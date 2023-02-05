package dto

import (
	"telegram-api/internal/domain/model"
	"time"
)

type Office struct {
	ID        int64     `db:"id,omitempty"`
	Name      string    `db:"name,omitempty"`
	City      string    `db:"city,omitempty"`
	TimeZone  string    `db:"time_zone,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (o *Office) ToModel() *model.Office {
	return &model.Office{
		ID:       o.ID,
		Name:     o.Name,
		City:     o.City,
		TimeZone: o.TimeZone,
	}
}

func ToOfficeModels(array []Office) []*model.Office {
	var models []*model.Office
	for _, item := range array {
		models = append(models, item.ToModel())
	}
	return models
}
