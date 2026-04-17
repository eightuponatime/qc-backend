package repository

import (
	"context"
	"qc/internal/domain"
	"time"
)

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type VoteRepository interface {
	CreateVote(ctx context.Context, vote domain.VoteModel) error
	UpdateVote(ctx context.Context, vote domain.VoteUpdateModel, businessDate time.Time) error
	GetTodayVote(ctx context.Context, deviceId string, businessDate time.Time) (*domain.VoteModel, error)
}
