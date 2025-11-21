package props

type VerifyCardReq struct {
	OTP      string `json:"otp"`
	UserUUID string `json:"-"`
}

type VerifyCardResp struct {
}
