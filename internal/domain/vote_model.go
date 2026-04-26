package domain

import (
	"time"

	"github.com/google/uuid"
)

type VoteModel struct {
	Id           uuid.UUID `db:"id"`
	DeviceId     string    `db:"device_id"`
	ShiftType    string    `db:"shift_type"`
	PhoneModel   string    `db:"phone_model"`
	Browser      string    `db:"browser"`
	ExternalIP   string    `db:"external_ip"`
	BusinessDate time.Time `db:"business_date"`
	CreatedAt    time.Time `db:"created_at"`
}

type VoteItemModel struct {
	Id        uuid.UUID `db:"id"`
	VoteId    uuid.UUID `db:"vote_id"`
	MealType  string    `db:"meal_type"`
	Rating    *int16    `db:"rating"`
	Review    *string   `db:"review"`
	CreatedAt time.Time `db:"created_at"`
}
