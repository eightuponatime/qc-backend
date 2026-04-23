package impl

import (
	"context"
	"log/slog"
	"qc/config"
	"qc/internal/repository"
	"qc/internal/service"
	"time"
)

type ReportScheduler struct {
	dispatchService service.ReportDispatchService
	sentReportRepo  repository.SentReportRepository
	cfg             *config.Config
	interval        time.Duration
	now             func() time.Time
}

func NewReportScheduler(
	dispatchService service.ReportDispatchService,
	sentReportRepo repository.SentReportRepository,
	cfg *config.Config,
	interval time.Duration,
) *ReportScheduler {
	return &ReportScheduler{
		dispatchService: dispatchService,
		sentReportRepo:  sentReportRepo,
		cfg:             cfg,
		interval:        interval,
		now:             time.Now,
	}
}

func (s *ReportScheduler) Start(ctx context.Context) {
	s.runTick(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.runTick(ctx)
		}
	}
}

func (s *ReportScheduler) runTick(ctx context.Context) {
	location, err := time.LoadLocation(s.cfg.BusinessTimezone)
	if err != nil {
		slog.Error("report scheduler timezone load failed", slog.Any("error", err))
		return
	}

	shiftStart, err := time.ParseInLocation("2006-01-02", s.cfg.ShiftStartDate, location)
	if err != nil {
		slog.Error("report scheduler shift start parse failed", slog.Any("error", err))
		return
	}

	currentDate := normalizeBusinessDate(s.now(), location)
	currentPeriodStart := shiftStart

	if !currentDate.Before(shiftStart) {
		daysSinceStart := int(currentDate.Sub(shiftStart).Hours() / 24)
		periodIndex := daysSinceStart / 15
		currentPeriodStart = shiftStart.AddDate(0, 0, periodIndex*15)
	}

	if !currentDate.Equal(currentPeriodStart) {
		return
	}

	if !currentPeriodStart.After(shiftStart) {
		return
	}

	previousPeriodStart := currentPeriodStart.AddDate(0, 0, -15)
	previousPeriodEnd := currentPeriodStart.AddDate(0, 0, -1)

	alreadySent, err := s.sentReportRepo.ExistsByPeriod(ctx, previousPeriodStart, previousPeriodEnd)
	if err != nil {
		slog.Error(
			"report scheduler sent report check failed",
			slog.String("period_start", previousPeriodStart.Format("2006-01-02")),
			slog.String("period_end", previousPeriodEnd.Format("2006-01-02")),
			slog.Any("error", err),
		)
		return
	}
	if alreadySent {
		slog.Info(
			"report already sent for period",
			slog.String("period_start", previousPeriodStart.Format("2006-01-02")),
			slog.String("period_end", previousPeriodEnd.Format("2006-01-02")),
		)
		return
	}

	slog.Info(
		"report scheduler sending period report",
		slog.String("period_start", previousPeriodStart.Format("2006-01-02")),
		slog.String("period_end", previousPeriodEnd.Format("2006-01-02")),
	)

	if err := s.dispatchService.SendPeriodReport(ctx, previousPeriodStart, previousPeriodEnd); err != nil {
		slog.Error(
			"report scheduler send failed",
			slog.String("period_start", previousPeriodStart.Format("2006-01-02")),
			slog.String("period_end", previousPeriodEnd.Format("2006-01-02")),
			slog.Any("error", err),
		)
		return
	}

	if err := s.sentReportRepo.MarkAsSent(ctx, previousPeriodStart, previousPeriodEnd); err != nil {
		slog.Error(
			"report scheduler mark-as-sent failed",
			slog.String("period_start", previousPeriodStart.Format("2006-01-02")),
			slog.String("period_end", previousPeriodEnd.Format("2006-01-02")),
			slog.Any("error", err),
		)
		return
	}

	slog.Info(
		"period report sent successfully",
		slog.String("period_start", previousPeriodStart.Format("2006-01-02")),
		slog.String("period_end", previousPeriodEnd.Format("2006-01-02")),
	)
}
