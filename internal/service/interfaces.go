package service

import (
	"context"
	"qc/internal/dto"
	"time"
)

type VoteService interface {
	CreateVote(ctx context.Context, req dto.VoteRequestDto, externalIp string) error
	GetTodayVote(ctx context.Context, deviceId string) (*dto.VoteResponseDto, error)
}

type ReportService interface {
	CreateReport(ctx context.Context) (map[string]map[string][]dto.ReportVoteItemDto, error)
	CreateSummary(ctx context.Context) (*dto.ReportSummaryDto, error)
	CreateSummaryForPeriod(ctx context.Context, periodStart, periodEnd time.Time) (*dto.ReportSummaryDto, error)
	CreateAnalyticsSummary(ctx context.Context) (*dto.ReportAnalyticsSummaryDto, error)
	CreateAnalyticsSummaryForPeriod(ctx context.Context, periodStart, periodEnd time.Time) (*dto.ReportAnalyticsSummaryDto, error)
}

type ReportDispatchService interface {
	SendPeriodReport(ctx context.Context, periodStart, periodEnd time.Time) error
}

type AnalyticsAccessService interface {
	CreateAccessCode(ctx context.Context, validFrom, validUntil time.Time) (string, error)
	ValidateAccessCode(ctx context.Context, code string) (bool, error)
}
