package postgres

import (
	"context"
	"fmt"
	"qc/internal/domain"

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
	_, err := sqlx.NamedExecContext(ctx, db, `
		INSERT INTO votes (
			device_id, phone_model, browser,
			breakfast, lunch, dinner,
			external_ip,
			breakfast_at, lunch_at, dinner_at
		) VALUES (
			:device_id, :phone_model, :browser,
			:breakfast, :lunch, :dinner,
			:external_ip,
			:breakfast_at, :lunch_at, :dinner_at
		)
	`, vote)
	if err != nil {
		return fmt.Errorf("create vote: %w", err)
	}
	return nil
}

func (r *VoteRepository) UpdateVote(ctx context.Context, vote domain.VoteUpdateModel) error {
	db := extractTransaction(ctx, r.db)
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
			AND DATE(created_at AT TIME ZONE 'Asia/Almaty') = CURRENT_DATE
	`, vote)
	if err != nil {
		return fmt.Errorf("update vote: %w", err)
	}
	return nil
}

func (r *VoteRepository) GetTodayVote(ctx context.Context, deviceId string) (*domain.VoteModel, error) {
	db := extractTransaction(ctx, r.db)
	var vote domain.VoteModel
	err := sqlx.GetContext(ctx, db, &vote, `
        SELECT * FROM votes
        WHERE device_id = $1
        AND DATE(created_at AT TIME ZONE 'Asia/Almaty') = CURRENT_DATE
    `, deviceId)
	if err != nil {
		return nil, fmt.Errorf("get today vote: %w", err)
	}
	return &vote, nil
}
