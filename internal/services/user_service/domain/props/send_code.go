package props

type SendCodeReq struct {
	Email string `json:"email"`
}

type SendCodeResponse struct {
	Status string `json:"status"`
}
