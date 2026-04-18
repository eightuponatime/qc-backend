package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"qc/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type VoteRepository struct {
	db *sqlx.DB
}

func NewVoteRepository(db *sqlx.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) CreateVote(ctx context.Context, vote domain.VoteModel) (*domain.VoteModel, error) {
	db := extractTransaction(ctx, r.db)

	params := map[string]any{
		"device_id":     vote.DeviceId,
		"phone_model":   vote.PhoneModel,
		"browser":       vote.Browser,
		"external_ip":   vote.ExternalIP,
		"business_date": vote.BusinessDate.Format("2006-01-02"),
	}

	rows, err := sqlx.NamedQueryContext(ctx, db, `
		INSERT INTO votes (
			device_id,
			phone_model,
			browser,
			external_ip,
			business_date
		) VALUES (
			:device_id,
			:phone_model,
			:browser,
			:external_ip,
			:business_date
		)
		RETURNING id, device_id, phone_model, browser, external_ip, business_date, created_at
	`, params)
	if err != nil {
		return nil, fmt.Errorf("create vote: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var created domain.VoteModel
		if err := rows.StructScan(&created); err != nil {
			return nil, fmt.Errorf("scan created vote: %w", err)
		}
		return &created, nil
	}

	return nil, fmt.Errorf("create vote: no row returned")
}

func (r *VoteRepository) GetVoteByDay(
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
			external_ip,
			business_date,
			created_at
		FROM votes
		WHERE device_id = $1
		  AND business_date = $2
		LIMIT 1
	`, deviceId, businessDate.Format("2006-01-02"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("get vote by day: %w", err)
	}

	return &vote, nil
}

func (r *VoteRepository) UpsertVoteItem(ctx context.Context, item domain.VoteItemModel) error {
	db := extractTransaction(ctx, r.db)

	params := map[string]any{
		"vote_id":    item.VoteId,
		"meal_type":  item.MealType,
		"rating":     item.Rating,
		"review":     item.Review,
	}

	_, err := sqlx.NamedExecContext(ctx, db, `
		INSERT INTO vote_items (
			vote_id,
			meal_type,
			rating,
			review
		) VALUES (
			:vote_id,
			:meal_type,
			:rating,
			:review
		)
		ON CONFLICT (vote_id, meal_type)
		DO UPDATE SET
			rating = EXCLUDED.rating,
			review = EXCLUDED.review
	`, params)
	if err != nil {
		return fmt.Errorf("upsert vote item: %w", err)
	}

	return nil
}

func (r *VoteRepository) GetVoteItems(ctx context.Context, voteId uuid.UUID) ([]domain.VoteItemModel, error) {
	db := extractTransaction(ctx, r.db)

	var items []domain.VoteItemModel
	err := sqlx.SelectContext(ctx, db, &items, `
		SELECT
			id,
			vote_id,
			meal_type,
			rating,
			review,
			created_at
		FROM vote_items
		WHERE vote_id = $1
		ORDER BY
			CASE meal_type
				WHEN 'breakfast' THEN 1
				WHEN 'lunch' THEN 2
				WHEN 'dinner' THEN 3
				ELSE 99
			END
	`, voteId)
	if err != nil {
		return nil, fmt.Errorf("get vote items: %w", err)
	}

	return items, nil
}