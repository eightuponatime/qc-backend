package handler

import (
	"html/template"
	"net/http"
	"net/url"
	"qc/config"
	"qc/internal/domain"
	"qc/internal/dto"
	"qc/internal/repository"
	"qc/internal/service"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"
)

const analyticsAccessCookieName = "analytics_access_code"

type AnalyticsHandler struct {
	accessService service.AnalyticsAccessService
	reportService service.ReportService
	sentReports   repository.SentReportRepository
	tmpl          *template.Template
	cfg           *config.Config
}

type analyticsPeriodView struct {
	PeriodStart string
	PeriodEnd   string
	Display     string
	URL         string
	IsSelected  bool
}

type analyticsPageData struct {
	Title                  string
	AnalyticsURL           string
	Periods                []analyticsPeriodView
	ShiftTabs              []analyticsShiftTabView
	SelectedShift          string
	SelectedShiftLabel     string
	SelectedShiftSummary   dto.ReportShiftSummaryDto
	AttentionRequiredItems []dto.ReportReviewDto
	DetailedReviewsByDate  []dto.ReportDateReviewsDto
	Report                 *dto.ReportAnalyticsSummaryDto
}

type analyticsShiftTabView struct {
	ShiftType    string
	Label        string
	URL          string
	TotalRatings int
	IsSelected   bool
}

type analyticsLoginPageData struct {
	Title   string
	Error   string
	Code    string
	PostURL string
}

func NewAnalyticsHandler(
	accessService service.AnalyticsAccessService,
	reportService service.ReportService,
	sentReports repository.SentReportRepository,
	tmpl *template.Template,
	cfg *config.Config,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		accessService: accessService,
		reportService: reportService,
		sentReports:   sentReports,
		tmpl:          tmpl,
		cfg:           cfg,
	}
}

func (h *AnalyticsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/analytics", h.AnalyticsPage)
	r.Post("/analytics/login", h.Login)
	r.Post("/analytics/logout", h.Logout)
}

func (h *AnalyticsHandler) AnalyticsPage(w http.ResponseWriter, r *http.Request) {
	hasAccess, err := h.hasAnalyticsAccess(r)
	if err != nil {
		http.Error(w, "failed to check analytics access", http.StatusInternalServerError)
		return
	}
	if !hasAccess {
		h.renderLogin(w, analyticsLoginPageData{
			Title:   "Доступ к аналитике",
			PostURL: "/analytics/login",
		})
		return
	}

	sentReports, err := h.sentReports.List(r.Context())
	if err != nil {
		http.Error(w, "failed to load report periods", http.StatusInternalServerError)
		return
	}

	periodStart, periodEnd, hasSelectedPeriod := selectAnalyticsPeriod(r, sentReports)

	var report *dto.ReportAnalyticsSummaryDto
	if hasSelectedPeriod {
		report, err = h.reportService.CreateAnalyticsSummaryForPeriod(r.Context(), periodStart, periodEnd)
	} else {
		report, err = h.reportService.CreateAnalyticsSummary(r.Context())
	}
	if err != nil {
		http.Error(w, "failed to build analytics report", http.StatusInternalServerError)
		return
	}

	data := analyticsPageData{
		Title:                  "Аналитика качества питания",
		AnalyticsURL:           h.cfg.AnalyticsURL,
		Periods:                buildAnalyticsPeriods(sentReports, report.PeriodStart, report.PeriodEnd, r.URL.Query().Get("shift")),
		ShiftTabs:              buildAnalyticsShiftTabs(report.Summary.ShiftSummaries, report.PeriodStart, report.PeriodEnd, r.URL.Query().Get("shift")),
		SelectedShift:          selectedAnalyticsShift(r, report.Summary.ShiftSummaries),
		SelectedShiftLabel:     analyticsShiftLabel(selectedAnalyticsShift(r, report.Summary.ShiftSummaries)),
		SelectedShiftSummary:   selectShiftSummary(report.Summary.ShiftSummaries, selectedAnalyticsShift(r, report.Summary.ShiftSummaries)),
		AttentionRequiredItems: filterReviewsByShift(report.AttentionRequiredItems, selectedAnalyticsShift(r, report.Summary.ShiftSummaries)),
		DetailedReviewsByDate:  filterDetailedReviewsByShift(report.DetailedReviewsByDate, selectedAnalyticsShift(r, report.Summary.ShiftSummaries)),
		Report:                 report,
	}

	if err := h.tmpl.ExecuteTemplate(w, "analytics.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AnalyticsHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	isValid, err := h.accessService.ValidateAccessCode(r.Context(), code)
	if err != nil {
		http.Error(w, "failed to validate access code", http.StatusInternalServerError)
		return
	}
	if !isValid {
		h.renderLogin(w, analyticsLoginPageData{
			Title:   "Доступ к аналитике",
			Error:   "Код не найден или срок действия истек.",
			Code:    code,
			PostURL: "/analytics/login",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     analyticsAccessCookieName,
		Value:    code,
		Path:     "/analytics",
		MaxAge:   60 * 60 * 24 * 15,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/analytics", http.StatusSeeOther)
}

func (h *AnalyticsHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     analyticsAccessCookieName,
		Value:    "",
		Path:     "/analytics",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/analytics", http.StatusSeeOther)
}

func (h *AnalyticsHandler) hasAnalyticsAccess(r *http.Request) (bool, error) {
	cookie, err := r.Cookie(analyticsAccessCookieName)
	if err == http.ErrNoCookie {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return h.accessService.ValidateAccessCode(r.Context(), cookie.Value)
}

func (h *AnalyticsHandler) renderLogin(w http.ResponseWriter, data analyticsLoginPageData) {
	if err := h.tmpl.ExecuteTemplate(w, "analytics_login.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func selectAnalyticsPeriod(
	r *http.Request,
	reports []domain.SentReportModel,
) (time.Time, time.Time, bool) {
	startParam := r.URL.Query().Get("period_start")
	endParam := r.URL.Query().Get("period_end")
	if startParam != "" && endParam != "" {
		periodStart, startErr := time.Parse("2006-01-02", startParam)
		periodEnd, endErr := time.Parse("2006-01-02", endParam)
		if startErr == nil && endErr == nil {
			return periodStart, periodEnd, true
		}
	}

	if len(reports) == 0 {
		return time.Time{}, time.Time{}, false
	}

	return reports[0].PeriodStart, reports[0].PeriodEnd, true
}

func buildAnalyticsPeriods(
	reports []domain.SentReportModel,
	selectedStart string,
	selectedEnd string,
	selectedShift string,
) []analyticsPeriodView {
	periods := make([]analyticsPeriodView, 0, len(reports))
	for _, report := range reports {
		periodStart := report.PeriodStart.Format("2006-01-02")
		periodEnd := report.PeriodEnd.Format("2006-01-02")
		link := buildAnalyticsLink(periodStart, periodEnd, selectedShift)
		periods = append(periods, analyticsPeriodView{
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
			Display:     formatShortPeriod(report.PeriodStart, report.PeriodEnd),
			IsSelected:  periodStart == selectedStart && periodEnd == selectedEnd,
			URL:         link,
		})
	}

	return periods
}

func selectedAnalyticsShift(r *http.Request, summaries []dto.ReportShiftSummaryDto) string {
	selected := r.URL.Query().Get("shift")
	if selected == "day" || selected == "night" {
		return selected
	}

	for _, summary := range summaries {
		if summary.ShiftType == "day" || summary.ShiftType == "night" {
			return summary.ShiftType
		}
	}

	return "day"
}

func selectShiftSummary(summaries []dto.ReportShiftSummaryDto, shiftType string) dto.ReportShiftSummaryDto {
	for _, summary := range summaries {
		if summary.ShiftType == shiftType {
			return summary
		}
	}

	return dto.ReportShiftSummaryDto{
		ShiftType:        shiftType,
		WeekdayStats:     []dto.ReportWeekdayStatsDto{},
		MealStats:        []dto.ReportMealStatsDto{},
		TotalRatings:     0,
		TextReviewsCount: 0,
	}
}

func buildAnalyticsShiftTabs(
	summaries []dto.ReportShiftSummaryDto,
	periodStart string,
	periodEnd string,
	selectedShift string,
) []analyticsShiftTabView {
	tabs := make([]analyticsShiftTabView, 0, 2)
	order := []string{"day", "night"}

	for _, shiftType := range order {
		summary := selectShiftSummary(summaries, shiftType)
		tabs = append(tabs, analyticsShiftTabView{
			ShiftType:    shiftType,
			Label:        analyticsShiftLabel(shiftType),
			TotalRatings: summary.TotalRatings,
			IsSelected:   shiftType == normalizeShiftType(selectedShift),
			URL:          buildAnalyticsLink(periodStart, periodEnd, shiftType),
		})
	}

	return tabs
}

func analyticsShiftLabel(shiftType string) string {
	if shiftType == "night" {
		return "Ночная смена"
	}
	return "Дневная смена"
}

func normalizeShiftType(shiftType string) string {
	if shiftType == "night" {
		return "night"
	}
	return "day"
}

func buildAnalyticsLink(periodStart string, periodEnd string, shiftType string) string {
	values := url.Values{}
	if periodStart != "" {
		values.Set("period_start", periodStart)
	}
	if periodEnd != "" {
		values.Set("period_end", periodEnd)
	}
	if normalized := normalizeShiftType(shiftType); normalized != "" {
		values.Set("shift", normalized)
	}

	encoded := values.Encode()
	if encoded == "" {
		return "/analytics"
	}
	return "/analytics?" + encoded
}

func filterReviewsByShift(reviews []dto.ReportReviewDto, shiftType string) []dto.ReportReviewDto {
	filtered := make([]dto.ReportReviewDto, 0, len(reviews))
	for _, review := range reviews {
		if review.ShiftType == normalizeShiftType(shiftType) {
			filtered = append(filtered, review)
		}
	}
	return filtered
}

func filterDetailedReviewsByShift(items []dto.ReportDateReviewsDto, shiftType string) []dto.ReportDateReviewsDto {
	filtered := make([]dto.ReportDateReviewsDto, 0, len(items))
	for _, item := range items {
		reviews := filterReviewsByShift(item.Reviews, shiftType)
		sort.SliceStable(reviews, func(i, j int) bool {
			if reviews[i].Rating != reviews[j].Rating {
				return reviews[i].Rating < reviews[j].Rating
			}
			return reviews[i].MealType < reviews[j].MealType
		})

		positiveCount := 0
		lowCount := 0
		for _, review := range reviews {
			if review.Rating >= 4 {
				positiveCount++
			}
			if review.Rating <= 3 {
				lowCount++
			}
		}

		filtered = append(filtered, dto.ReportDateReviewsDto{
			BusinessDate:         item.BusinessDate,
			BusinessDateDisplay:  item.BusinessDateDisplay,
			TotalReviews:         len(reviews),
			PositiveReviewsCount: positiveCount,
			LowReviewsCount:      lowCount,
			Reviews:              reviews,
		})
	}

	return filtered
}

func formatShortPeriod(periodStart, periodEnd time.Time) string {
	return periodStart.Format("02.01.06") + " - " + periodEnd.Format("02.01.06")
}
