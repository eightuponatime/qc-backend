package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"qc/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
)

type ReportHandler struct {
	reportService   service.ReportService
	dispatchService service.ReportDispatchService
}

func NewReportHander(
	reportService service.ReportService,
	dispatchService service.ReportDispatchService,
) *ReportHandler {
	return &ReportHandler{
		reportService:   reportService,
		dispatchService: dispatchService,
	}
}

func (h *ReportHandler) RegisterRoutes(r chi.Router) {
	r.Get("/checker", h.Checker)
	r.Get("/checker/summary", h.CheckerSummary)
	r.Get("/checker/analytics-summary", h.CheckerAnalyticsSummary)
	r.Post("/checker/send-demo-report", h.SendDemoReport)
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
