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

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
