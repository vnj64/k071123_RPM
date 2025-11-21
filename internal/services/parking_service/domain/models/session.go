package models

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	UUID        uuid.UUID     `json:"uuid"`
	ParkingUUID uuid.UUID     `json:"parking_uuid"`
	CarUUID     uuid.UUID     `json:"car_uuid"`
	Status      SessionStatus `json:"status"`
	StartAt     time.Time     `json:"start_at"`
	FinishAt    *time.Time    `json:"finish_at"`
	Cost        *float64      `json:"cost"`
}

type SessionStatus string

const (
	Active         SessionStatus = "active"
	WaitingPayment SessionStatus = "waiting_payment"
	Failed         SessionStatus = "failed"
	Finished       SessionStatus = "finished"
)
