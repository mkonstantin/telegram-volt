package dto

type InlineRequest struct {
	Type       string `json:"type,omitempty"`
	OfficeID   int64  `json:"office_id,omitempty"`
	BookSeatID int64  `json:"book_seat_id,omitempty"`
	Action     int    `json:"action,omitempty"`
}

const (
	OfficeMenuFreeSeats           = 1
	OfficeMenuSubscribe           = 2
	OfficeMenuChooseAnotherOffice = 3

	ActionCancelBookYes = 11
	ActionCancelBookNo  = 12

	ActionBookYes = 21
	ActionBookNo  = 22
)
