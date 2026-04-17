// service/impl/vote_service.go
package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"qc/config"
	"qc/internal/domain"
	"qc/internal/dto"
	"qc/internal/repository"
	"time"
)

type VoteServiceImpl struct {
	voteRepo  repository.VoteRepository
	txManager repository.TransactionManager
	cfg       *config.Config
}

func NewVoteService(
	voteRepo repository.VoteRepository,
	txManager repository.TransactionManager,
	cfg *config.Config,
) *VoteServiceImpl {
	return &VoteServiceImpl{
		voteRepo:  voteRepo,
		txManager: txManager,
		cfg:       cfg,
	}
}

func (s *VoteServiceImpl) CreateVote(ctx context.Context, req dto.VoteRequestDto, externalIp string) error {
	return s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		now := time.Now()

		businessDate, err := s.getBusinessDate(now)
		if err != nil {
			return fmt.Errorf("get business date: %w", err)
		}

		existing, err := s.voteRepo.GetTodayVote(ctx, req.DeviceId, businessDate)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("check existing vote: %w", err)
		}

		if existing == nil {
			vote := domain.VoteModel{
				DeviceId:     req.DeviceId,
				PhoneModel:   req.PhoneModel,
				Browser:      req.Browser,
				ExternalIP:   externalIp,
				BusinessDate: businessDate,
				Breakfast:    req.Breakfast,
				Lunch:        req.Lunch,
				Dinner:       req.Dinner,
				BreakfastAt:  timestampIfNotNil(req.Breakfast, now),
				LunchAt:      timestampIfNotNil(req.Lunch, now),
				DinnerAt:     timestampIfNotNil(req.Dinner, now),
			}
			return s.voteRepo.CreateVote(ctx, vote)
		}

		update := domain.VoteUpdateModel{
			DeviceId:    req.DeviceId,
			Breakfast:   req.Breakfast,
			Lunch:       req.Lunch,
			Dinner:      req.Dinner,
			BreakfastAt: timestampIfNotNil(req.Breakfast, now),
			LunchAt:     timestampIfNotNil(req.Lunch, now),
			DinnerAt:    timestampIfNotNil(req.Dinner, now),
		}
		return s.voteRepo.UpdateVote(ctx, update, businessDate)
	})
}

func (s *VoteServiceImpl) GetTodayVote(ctx context.Context, deviceId string) (*dto.VoteResponseDto, error) {
	businessDate, err := s.getBusinessDate(time.Now())
	if err != nil {
		return nil, fmt.Errorf("get business date: %w", err)
	}

	vote, err := s.voteRepo.GetTodayVote(ctx, deviceId, businessDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get today vote: %w", err)
	}

	return &dto.VoteResponseDto{
		Breakfast: vote.Breakfast,
		Lunch:     vote.Lunch,
		Dinner:    vote.Dinner,
	}, nil
}

func timestampIfNotNil(val *int16, t time.Time) *time.Time {
	if val == nil {
		return nil
	}
	return &t
}

func (s *VoteServiceImpl) getBusinessDate(now time.Time) (time.Time, error) {
	location, err := time.LoadLocation(s.cfg.BusinessTimezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("load business timezone: %w", err)
	}

	localNow := now.In(location)

	businessDate := time.Date(
		localNow.Year(),
		localNow.Month(),
		localNow.Day(),
		0, 0, 0, 0,
		location,
	)

	return businessDate, nil
}
