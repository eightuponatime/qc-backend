package repository

import (
	"context"
	"qc/internal/domain"
)

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type VoteRepository interface {
	CreateVote(ctx context.Context, voteRequest domain.VoteModel) error
	UpdateVote(ctx context.Context, voteUpdate domain.VoteUpdateModel) error
	GetTodayVote(ctx context.Context, deviceId string) (*domain.VoteModel, error)
}
