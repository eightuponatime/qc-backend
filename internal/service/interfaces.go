package service

import (
	"context"
	"qc/internal/dto"
)

type VoteService interface {
	CreateVote(ctx context.Context, req dto.VoteRequestDto, externalIp string) error
	GetTodayVote(ctx context.Context, deviceId string) (*dto.VoteResponseDto, error)
}