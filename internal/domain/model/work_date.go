package model

import "time"

type DateStatus string

const (
	StatusWait         DateStatus = "wait"
	StatusSetBookSeats            = "set_book_seats"
	StatusAccept                  = "accept_to_book"
	StatusDone                    = "done"
)

type WorkDate struct {
	ID     int64
	Status string
	Date   time.Time
}
