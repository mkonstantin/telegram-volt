package model

type Office struct {
	ID   int64  `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}

type OfficeRepository interface {
	Create(office Office) (Office, error)
	Read(id int64) (Office, error)
	Update(office Office) (Office, error)
	Delete(id int64) error
}
