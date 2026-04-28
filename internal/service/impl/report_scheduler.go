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

	currentTime := s.now().In(location)
	periodStart, periodEnd, ok := resolveLatestReportPeriod(currentTime, location, s.cfg.ReportSendHour)
	if !ok {
		return
	}

	alreadySent, err := s.sentReportRepo.ExistsByPeriod(ctx, periodStart, periodEnd)
	if err != nil {
		slog.Error(
			"report scheduler sent report check failed",
			slog.String("period_start", periodStart.Format("2006-01-02")),
			slog.String("period_end", periodEnd.Format("2006-01-02")),
			slog.Any("error", err),
		)
		return
	}
	if alreadySent {
		return
	}

	slog.Info(
		"report scheduler sending period report",
		slog.String("period_start", periodStart.Format("2006-01-02")),
		slog.String("period_end", periodEnd.Format("2006-01-02")),
	)

	if err := s.dispatchService.SendPeriodReport(ctx, periodStart, periodEnd); err != nil {
		slog.Error(
			"report scheduler send failed",
			slog.String("period_start", periodStart.Format("2006-01-02")),
			slog.String("period_end", periodEnd.Format("2006-01-02")),
			slog.Any("error", err),
		)
		return
	}

	if err := s.sentReportRepo.MarkAsSent(ctx, periodStart, periodEnd); err != nil {
		slog.Error(
			"report scheduler mark-as-sent failed",
			slog.String("period_start", periodStart.Format("2006-01-02")),
			slog.String("period_end", periodEnd.Format("2006-01-02")),
			slog.Any("error", err),
		)
		return
	}

	slog.Info(
		"period report sent successfully",
		slog.String("period_start", periodStart.Format("2006-01-02")),
		slog.String("period_end", periodEnd.Format("2006-01-02")),
	)
}

func resolveLatestReportPeriod(currentTime time.Time, location *time.Location, reportSendHour int) (time.Time, time.Time, bool) {
	currentDate := normalizeBusinessDate(currentTime, location)

	periodYear := currentDate.Year()
	periodMonth := currentDate.Month()

	if currentDate.Day() < 16 || (currentDate.Day() == 16 && currentTime.Hour() < reportSendHour) {
		previousMonthDate := currentDate.AddDate(0, -1, 0)
		periodYear = previousMonthDate.Year()
		periodMonth = previousMonthDate.Month()
	}

	periodStart := time.Date(periodYear, periodMonth, 1, 0, 0, 0, 0, location)
	periodEnd := time.Date(periodYear, periodMonth, 15, 0, 0, 0, 0, location)

	return periodStart, periodEnd, true
}
