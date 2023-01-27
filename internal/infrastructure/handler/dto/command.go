package dto

type CommandResponse struct {
	Type          string         `json:"type"`
	ConfirmOffice *ConfirmOffice `json:"confirm,omitempty"`
	ChooseOffice  *ChooseOffice  `json:"choose,omitempty"`
}

type CommandType string

type ConfirmOffice struct {
	IsConfirm bool `json:"is_confirm"`
}

type ChooseOffice struct {
	OfficeID int64 `json:"office_id"`
}
