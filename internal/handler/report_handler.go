package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"qc/internal/service"

	"github.com/go-chi/chi/v5"
)

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHander(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

func (h *ReportHandler) RegisterRoutes(r chi.Router) {
	r.Get("/checker", h.Checker)
	r.Get("/checker/summary", h.CheckerSummary)
	r.Get("/checker/analytics-summary", h.CheckerAnalyticsSummary)
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
