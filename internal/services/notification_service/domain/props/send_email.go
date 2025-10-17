package props

type SendEmailRequest struct {
	Subject string   `json:"subject"`
	To      []string `json:"to"`
	Data    string   `json:"data"`
}

type SendEmailResp struct {
	Status string `json:"status"`
}
