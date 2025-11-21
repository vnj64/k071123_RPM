package unit_statuses

type UnitStatus string

const (
	Active    UnitStatus = "active"
	Inactive  UnitStatus = "inactive"
	Warehouse UnitStatus = "warehouse"
	Repair    UnitStatus = "repair"
)

type NetworkUnitStatus string

const (
	Online  NetworkUnitStatus = "online"
	Offline NetworkUnitStatus = "offline"
)

type UnitDirection string

const (
	In  UnitDirection = "in"
	Out UnitDirection = "out"
)
