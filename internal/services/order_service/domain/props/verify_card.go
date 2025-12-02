package props

type VerifyCardReq struct {
	OTP      string `json:"otp"`
	CVC      string `json:"cvc"`
	UserUUID string `json:"-"`
}

type VerifyCardResp struct {
}
