package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"qc/internal/dto"
	"qc/internal/repository"
	"qc/internal/service"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

var errInvalidReportPeriod = errors.New("period_start and period_end must be valid YYYY-MM-DD dates")
var errNoSentReports = errors.New("no sent reports found")

type ReportHandler struct {
	reportService   service.ReportService
	dispatchService service.ReportDispatchService
	sentReportRepo  repository.SentReportRepository
}

func NewReportHander(
	reportService service.ReportService,
	dispatchService service.ReportDispatchService,
	sentReportRepo repository.SentReportRepository,
) *ReportHandler {
	return &ReportHandler{
		reportService:   reportService,
		dispatchService: dispatchService,
		sentReportRepo:  sentReportRepo,
	}
}

func (h *ReportHandler) RegisterRoutes(r chi.Router) {
	r.Get("/checker", h.Checker)
	r.Get("/checker/summary", h.CheckerSummary)
	r.Get("/checker/analytics-summary", h.CheckerAnalyticsSummary)
	r.Post("/checker/send-demo-report", h.SendDemoReport)
	r.Post("/checker/send-report-to", h.SendReportTo)
}

func (h *ReportHandler) Checker(w http.ResponseWriter, r *http.Request) {
	resp, err := h.reportService.CreateReport(r.Context())
	if err != nil {
		log.Printf("create report: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ReportHandler) CheckerSummary(w http.ResponseWriter, r *http.Request) {
	resp, err := h.reportService.CreateSummary(r.Context())
	if err != nil {
		log.Printf("create report summary: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ReportHandler) CheckerAnalyticsSummary(w http.ResponseWriter, r *http.Request) {
	resp, err := h.reportService.CreateAnalyticsSummary(r.Context())
	if err != nil {
		log.Printf("create analytics report summary: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ReportHandler) SendDemoReport(w http.ResponseWriter, r *http.Request) {
	periodStart, err := time.Parse("2006-01-02", "2026-04-01")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	periodEnd, err := time.Parse("2006-01-02", "2026-04-15")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := h.dispatchService.SendPeriodReport(r.Context(), periodStart, periodEnd); err != nil {
		log.Printf("send demo period report: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":       "ok",
		"period_start": "2026-04-01",
		"period_end":   "2026-04-15",
	})
}

func (h *ReportHandler) SendReportTo(w http.ResponseWriter, r *http.Request) {
	var req dto.SendReportToRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	address, err := mail.ParseAddress(email)
	if err != nil {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	periodStart, periodEnd, err := h.resolveRequestedReportPeriod(r, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.dispatchService.SendPeriodReportTo(
		r.Context(),
		periodStart,
		periodEnd,
		[]string{address.Address},
	); err != nil {
		log.Printf("send report to %s: %v", address.Address, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":       "ok",
		"email":        address.Address,
		"period_start": periodStart.Format("2006-01-02"),
		"period_end":   periodEnd.Format("2006-01-02"),
	})
}

func (h *ReportHandler) resolveRequestedReportPeriod(r *http.Request, req dto.SendReportToRequestDto) (time.Time, time.Time, error) {
	periodStartRaw := strings.TrimSpace(req.PeriodStart)
	periodEndRaw := strings.TrimSpace(req.PeriodEnd)

	if periodStartRaw == "" && periodEndRaw == "" {
		reports, err := h.sentReportRepo.List(r.Context())
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		if len(reports) == 0 {
			return time.Time{}, time.Time{}, errNoSentReports
		}

		return reports[0].PeriodStart, reports[0].PeriodEnd, nil
	}

	if periodStartRaw == "" || periodEndRaw == "" {
		return time.Time{}, time.Time{}, errInvalidReportPeriod
	}

	periodStart, err := time.Parse("2006-01-02", periodStartRaw)
	if err != nil {
		return time.Time{}, time.Time{}, errInvalidReportPeriod
	}

	periodEnd, err := time.Parse("2006-01-02", periodEndRaw)
	if err != nil {
		return time.Time{}, time.Time{}, errInvalidReportPeriod
	}

	if periodEnd.Before(periodStart) {
		return time.Time{}, time.Time{}, errInvalidReportPeriod
	}

	return periodStart, periodEnd, nil
}
