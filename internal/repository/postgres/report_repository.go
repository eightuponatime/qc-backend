package postgres

import (
	"context"
	"fmt"
	"qc/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ReportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetAllVotes(
	ctx context.Context,
) (*[]domain.ReportModel, error) {
	db := extractTransaction(ctx, r.db)

	var report []domain.ReportModel
	err := sqlx.SelectContext(ctx, db, &report, `
	select v.id as vote_id, vt.meal_type, vt.rating, 
		  vt.review, v.business_date
	from votes as v
	left join vote_items as vt on 
  		v.id = vt.vote_id
	order by v.business_date, v.id,
		case vt.meal_type
			when 'breakfast' then 1
			when 'lunch' then 2
			when 'dinner' then 3
			else 99
		end
	`)

	if err != nil {
		return nil, fmt.Errorf("get report info: %w", err)
	}

	return &report, nil
}
