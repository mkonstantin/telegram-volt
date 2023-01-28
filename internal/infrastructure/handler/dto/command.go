package dto

type CommandResponse struct {
	Type     string `json:"type,omitempty"`
	OfficeID int64  `json:"office_id,omitempty"`
	Action   int    `json:"action,omitempty"`
}

const (
	OfficeMenuFreeSeats           = 1
	OfficeMenuSubscribe           = 2
	OfficeMenuChooseAnotherOffice = 3
)
