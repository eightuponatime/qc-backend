package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type AnalyticsAccessRepository struct {
	db *sqlx.DB
}

func NewAnalyticsAccessRepository(db *sqlx.DB) *AnalyticsAccessRepository {
	return &AnalyticsAccessRepository{db: db}
}

func (r *AnalyticsAccessRepository) Create(
	ctx context.Context,
	codeHash string,
	validFrom time.Time,
	validUntil time.Time,
) error {
	db := extractTransaction(ctx, r.db)

	_, err := db.ExecContext(ctx, `
		insert into analytics_access_codes (code_hash, valid_from, valid_until)
		values ($1, $2, $3)
	`, codeHash, validFrom.Format("2006-01-02"), validUntil.Format("2006-01-02"))
	if err != nil {
		return fmt.Errorf("create analytics access code: %w", err)
	}

	return nil
}

func (r *AnalyticsAccessRepository) ExistsValid(
	ctx context.Context,
	codeHash string,
	businessDate time.Time,
) (bool, error) {
	db := extractTransaction(ctx, r.db)

	var exists bool
	err := sqlx.GetContext(ctx, db, &exists, `
		select exists(
			select 1
			from analytics_access_codes
			where code_hash = $1
			  and valid_from <= $2
			  and valid_until >= $2
		)
	`, codeHash, businessDate.Format("2006-01-02"))
	if err != nil {
		return false, fmt.Errorf("check analytics access code: %w", err)
	}

	return exists, nil
}
