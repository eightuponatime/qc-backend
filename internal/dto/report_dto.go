package dto

type ReportVoteItemDto struct {
	MealType     string `json:"meal_type"`
	Rating       *int16 `json:"rating"`
	Review       string `json:"review"`
	BusinessDate string `json:"business_date"`
}

type ReportRatingDistributionDto struct {
	Five  int `json:"5"`
	Four  int `json:"4"`
	Three int `json:"3"`
	Two   int `json:"2"`
	One   int `json:"1"`
}

type ReportWeekdayStatsDto struct {
	Weekday            string                      `json:"weekday"`
	TotalRatings       int                         `json:"total_ratings"`
	TextReviewsCount   int                         `json:"text_reviews_count"`
	RatingDistribution ReportRatingDistributionDto `json:"rating_distribution"`
	MealStats          []ReportMealStatsDto        `json:"meal_stats"`
}

type ReportMealStatsDto struct {
	MealType        string `json:"meal_type"`
	TotalRatings    int    `json:"total_ratings"`
	LowRatingsCount int    `json:"low_ratings_count"`
}

type ReportShiftSummaryDto struct {
	ShiftType          string                      `json:"shift_type"`
	TotalRatings       int                         `json:"total_ratings"`
	TextReviewsCount   int                         `json:"text_reviews_count"`
	RatingDistribution ReportRatingDistributionDto `json:"rating_distribution"`
	WeekdayStats       []ReportWeekdayStatsDto     `json:"weekday_stats"`
	MealStats          []ReportMealStatsDto        `json:"meal_stats"`
}

type ReportSummaryDto struct {
	PeriodStart        string                      `json:"period_start"`
	PeriodEnd          string                      `json:"period_end"`
	PeriodStartDisplay string                      `json:"period_start_display"`
	PeriodEndDisplay   string                      `json:"period_end_display"`
	PeriodShortDisplay string                      `json:"period_short_display"`
	TotalVotes         int                         `json:"total_votes"`
	TotalRatings       int                         `json:"total_ratings"`
	TextReviewsCount   int                         `json:"text_reviews_count"`
	RatingDistribution ReportRatingDistributionDto `json:"rating_distribution"`
	WeekdayStats       []ReportWeekdayStatsDto     `json:"weekday_stats"`
	MealStats          []ReportMealStatsDto        `json:"meal_stats"`
	ShiftSummaries     []ReportShiftSummaryDto     `json:"shift_summaries"`
	Insights           []string                    `json:"insights"`
}

type ReportCalendarDateStatsDto struct {
	BusinessDate       string                      `json:"business_date"`
	TotalRatings       int                         `json:"total_ratings"`
	RatingDistribution ReportRatingDistributionDto `json:"rating_distribution"`
}

type ReportReviewDto struct {
	BusinessDate        string `json:"business_date"`
	BusinessDateDisplay string `json:"business_date_display"`
	BusinessWeekday     string `json:"business_weekday"`
	VoteID              string `json:"vote_id"`
	ShiftType           string `json:"shift_type"`
	MealType            string `json:"meal_type"`
	Rating              int16  `json:"rating"`
	Review              string `json:"review"`
}

type ReportDateReviewsDto struct {
	BusinessDate         string            `json:"business_date"`
	BusinessDateDisplay  string            `json:"business_date_display"`
	TotalReviews         int               `json:"total_reviews"`
	PositiveReviewsCount int               `json:"positive_reviews_count"`
	LowReviewsCount      int               `json:"low_reviews_count"`
	Reviews              []ReportReviewDto `json:"reviews"`
}

type ReportAnalyticsSummaryDto struct {
	PeriodStart            string                       `json:"period_start"`
	PeriodEnd              string                       `json:"period_end"`
	GeneratedAt            string                       `json:"generated_at"`
	Summary                ReportSummaryDto             `json:"summary"`
	CalendarDateStats      []ReportCalendarDateStatsDto `json:"calendar_date_stats"`
	AttentionRequiredItems []ReportReviewDto            `json:"attention_required_items"`
	DetailedReviewsByDate  []ReportDateReviewsDto       `json:"detailed_reviews_by_date"`
}
