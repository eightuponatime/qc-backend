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

		vote, err := s.voteRepo.GetVoteByDay(ctx, req.DeviceId, businessDate)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("get vote by day: %w", err)
			}

			voteToCreate := domain.VoteModel{
				DeviceId:     req.DeviceId,
				PhoneModel:   req.PhoneModel,
				Browser:      req.Browser,
				ExternalIP:   externalIp,
				BusinessDate: businessDate,
			}

			vote, err = s.voteRepo.CreateVote(ctx, voteToCreate)
			if err != nil {
				return fmt.Errorf("create vote: %w", err)
			}
		}

		for _, item := range req.Items {
			if !isValidMealType(item.MealType) {
				return fmt.Errorf("invalid meal type: %s", item.MealType)
			}

			voteItem := domain.VoteItemModel{
				VoteId:   vote.Id,
				MealType: item.MealType,
				Rating:   item.Rating,
				Review:   item.Review,
			}

			if err := s.voteRepo.UpsertVoteItem(ctx, voteItem); err != nil {
				return fmt.Errorf("upsert vote item (%s): %w", item.MealType, err)
			}
		}

		return nil
	})
}

func (s *VoteServiceImpl) GetTodayVote(ctx context.Context, deviceId string) (*dto.VoteResponseDto, error) {
	businessDate, err := s.getBusinessDate(time.Now())
	if err != nil {
		return nil, fmt.Errorf("get business date: %w", err)
	}

	vote, err := s.voteRepo.GetVoteByDay(ctx, deviceId, businessDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get vote by day: %w", err)
	}

	items, err := s.voteRepo.GetVoteItems(ctx, vote.Id)
	if err != nil {
		return nil, fmt.Errorf("get vote items: %w", err)
	}

	responseItems := make([]dto.VoteMealItemResponseDto, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, dto.VoteMealItemResponseDto{
			MealType: item.MealType,
			Rating:   item.Rating,
			Review:   item.Review,
		})
	}

	return &dto.VoteResponseDto{
		Items: responseItems,
	}, nil
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

func isValidMealType(mealType string) bool {
	switch mealType {
	case "breakfast", "lunch", "dinner":
		return true
	default:
		return false
	}
}