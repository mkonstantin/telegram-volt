package dto

import (
	"telegram-api/internal/domain/model"
	"time"
)

type WorkDate struct {
	ID        int64     `db:"id,omitempty"`
	Status    string    `db:"status,omitempty"`
	WorkDate  time.Time `db:"work_date,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (o *WorkDate) ToModel() *model.WorkDate {
	return &model.WorkDate{
		ID:       o.ID,
		Status:   o.Status,
		WorkDate: o.WorkDate,
	}
}

func ToWorkDateModels(array []WorkDate) []model.WorkDate {
	var models []model.WorkDate
	for _, item := range array {
		models = append(models, *item.ToModel())
	}
	return models
}
