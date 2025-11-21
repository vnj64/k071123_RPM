package http

type Handlers struct {
	ParkingHandler *ParkingHandler
	TariffHandler  *TariffHandler
	SessionHandler *SessionHandler
	UnitHandler    *UnitHandler
}
