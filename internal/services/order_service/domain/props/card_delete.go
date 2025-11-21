package props

type DeleteCard struct {
	UserUUID string `json:"-"`
	CardUUID string `json:"card_uuid"`
}
