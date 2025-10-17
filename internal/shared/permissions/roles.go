package permissions

type Role string

const (
	SuperAdmin   Role = "superadmin"
	ParkingAdmin Role = "parkingadmin"
	Default      Role = "user"
)
