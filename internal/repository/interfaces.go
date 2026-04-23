package repository

import (
	"context"
	"qc/internal/domain"
	"time"

	"github.com/google/uuid"
)

type VoteRepository interface {
	CreateVote(ctx context.Context, vote domain.VoteModel) (*domain.VoteModel, error)
	GetVoteByDay(ctx context.Context, deviceId string, businessDate time.Time) (*domain.VoteModel, error)

	UpsertVoteItem(ctx context.Context, item domain.VoteItemModel) error
	GetVoteItems(ctx context.Context, voteId uuid.UUID) ([]domain.VoteItemModel, error)
}

type ReportRepository interface {
	GetAllVotes(ctx context.Context) (*[]domain.ReportModel, error)
}

type SentReportRepository interface {
	ExistsByPeriod(ctx context.Context, periodStart, periodEnd time.Time) (bool, error)
	MarkAsSent(ctx context.Context, periodStart, periodEnd time.Time) error
	List(ctx context.Context) ([]domain.SentReportModel, error)
}

type AnalyticsAccessRepository interface {
	Create(ctx context.Context, codeHash string, validFrom, validUntil time.Time) error
	ExistsValid(ctx context.Context, codeHash string, businessDate time.Time) (bool, error)
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
