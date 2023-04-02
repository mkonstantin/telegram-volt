package common

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrBookSeatsNotFound        = errors.New("book seats not found")
	ErrDateSetBookSeatsNotFound = errors.New("date set book seats not found")
	ErrOfficeNotFound           = errors.New("office not found")
)
