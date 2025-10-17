package props

type ConfirmCodeReq struct {
	Code string `json:"code"`
	Email string `json:"email"`
}

type ConfirmCodeResp struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}