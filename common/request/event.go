package request

type ApproveEventRequest struct {
	ConfirmedDate string `json:"confirmed_date"`
}

type RejectEventRequest struct {
	Remarks string `json:"remarks"`
}
