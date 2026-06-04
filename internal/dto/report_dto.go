package dto

type SendReportToRequestDto struct {
	Email       string `json:"email"`
	PeriodStart string `json:"period_start,omitempty"`
	PeriodEnd   string `json:"period_end,omitempty"`
}

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

type ReportDateStatDto struct {
	BusinessDate        string  `json:"business_date"`
	BusinessDateDisplay string  `json:"business_date_display"`
	BusinessWeekday     string  `json:"business_weekday"`
	TotalRatings        int     `json:"total_ratings"`
	TextReviewsCount    int     `json:"text_reviews_count"`
	AverageRating       float64 `json:"average_rating"`
}

type ReportShiftSummaryDto struct {
	ShiftType          string                      `json:"shift_type"`
	TotalRatings       int                         `json:"total_ratings"`
	TextReviewsCount   int                         `json:"text_reviews_count"`
	RatingDistribution ReportRatingDistributionDto `json:"rating_distribution"`
	WeekdayStats       []ReportWeekdayStatsDto     `json:"weekday_stats"`
	DateStats          []ReportDateStatDto         `json:"date_stats"`
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

type ReportMealReviewsDto struct {
	MealType           string                      `json:"meal_type"`
	TotalRatings       int                         `json:"total_ratings"`
	TextReviewsCount   int                         `json:"text_reviews_count"`
	RatingDistribution ReportRatingDistributionDto `json:"rating_distribution"`
	AverageRating      float64                     `json:"average_rating"`
	Reviews            []ReportReviewDto           `json:"reviews"`
}

type ReportDateReviewsDto struct {
	BusinessDate         string                 `json:"business_date"`
	BusinessDateDisplay  string                 `json:"business_date_display"`
	BusinessWeekday      string                 `json:"business_weekday"`
	TotalRatings         int                    `json:"total_ratings"`
	AverageRating        float64                `json:"average_rating"`
	TotalReviews         int                    `json:"total_reviews"`
	PositiveReviewsCount int                    `json:"positive_reviews_count"`
	LowReviewsCount      int                    `json:"low_reviews_count"`
	Meals                []ReportMealReviewsDto `json:"meals"`
	Reviews              []ReportReviewDto      `json:"reviews"`
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
