package model

import "time"

type BookSeat struct {
	ID            int64
	Office        Office
	Seat          Seat
	User          *User
	BookDate      time.Time
	BookStartTime *time.Time
	BookEndTime   *time.Time
	Confirm       bool
	Hold          bool
}
