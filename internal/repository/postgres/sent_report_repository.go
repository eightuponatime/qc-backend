package postgres

import (
	"context"
	"fmt"
	"qc/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

type SentReportRepository struct {
	db *sqlx.DB
}

func NewSentReportRepository(db *sqlx.DB) *SentReportRepository {
	return &SentReportRepository{db: db}
}

func (r *SentReportRepository) ExistsByPeriod(ctx context.Context, periodStart, periodEnd time.Time) (bool, error) {
	db := extractTransaction(ctx, r.db)

	var exists bool
	err := sqlx.GetContext(ctx, db, &exists, `
		select exists(
			select 1
			from sent_reports
			where period_start = $1
			  and period_end = $2
		)
	`, periodStart.Format("2006-01-02"), periodEnd.Format("2006-01-02"))
	if err != nil {
		return false, fmt.Errorf("check sent report by period: %w", err)
	}

	return exists, nil
}

func (r *SentReportRepository) MarkAsSent(ctx context.Context, periodStart, periodEnd time.Time) error {
	db := extractTransaction(ctx, r.db)

	_, err := db.ExecContext(ctx, `
		insert into sent_reports (period_start, period_end)
		values ($1, $2)
		on conflict (period_start, period_end) do nothing
	`, periodStart.Format("2006-01-02"), periodEnd.Format("2006-01-02"))
	if err != nil {
		return fmt.Errorf("mark report as sent: %w", err)
	}

	return nil
}

func (r *SentReportRepository) List(ctx context.Context) ([]domain.SentReportModel, error) {
	db := extractTransaction(ctx, r.db)

	var reports []domain.SentReportModel
	err := sqlx.SelectContext(ctx, db, &reports, `
		select period_start, period_end, sent_at
		from sent_reports
		order by period_start desc
	`)
	if err != nil {
		return nil, fmt.Errorf("list sent reports: %w", err)
	}

	return reports, nil
}
