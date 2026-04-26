package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReportModel struct {
	VoteID       uuid.UUID `db:"vote_id"`
	ShiftType    string    `db:"shift_type"`
	MealType     *string   `db:"meal_type"`
	Rating       *int16    `db:"rating"`
	Review       *string   `db:"review"`
	BusinessDate time.Time `db:"business_date"`
}

type SentReportModel struct {
	PeriodStart time.Time `db:"period_start"`
	PeriodEnd   time.Time `db:"period_end"`
	SentAt      time.Time `db:"sent_at"`
}
