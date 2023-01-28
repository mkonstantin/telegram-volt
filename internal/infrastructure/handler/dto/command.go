package dto

type CommandResponse struct {
	Type      string `json:"type,omitempty"`
	OfficeID  int64  `json:"office_id,omitempty"`
	IsConfirm bool   `json:"is_confirm,omitempty"`
}
