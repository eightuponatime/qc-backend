package domain

import (
	"time"
	"github.com/google/uuid"
)

type VoteModel struct {
	Id       uuid.UUID `db:"id"`
	DeviceId string    `db:"device_id"`
	Breakfast *int16 `db:"breakfast"` 
	Lunch *int16 `db:"lunch"`
	Dinner *int16 `db:"dinner"`
	ExternalIP string `db:"external_ip"`
	VotedAt *time.Time `db:"voted_at"`
}
