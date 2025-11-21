package props

type SendCodeReq struct {
	Email     string `json:"email"`
	CarNumber string `json:"car_number"`
}

type SendCodeResponse struct {
	Status string `json:"status"`
}
