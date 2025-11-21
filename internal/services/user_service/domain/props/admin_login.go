package props

import "errors"

type AdminLoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *AdminLoginReq) Validate() error {
	if r.Login == "" || r.Password == "" {
		return errors.New("login or password is empty")
	}
	return nil
}

type AdminLoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
