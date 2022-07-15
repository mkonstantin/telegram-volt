package model

type User struct {
	ID   int64  `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}

type UserRepository interface {
	Create(user User) (User, error)
	Read(id int64) (User, error)
	Update(user User) (User, error)
	Delete(id int64) error
}
