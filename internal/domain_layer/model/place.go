package model

type Place struct {
	ID   int64  `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}

type PlaceRepository interface {
	Create(place Place) (Place, error)
	Read(id int64) (Place, error)
	Update(place Place) (Place, error)
	Delete(id int64) error
}
