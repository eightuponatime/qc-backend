package impl

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"qc/config"
	"qc/internal/repository"
	"strings"
	"time"
)

type AnalyticsAccessService struct {
	rp  repository.AnalyticsAccessRepository
	cfg *config.Config
}

func NewAnalyticsAccessService(
	rp repository.AnalyticsAccessRepository,
	cfg *config.Config,
) *AnalyticsAccessService {
	return &AnalyticsAccessService{
		rp:  rp,
		cfg: cfg,
	}
}

func (s *AnalyticsAccessService) CreateAccessCode(
	ctx context.Context,
	validFrom time.Time,
	validUntil time.Time,
) (string, error) {
	code, err := generateAnalyticsCode()
	if err != nil {
		return "", err
	}

	if err := s.rp.Create(ctx, hashAnalyticsCode(code), validFrom, validUntil); err != nil {
		return "", err
	}

	return code, nil
}

func (s *AnalyticsAccessService) ValidateAccessCode(ctx context.Context, code string) (bool, error) {
	code = normalizeAnalyticsCode(code)
	if code == "" {
		return false, nil
	}

	location, err := time.LoadLocation(s.cfg.BusinessTimezone)
	if err != nil {
		return false, fmt.Errorf("load business timezone: %w", err)
	}

	businessDate := normalizeBusinessDate(time.Now(), location)
	return s.rp.ExistsValid(ctx, hashAnalyticsCode(code), businessDate)
}

func generateAnalyticsCode() (string, error) {
	buffer := make([]byte, 6)
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate analytics code: %w", err)
	}

	raw := strings.ToUpper(hex.EncodeToString(buffer))
	return raw[:4] + "-" + raw[4:8] + "-" + raw[8:12], nil
}

func normalizeAnalyticsCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func hashAnalyticsCode(code string) string {
	sum := sha256.Sum256([]byte(normalizeAnalyticsCode(code)))
	return hex.EncodeToString(sum[:])
}
