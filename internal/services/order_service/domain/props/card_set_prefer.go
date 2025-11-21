package props

type SetCardAsPreferred struct {
	UserUUID string `json:"-"`
	CardUUID string `json:"card_uuid"`
}
