package model

import "time"

const (
	StatusWait         = "wait"
	StatusSetBookSeats = "set_book_seats"
	StatusAccept       = "accept_to_book"
	StatusDone         = "done"
)

type WorkDate struct {
	ID       int64
	Status   string
	WorkDate time.Time
}
