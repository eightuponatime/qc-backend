package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"qc/internal/dto"
	"qc/internal/service"
	"log"

	"github.com/go-chi/chi/v5"
)

type VoteHandler struct {
	voteService service.VoteService
}

func NewVoteHandler(voteService service.VoteService) *VoteHandler {
	return &VoteHandler{voteService: voteService}
}

func (v *VoteHandler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler) {
	r.Route("/api/vote", func(r chi.Router) {
		r.With(authRequired).Post("/", v.Vote)
		r.Get("/today", v.GetTodayVote)
	})
}

func (v *VoteHandler) Vote(w http.ResponseWriter, r *http.Request) {
	var req dto.VoteRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	externalIp := extractIp(r)

	if err := v.voteService.CreateVote(r.Context(), req, externalIp); err != nil {
		http.Error(w, `{"error":"failed to save vote"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (v *VoteHandler) GetTodayVote(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("device_id")
	if deviceId == "" {
		http.Error(w, `{"error":"device_id is required"}`, http.StatusBadRequest)
		return
	}

	vote, err := v.voteService.GetTodayVote(r.Context(), deviceId)
	if err != nil {
		log.Printf("GetTodayVote error: %+v", err)
		http.Error(w, `{"error":"failed to get vote"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vote)
}

func extractIp(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
