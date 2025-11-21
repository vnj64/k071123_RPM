package parking_statuses

type ParkingStatus string

const (
	Active   ParkingStatus = "active"
	Inactive ParkingStatus = "inactive"
	Closed   ParkingStatus = "closed"
)
