package props

type SearchParkingReq struct {
	UUIDs []string `json:"uuids"`
	Query string   `json:"query"`
}
