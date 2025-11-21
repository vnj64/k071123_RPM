package props

type FindPayments struct {
	SessionUUIDs []string `json:"session_uuids"`
	Statuses     []string `json:"statuses"`
}
