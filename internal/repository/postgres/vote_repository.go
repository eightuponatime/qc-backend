package postgres

import (
	"context"
	"fmt"
	"qc/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

type VoteRepository struct {
	db *sqlx.DB
}

func NewVoteRepository(db *sqlx.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) CreateVote(ctx context.Context, vote domain.VoteModel) error {
	db := extractTransaction(ctx, r.db)

	params := map[string]any{
		"device_id":     vote.DeviceId,
		"phone_model":   vote.PhoneModel,
		"browser":       vote.Browser,
		"breakfast":     vote.Breakfast,
		"lunch":         vote.Lunch,
		"dinner":        vote.Dinner,
		"external_ip":   vote.ExternalIP,
		"business_date": vote.BusinessDate.Format("2006-01-02"),
		"breakfast_at":  vote.BreakfastAt,
		"lunch_at":      vote.LunchAt,
		"dinner_at":     vote.DinnerAt,
	}

	_, err := sqlx.NamedExecContext(ctx, db, `
		INSERT INTO votes (
			device_id, phone_model, browser,
			breakfast, lunch, dinner,
			external_ip, business_date,
			breakfast_at, lunch_at, dinner_at
		) VALUES (
			:device_id, :phone_model, :browser,
			:breakfast, :lunch, :dinner,
			:external_ip, :business_date,
			:breakfast_at, :lunch_at, :dinner_at
		)
	`, params)
	if err != nil {
		return fmt.Errorf("create vote: %w", err)
	}

	return nil
}

func (r *VoteRepository) UpdateVote(
	ctx context.Context,
	vote domain.VoteUpdateModel,
	businessDate time.Time,
) error {
	db := extractTransaction(ctx, r.db)

	params := map[string]any{
		"device_id":     vote.DeviceId,
		"business_date": businessDate.Format("2006-01-02"),
		"breakfast":     vote.Breakfast,
		"lunch":         vote.Lunch,
		"dinner":        vote.Dinner,
		"breakfast_at":  vote.BreakfastAt,
		"lunch_at":      vote.LunchAt,
		"dinner_at":     vote.DinnerAt,
	}

	_, err := sqlx.NamedExecContext(ctx, db, `
		UPDATE votes SET
			breakfast    = COALESCE(:breakfast, breakfast),
			lunch        = COALESCE(:lunch, lunch),
			dinner       = COALESCE(:dinner, dinner),
			breakfast_at = COALESCE(:breakfast_at, breakfast_at),
			lunch_at     = COALESCE(:lunch_at, lunch_at),
			dinner_at    = COALESCE(:dinner_at, dinner_at)
		WHERE
			device_id = :device_id
			AND business_date = :business_date
	`, params)
	if err != nil {
		return fmt.Errorf("update vote: %w", err)
	}

	return nil
}

func (r *VoteRepository) GetTodayVote(
	ctx context.Context,
	deviceId string,
	businessDate time.Time,
) (*domain.VoteModel, error) {
	db := extractTransaction(ctx, r.db)

	var vote domain.VoteModel
	err := sqlx.GetContext(ctx, db, &vote, `
		SELECT
			id,
			device_id,
			phone_model,
			browser,
			breakfast,
			lunch,
			dinner,
			external_ip,
			business_date,
			breakfast_at,
			lunch_at,
			dinner_at,
			created_at
		FROM votes
		WHERE device_id = $1
		  AND business_date = $2
		LIMIT 1
	`, deviceId, businessDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("get today vote: %w", err)
	}

	return &vote, nil
}
