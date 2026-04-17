package domain

import (
	"time"

	"github.com/google/uuid"
)

type VoteModel struct {
	Id           uuid.UUID  `db:"id"`
	DeviceId     string     `db:"device_id"`
	PhoneModel   string     `db:"phone_model"`
	Browser      string     `db:"browser"`
	Breakfast    *int16     `db:"breakfast"`
	Lunch        *int16     `db:"lunch"`
	Dinner       *int16     `db:"dinner"`
	ExternalIP   string     `db:"external_ip"`
	BusinessDate time.Time  `db:"business_date"`
	BreakfastAt  *time.Time `db:"breakfast_at"`
	LunchAt      *time.Time `db:"lunch_at"`
	DinnerAt     *time.Time `db:"dinner_at"`
	CreatedAt    time.Time  `db:"created_at"`
}

type VoteUpdateModel struct {
	DeviceId    string     `db:"device_id"`
	Breakfast   *int16     `db:"breakfast"`
	Lunch       *int16     `db:"lunch"`
	Dinner      *int16     `db:"dinner"`
	BreakfastAt *time.Time `db:"breakfast_at"`
	LunchAt     *time.Time `db:"lunch_at"`
	DinnerAt    *time.Time `db:"dinner_at"`
}
