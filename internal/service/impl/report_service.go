package impl

import (
	"context"
	"fmt"
	"qc/config"
	"qc/internal/dto"
	"qc/internal/repository"
	"sort"
	"strings"
	"time"
)

type ReportService struct {
	rp  repository.ReportRepository
	cfg *config.Config
}

func NewReportService(rp repository.ReportRepository, cfg *config.Config) *ReportService {
	return &ReportService{
		rp:  rp,
		cfg: cfg,
	}
}

func (r *ReportService) CreateReport(
	ctx context.Context,
) (map[string]map[string][]dto.ReportVoteItemDto, error) {
	reportModels, err := r.rp.GetAllVotes(ctx)
	if err != nil {
		return nil, err
	}

	location, err := time.LoadLocation(r.cfg.BusinessTimezone)
	if err != nil {
		return nil, fmt.Errorf("load business timezone: %w", err)
	}

	dateMap := make(map[string]map[string][]dto.ReportVoteItemDto)

	for _, model := range *reportModels {
		businessDate := model.BusinessDate.In(location).Format("2006-01-02")

		if _, ok := dateMap[businessDate]; !ok {
			dateMap[businessDate] = make(map[string][]dto.ReportVoteItemDto)
		}

		voteID := model.VoteID.String()
		if _, ok := dateMap[businessDate][voteID]; !ok {
			dateMap[businessDate][voteID] = []dto.ReportVoteItemDto{}
		}

		if model.MealType == nil {
			continue
		}

		review := ""
		if model.Review != nil {
			review = *model.Review
		}

		dateMap[businessDate][voteID] = append(dateMap[businessDate][voteID], dto.ReportVoteItemDto{
			MealType:     *model.MealType,
			Rating:       model.Rating,
			Review:       review,
			BusinessDate: businessDate,
		})
	}

	return dateMap, nil
}

func (r *ReportService) CreateSummary(ctx context.Context) (*dto.ReportSummaryDto, error) {
	_, periodStart, periodEnd, err := r.getCurrentPeriodBounds()
	if err != nil {
		return nil, err
	}

	return r.CreateSummaryForPeriod(ctx, periodStart, periodEnd)
}

func (r *ReportService) CreateSummaryForPeriod(
	ctx context.Context,
	periodStart time.Time,
	periodEnd time.Time,
) (*dto.ReportSummaryDto, error) {
	reportModels, err := r.rp.GetAllVotes(ctx)
	if err != nil {
		return nil, err
	}

	location, err := time.LoadLocation(r.cfg.BusinessTimezone)
	if err != nil {
		return nil, fmt.Errorf("load business timezone: %w", err)
	}

	periodStart = normalizeBusinessDate(periodStart, location)
	periodEnd = normalizeBusinessDate(periodEnd, location)

	summary := &dto.ReportSummaryDto{
		PeriodStart:        periodStart.Format("2006-01-02"),
		PeriodEnd:          periodEnd.Format("2006-01-02"),
		PeriodStartDisplay: formatRussianDate(periodStart),
		PeriodEndDisplay:   formatRussianDate(periodEnd),
		PeriodShortDisplay: formatShortDateRange(periodStart, periodEnd),
		WeekdayStats: []dto.ReportWeekdayStatsDto{
			newWeekdayStats("Monday"),
			newWeekdayStats("Tuesday"),
			newWeekdayStats("Wednesday"),
			newWeekdayStats("Thursday"),
			newWeekdayStats("Friday"),
			newWeekdayStats("Saturday"),
			newWeekdayStats("Sunday"),
		},
		MealStats: []dto.ReportMealStatsDto{
			{MealType: "breakfast"},
			{MealType: "lunch"},
			{MealType: "dinner"},
		},
		Insights: []string{},
	}

	uniqueVotes := make(map[string]struct{})
	weekdayIndexByName := map[string]int{
		"Monday":    0,
		"Tuesday":   1,
		"Wednesday": 2,
		"Thursday":  3,
		"Friday":    4,
		"Saturday":  5,
		"Sunday":    6,
	}
	mealIndexByType := map[string]int{
		"breakfast": 0,
		"lunch":     1,
		"dinner":    2,
	}

	weekdayLowCounts := make(map[string]int)
	weekdayHighCounts := make(map[string]int)
	weekdayTextReviews := make(map[string]int)
	mealLowCounts := make(map[string]int)

	for _, model := range *reportModels {
		businessDate := model.BusinessDate.In(location)
		businessDate = time.Date(
			businessDate.Year(),
			businessDate.Month(),
			businessDate.Day(),
			0, 0, 0, 0,
			location,
		)

		if businessDate.Before(periodStart) || businessDate.After(periodEnd) {
			continue
		}

		uniqueVotes[model.VoteID.String()] = struct{}{}

		if model.MealType == nil || model.Rating == nil {
			continue
		}

		summary.TotalRatings++
		incrementDistribution(&summary.RatingDistribution, *model.Rating)

		review := ""
		if model.Review != nil {
			review = strings.TrimSpace(*model.Review)
		}
		if review != "" {
			summary.TextReviewsCount++
			weekdayTextReviews[businessDate.Weekday().String()]++
		}

		weekdayName := businessDate.Weekday().String()
		if idx, ok := weekdayIndexByName[weekdayName]; ok {
			summary.WeekdayStats[idx].TotalRatings++
			summary.WeekdayStats[idx].TextReviewsCount += boolToInt(review != "")
			incrementDistribution(&summary.WeekdayStats[idx].RatingDistribution, *model.Rating)

			if mealIdx, ok := mealIndexByType[*model.MealType]; ok {
				summary.WeekdayStats[idx].MealStats[mealIdx].TotalRatings++
				if *model.Rating <= 3 {
					summary.WeekdayStats[idx].MealStats[mealIdx].LowRatingsCount++
				}
			}
		}

		if idx, ok := mealIndexByType[*model.MealType]; ok {
			summary.MealStats[idx].TotalRatings++
			if *model.Rating <= 3 {
				summary.MealStats[idx].LowRatingsCount++
				mealLowCounts[*model.MealType]++
			}
		}

		if *model.Rating <= 3 {
			weekdayLowCounts[weekdayName]++
		}
		if *model.Rating >= 4 {
			weekdayHighCounts[weekdayName]++
		}
	}

	summary.TotalVotes = len(uniqueVotes)
	summary.Insights = buildInsights(
		weekdayLowCounts,
		weekdayHighCounts,
		weekdayTextReviews,
		mealLowCounts,
	)

	return summary, nil
}

func (r *ReportService) CreateAnalyticsSummary(ctx context.Context) (*dto.ReportAnalyticsSummaryDto, error) {
	_, periodStart, periodEnd, err := r.getCurrentPeriodBounds()
	if err != nil {
		return nil, err
	}

	return r.CreateAnalyticsSummaryForPeriod(ctx, periodStart, periodEnd)
}

func (r *ReportService) CreateAnalyticsSummaryForPeriod(
	ctx context.Context,
	periodStart time.Time,
	periodEnd time.Time,
) (*dto.ReportAnalyticsSummaryDto, error) {
	reportModels, err := r.rp.GetAllVotes(ctx)
	if err != nil {
		return nil, err
	}

	location, err := time.LoadLocation(r.cfg.BusinessTimezone)
	if err != nil {
		return nil, fmt.Errorf("load business timezone: %w", err)
	}

	periodStart = normalizeBusinessDate(periodStart, location)
	periodEnd = normalizeBusinessDate(periodEnd, location)

	summary, err := r.CreateSummaryForPeriod(ctx, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}

	analyticsSummary := &dto.ReportAnalyticsSummaryDto{
		PeriodStart:            periodStart.Format("2006-01-02"),
		PeriodEnd:              periodEnd.Format("2006-01-02"),
		GeneratedAt:            time.Now().In(location).Format("2006-01-02 15:04:05"),
		Summary:                *summary,
		CalendarDateStats:      []dto.ReportCalendarDateStatsDto{},
		AttentionRequiredItems: []dto.ReportReviewDto{},
		DetailedReviewsByDate:  []dto.ReportDateReviewsDto{},
	}

	calendarStatsByDate := make(map[string]*dto.ReportCalendarDateStatsDto)
	detailedReviewsByDate := make(map[string][]dto.ReportReviewDto)

	for _, model := range *reportModels {
		businessDate := normalizeBusinessDate(model.BusinessDate, location)
		if businessDate.Before(periodStart) || businessDate.After(periodEnd) {
			continue
		}

		businessDateString := businessDate.Format("2006-01-02")
		if _, ok := calendarStatsByDate[businessDateString]; !ok {
			calendarStatsByDate[businessDateString] = &dto.ReportCalendarDateStatsDto{
				BusinessDate: businessDateString,
			}
		}

		if model.MealType == nil || model.Rating == nil {
			continue
		}

		calendarStatsByDate[businessDateString].TotalRatings++
		incrementDistribution(&calendarStatsByDate[businessDateString].RatingDistribution, *model.Rating)

		review := ""
		if model.Review != nil {
			review = strings.TrimSpace(*model.Review)
		}

		reviewItem := dto.ReportReviewDto{
			BusinessDate:        businessDateString,
			BusinessDateDisplay: formatRussianDate(businessDate),
			BusinessWeekday:     weekdayNominativeRussian(businessDate.Weekday().String()),
			VoteID:              model.VoteID.String(),
			MealType:            *model.MealType,
			Rating:              *model.Rating,
			Review:              review,
		}

		if *model.Rating <= 3 {
			analyticsSummary.AttentionRequiredItems = append(analyticsSummary.AttentionRequiredItems, reviewItem)
		}

		detailedReviewsByDate[businessDateString] = append(detailedReviewsByDate[businessDateString], reviewItem)
	}

	analyticsSummary.CalendarDateStats = buildCalendarDateStats(periodStart, periodEnd, calendarStatsByDate)
	analyticsSummary.AttentionRequiredItems = sortAttentionReviews(analyticsSummary.AttentionRequiredItems)
	analyticsSummary.DetailedReviewsByDate = buildDetailedReviewsByDate(periodStart, periodEnd, detailedReviewsByDate)

	return analyticsSummary, nil
}

// ===============
// HELPER METHODS
// ===============

func newWeekdayStats(name string) dto.ReportWeekdayStatsDto {
	return dto.ReportWeekdayStatsDto{
		Weekday: name,
		MealStats: []dto.ReportMealStatsDto{
			{MealType: "breakfast"},
			{MealType: "lunch"},
			{MealType: "dinner"},
		},
	}
}

func incrementDistribution(distribution *dto.ReportRatingDistributionDto, rating int16) {
	switch rating {
	case 5:
		distribution.Five++
	case 4:
		distribution.Four++
	case 3:
		distribution.Three++
	case 2:
		distribution.Two++
	case 1:
		distribution.One++
	}
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func buildInsights(
	weekdayLowCounts map[string]int,
	weekdayHighCounts map[string]int,
	weekdayTextReviews map[string]int,
	mealLowCounts map[string]int,
) []string {
	insights := make([]string, 0, 4)

	if weekday, count := maxMapKey(weekdayLowCounts); count > 0 {
		insights = append(
			insights,
			fmt.Sprintf(
				"В %s зафиксировано наибольшее количество низких оценок (1-3): %d.",
				weekdayToRussian(weekday),
				count,
			),
		)
	}

	if weekday, count := maxMapKey(weekdayTextReviews); count > 0 {
		insights = append(
			insights,
			fmt.Sprintf(
				"В %s отмечено наибольшее количество текстовых отзывов: %d.",
				weekdayToRussian(weekday),
				count,
			),
		)
	}

	if weekday, count := maxMapKey(weekdayHighCounts); count > 0 {
		insights = append(
			insights,
			fmt.Sprintf(
				"В %s преобладали высокие оценки (4-5): %d.",
				weekdayToRussian(weekday),
				count,
			),
		)
	}

	if meal, count := maxMapKey(mealLowCounts); count > 0 {
		insights = append(
			insights,
			fmt.Sprintf(
				"По приему пищи %q отмечено наибольшее количество низких оценок: %d.",
				mealToRussian(meal),
				count,
			),
		)
	}

	return insights
}

func maxMapKey(values map[string]int) (string, int) {
	maxKey := ""
	maxValue := 0

	for key, value := range values {
		if value > maxValue {
			maxKey = key
			maxValue = value
		}
	}

	return maxKey, maxValue
}

func weekdayToRussian(weekday string) string {
	switch weekday {
	case "Monday":
		return "понедельник"
	case "Tuesday":
		return "вторник"
	case "Wednesday":
		return "среду"
	case "Thursday":
		return "четверг"
	case "Friday":
		return "пятницу"
	case "Saturday":
		return "субботу"
	case "Sunday":
		return "воскресенье"
	default:
		return weekday
	}
}

func mealToRussian(meal string) string {
	switch meal {
	case "breakfast":
		return "завтрак"
	case "lunch":
		return "обед"
	case "dinner":
		return "ужин"
	default:
		return meal
	}
}

func (r *ReportService) getCurrentPeriodBounds() (*time.Location, time.Time, time.Time, error) {
	location, err := time.LoadLocation(r.cfg.BusinessTimezone)
	if err != nil {
		return nil, time.Time{}, time.Time{}, fmt.Errorf("load business timezone: %w", err)
	}

	shiftStart, err := time.ParseInLocation("2006-01-02", r.cfg.ShiftStartDate, location)
	if err != nil {
		return nil, time.Time{}, time.Time{}, fmt.Errorf("parse shift start date: %w", err)
	}

	nowBusinessDate := normalizeBusinessDate(time.Now(), location)

	periodStart := shiftStart
	if !nowBusinessDate.Before(shiftStart) {
		daysSinceStart := int(nowBusinessDate.Sub(shiftStart).Hours() / 24)
		periodIndex := daysSinceStart / 15
		periodStart = shiftStart.AddDate(0, 0, periodIndex*15)
	}

	periodEnd := periodStart.AddDate(0, 0, 14)
	return location, periodStart, periodEnd, nil
}

func normalizeBusinessDate(input time.Time, location *time.Location) time.Time {
	local := input.In(location)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, location)
}

func buildCalendarDateStats(
	periodStart time.Time,
	periodEnd time.Time,
	statsByDate map[string]*dto.ReportCalendarDateStatsDto,
) []dto.ReportCalendarDateStatsDto {
	result := make([]dto.ReportCalendarDateStatsDto, 0, 15)

	for current := periodStart; !current.After(periodEnd); current = current.AddDate(0, 0, 1) {
		dateString := current.Format("2006-01-02")
		if stat, ok := statsByDate[dateString]; ok {
			result = append(result, *stat)
			continue
		}

		result = append(result, dto.ReportCalendarDateStatsDto{
			BusinessDate: dateString,
		})
	}

	return result
}

func sortAttentionReviews(items []dto.ReportReviewDto) []dto.ReportReviewDto {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Rating != items[j].Rating {
			return items[i].Rating < items[j].Rating
		}
		if items[i].BusinessDate != items[j].BusinessDate {
			return items[i].BusinessDate < items[j].BusinessDate
		}
		if hasTextReview(items[i]) != hasTextReview(items[j]) {
			return hasTextReview(items[i])
		}
		return items[i].MealType < items[j].MealType
	})

	return items
}

func buildDetailedReviewsByDate(
	periodStart time.Time,
	periodEnd time.Time,
	reviewsByDate map[string][]dto.ReportReviewDto,
) []dto.ReportDateReviewsDto {
	result := make([]dto.ReportDateReviewsDto, 0, 15)

	for current := periodStart; !current.After(periodEnd); current = current.AddDate(0, 0, 1) {
		dateString := current.Format("2006-01-02")
		reviews := reviewsByDate[dateString]
		textReviews := make([]dto.ReportReviewDto, 0, len(reviews))
		for _, review := range reviews {
			if hasTextReview(review) {
				textReviews = append(textReviews, review)
			}
		}

		sort.SliceStable(textReviews, func(i, j int) bool {
			if textReviews[i].Rating != textReviews[j].Rating {
				return textReviews[i].Rating < textReviews[j].Rating
			}
			return textReviews[i].MealType < textReviews[j].MealType
		})

		result = append(result, dto.ReportDateReviewsDto{
			BusinessDate:         dateString,
			BusinessDateDisplay:  formatRussianDate(current),
			TotalReviews:         len(textReviews),
			PositiveReviewsCount: countPositiveReviews(textReviews),
			LowReviewsCount:      countLowReviews(textReviews),
			Reviews:              textReviews,
		})
	}

	return result
}

func hasTextReview(review dto.ReportReviewDto) bool {
	return strings.TrimSpace(review.Review) != ""
}

func countPositiveReviews(reviews []dto.ReportReviewDto) int {
	count := 0
	for _, review := range reviews {
		if review.Rating >= 4 {
			count++
		}
	}
	return count
}

func countLowReviews(reviews []dto.ReportReviewDto) int {
	count := 0
	for _, review := range reviews {
		if review.Rating <= 3 {
			count++
		}
	}
	return count
}

func formatRussianDate(date time.Time) string {
	months := map[time.Month]string{
		time.January:   "января",
		time.February:  "февраля",
		time.March:     "марта",
		time.April:     "апреля",
		time.May:       "мая",
		time.June:      "июня",
		time.July:      "июля",
		time.August:    "августа",
		time.September: "сентября",
		time.October:   "октября",
		time.November:  "ноября",
		time.December:  "декабря",
	}

	return fmt.Sprintf("%d %s %d", date.Day(), months[date.Month()], date.Year())
}

func formatShortDateRange(periodStart, periodEnd time.Time) string {
	return fmt.Sprintf("%s - %s", periodStart.Format("02.01.06"), periodEnd.Format("02.01.06"))
}

func weekdayNominativeRussian(weekday string) string {
	switch weekday {
	case "Monday":
		return "понедельник"
	case "Tuesday":
		return "вторник"
	case "Wednesday":
		return "среда"
	case "Thursday":
		return "четверг"
	case "Friday":
		return "пятница"
	case "Saturday":
		return "суббота"
	case "Sunday":
		return "воскресенье"
	default:
		return weekday
	}
}
