package dto

import "time"

type InlineRequest struct {
	Type       string     `json:"type,omitempty"`
	OfficeID   int64      `json:"office_id,omitempty"`
	BookSeatID int64      `json:"book_seat_id,omitempty"`
	Action     int        `json:"action,omitempty"`
	BookDate   *time.Time `json:"book_date,omitempty"`
}

const (
	OfficeMenuFreeSeats           = 1
	OfficeMenuSubscribe           = 2
	OfficeMenuChooseAnotherOffice = 3
	OfficeMenuCancelBook          = 4
	OfficeMenuConfirm             = 5

	ActionCancelBookYes = 11
	ActionCancelBookNo  = 12
	ActionCancelHold    = 14

	ActionBookYes  = 21
	ActionBookNo   = 22
	ActionBookHold = 23

	ActionShowOfficeMenu = 31
	ActionShowSeatList   = 32

	DateListShowMap = 41
	Back            = 100
)
