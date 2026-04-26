package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"qc/config"
	"qc/internal/dto"
	"qc/internal/i18n"
	"qc/internal/middleware"
	"qc/internal/service"

	"github.com/go-chi/chi/v5"
)

type VoteHandler struct {
	voteService service.VoteService
	tmpl        *template.Template
	cfg         *config.Config
}

func NewVoteHandler(
	voteService service.VoteService,
	tmpl *template.Template,
	cfg *config.Config,
) *VoteHandler {
	return &VoteHandler{
		voteService: voteService,
		tmpl:        tmpl,
		cfg:         cfg,
	}
}

func (v *VoteHandler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler) {
	r.Get("/fragments/vote-ui", v.GetVoteUIFragment)
	r.With(authRequired).Post("/fragments/vote", v.SubmitVoteFragment)
}

//
// API (JSON)
//

func (v *VoteHandler) Vote(w http.ResponseWriter, r *http.Request) {
	var req dto.VoteRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	externalIp := extractIp(r)

	if err := v.voteService.CreateVote(r.Context(), req, externalIp); err != nil {
		slog.Error(
			"vote save failed",
			slog.String("channel", "api"),
			slog.String("device_id", req.DeviceId),
			slog.String("external_ip", externalIp),
			slog.Any("error", err),
		)
		http.Error(w, `{"error":"failed to save vote"}`, http.StatusInternalServerError)
		return
	}

	slog.Info(
		"vote submitted",
		slog.String("channel", "api"),
		slog.String("device_id", req.DeviceId),
		slog.String("external_ip", externalIp),
		slog.Int("items_count", len(req.Items)),
	)

	w.WriteHeader(http.StatusOK)
}

func (v *VoteHandler) GetTodayVote(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("device_id")
	if deviceId == "" {
		http.Error(w, `{"error":"device_id is required"}`, http.StatusBadRequest)
		return
	}

	vote, err := v.voteService.GetTodayVote(r.Context(), deviceId, requestedShiftType(r.URL.Query().Get("shift_type")))
	if err != nil {
		slog.Error("get today vote failed", slog.String("device_id", deviceId), slog.Any("error", err))
		http.Error(w, `{"error":"failed to get vote"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vote)
}

//
// HTMX UI
//

type VoteUIData struct {
	DeviceId   string
	ShiftType  string
	PhoneModel string
	Browser    string

	Vote *dto.VoteResponseDto

	SuccessMessage      string
	ErrorMessage        string
	ActiveMeal          string
	VotingClosed        bool
	VotingClosedMessage string

	AccessRestricted  bool
	AccessReason      string
	AccessDismissible bool

	PrefillMealType string
	PrefillRating   string
	PrefillReview   string

	ServerNowISO string
	Lang         string
	T            i18n.Translations

	BreakfastAvailability MealAvailability
	LunchAvailability     MealAvailability
	DinnerAvailability    MealAvailability
}

type MealAvailability struct {
	IsAvailable      bool
	OpensAt          string
	MinutesUntilOpen int
}

func (v *VoteHandler) GetVoteUIFragment(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("device_id")
	if deviceId == "" {
		translations := v.loadTranslations(r)
		http.Error(w, translateError(translations, "device_id_required"), http.StatusBadRequest)
		return
	}

	requestedShiftType := requestedShiftType(r.URL.Query().Get("shift_type"))
	vote, err := v.voteService.GetTodayVote(r.Context(), deviceId, requestedShiftType)
	if err != nil {
		slog.Error("vote ui lookup failed", slog.String("device_id", deviceId), slog.Any("error", err))

		translations := v.loadTranslations(r)
		http.Error(w, translateError(translations, "vote_lookup_failed"), http.StatusInternalServerError)
		return
	}

	translations := v.loadTranslations(r)
	lang := i18n.DetectLanguage(r)

	data := VoteUIData{
		DeviceId:        deviceId,
		ShiftType:       detectShiftType(vote, requestedShiftType),
		PhoneModel:      r.URL.Query().Get("phone_model"),
		Browser:         r.URL.Query().Get("browser"),
		Vote:            vote,
		T:               translations,
		Lang:            lang,
		PrefillMealType: r.URL.Query().Get("meal_type"),
		PrefillRating:   r.URL.Query().Get("rating"),
		PrefillReview:   r.URL.Query().Get("review"),
		ActiveMeal:      r.URL.Query().Get("meal_type"),
	}
	data.VotingClosed, data.VotingClosedMessage = v.getVotingClosedState(data.ShiftType, translations)

	v.renderVoteUI(w, data, http.StatusOK)
}

func (v *VoteHandler) SubmitVoteFragment(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		slog.Error("vote fragment parse form failed", slog.Any("error", err))

		translations := v.loadTranslations(r)
		v.renderErrorUI(
			w,
			r,
			translateError(translations, "invalid_form"),
		)
		return
	}

	translations := v.loadTranslations(r)

	authErrCode, _ := r.Context().Value(middleware.AuthErrorKey).(string)
	if authErrCode != "" {
		if authErrCode == "invalid_external_ip" {
			v.renderRestrictedUI(w, r, "corporate_wifi_required")
			return
		}

		v.renderErrorUI(
			w,
			r,
			translateError(translations, authErrCode),
		)
		return
	}

	req, err := buildVoteRequestFromForm(r)
	if err != nil {
		code := formErrorCode(err)

		slog.Warn("vote fragment validation failed", slog.String("code", code), slog.Any("error", err))

		v.renderErrorUI(
			w,
			r,
			translateError(translations, code),
		)
		return
	}

	externalIp, _ := r.Context().Value(middleware.ExternalIpKey).(string)
	if externalIp == "" {
		externalIp = extractIp(r)
	}

	if err := v.voteService.CreateVote(r.Context(), req, externalIp); err != nil {
		code := serviceErrorCode(err)

		slog.Error(
			"vote fragment save failed",
			slog.String("code", code),
			slog.String("device_id", req.DeviceId),
			slog.String("shift_type", req.ShiftType),
			slog.String("meal_type", firstMealType(req)),
			slog.String("external_ip", externalIp),
			slog.Any("error", err),
		)

		if code == "voting_closed" {
			v.renderVotingClosedUI(w, r, req.ShiftType)
			return
		}

		v.renderErrorUI(w, r, translateError(translations, code))
		return
	}

	vote, err := v.voteService.GetTodayVote(r.Context(), req.DeviceId, req.ShiftType)
	if err != nil {
		slog.Error("vote fragment reload failed", slog.String("device_id", req.DeviceId), slog.Any("error", err))

		v.renderErrorUI(
			w,
			r,
			translateError(translations, "vote_reload_failed"),
		)
		return
	}

	lang := i18n.DetectLanguage(r)
	activeMeal := firstMealType(req)

	slog.Info(
		"vote submitted",
		slog.String("channel", "fragment"),
		slog.String("device_id", req.DeviceId),
		slog.String("shift_type", req.ShiftType),
		slog.String("meal_type", activeMeal),
		slog.String("external_ip", externalIp),
		slog.String("phone_model", req.PhoneModel),
		slog.String("browser", req.Browser),
	)

	v.renderVoteUI(w, VoteUIData{
		DeviceId:       req.DeviceId,
		ShiftType:      req.ShiftType,
		PhoneModel:     req.PhoneModel,
		Browser:        req.Browser,
		Vote:           vote,
		SuccessMessage: translateUI(translations, "saved", "Сохранено"),
		ActiveMeal:     activeMeal,
		T:              translations,
		Lang:           lang,
	}, http.StatusOK)
}

//
// RENDER HELPERS
//

func (v *VoteHandler) renderVoteUI(w http.ResponseWriter, data VoteUIData, status int) {
	breakfastAvailability, lunchAvailability, dinnerAvailability, serverNow, err := getMealAvailabilities(
		time.Now(),
		v.cfg.BusinessTimezone,
	)
	if err == nil {
		data.BreakfastAvailability = breakfastAvailability
		data.LunchAvailability = lunchAvailability
		data.DinnerAvailability = dinnerAvailability
		data.ServerNowISO = serverNow.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if err := v.tmpl.ExecuteTemplate(w, "vote_ui.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (v *VoteHandler) renderErrorUI(w http.ResponseWriter, r *http.Request, message string) {
	activeMeal := r.FormValue("meal_type")
	translations := v.loadTranslations(r)
	lang := i18n.DetectLanguage(r)
	shiftType := requestedShiftType(r.FormValue("shift_type"))

	data := VoteUIData{
		DeviceId:     r.FormValue("device_id"),
		ShiftType:    shiftType,
		PhoneModel:   r.FormValue("phone_model"),
		Browser:      r.FormValue("browser"),
		ErrorMessage: message,
		ActiveMeal:   activeMeal,
		T:            translations,
		Lang:         lang,
	}
	data.VotingClosed, data.VotingClosedMessage = v.getVotingClosedState(shiftType, translations)

	v.renderVoteUI(w, data, http.StatusOK)
}

func (v *VoteHandler) renderRestrictedUI(w http.ResponseWriter, r *http.Request, reason string) {
	activeMeal := r.FormValue("meal_type")
	translations := v.loadTranslations(r)
	lang := i18n.DetectLanguage(r)

	data := VoteUIData{
		DeviceId:          r.FormValue("device_id"),
		ShiftType:         r.FormValue("shift_type"),
		PhoneModel:        r.FormValue("phone_model"),
		Browser:           r.FormValue("browser"),
		ActiveMeal:        activeMeal,
		AccessRestricted:  true,
		AccessReason:      reason,
		AccessDismissible: true,
		PrefillMealType:   r.FormValue("meal_type"),
		PrefillRating:     r.FormValue("rating"),
		PrefillReview:     r.FormValue("review"),
		T:                 translations,
		Lang:              lang,
	}

	v.renderVoteUI(w, data, http.StatusOK)
}

func (v *VoteHandler) renderVotingClosedUI(w http.ResponseWriter, r *http.Request, shiftType string) {
	translations := v.loadTranslations(r)
	lang := i18n.DetectLanguage(r)
	normalizedShiftType := requestedShiftType(shiftType)
	_, message := v.getVotingClosedState(normalizedShiftType, translations)

	data := VoteUIData{
		DeviceId:            r.FormValue("device_id"),
		ShiftType:           normalizedShiftType,
		PhoneModel:          r.FormValue("phone_model"),
		Browser:             r.FormValue("browser"),
		ActiveMeal:          r.FormValue("meal_type"),
		VotingClosed:        true,
		VotingClosedMessage: message,
		T:                   translations,
		Lang:                lang,
	}

	v.renderVoteUI(w, data, http.StatusOK)
}

//
// FORM PARSING
//

func buildVoteRequestFromForm(r *http.Request) (dto.VoteRequestDto, error) {
	deviceId := r.FormValue("device_id")
	if deviceId == "" {
		return dto.VoteRequestDto{}, fmt.Errorf("device_id is required")
	}

	mealType := r.FormValue("meal_type")
	if mealType == "" {
		return dto.VoteRequestDto{}, fmt.Errorf("meal_type is required")
	}

	item, err := buildSingleMealItemFromForm(r, mealType)
	if err != nil {
		return dto.VoteRequestDto{}, err
	}

	return dto.VoteRequestDto{
		DeviceId:   deviceId,
		ShiftType:  r.FormValue("shift_type"),
		PhoneModel: r.FormValue("phone_model"),
		Browser:    r.FormValue("browser"),
		Items:      []dto.VoteMealItemDto{item},
	}, nil
}

func buildSingleMealItemFromForm(r *http.Request, mealType string) (dto.VoteMealItemDto, error) {
	ratingStr := r.FormValue("rating")
	reviewStr := r.FormValue("review")

	if ratingStr == "" && reviewStr == "" {
		return dto.VoteMealItemDto{}, fmt.Errorf("rating or review is required")
	}

	var rating *int16
	if ratingStr != "" {
		val, err := strconv.ParseInt(ratingStr, 10, 16)
		if err != nil {
			return dto.VoteMealItemDto{}, fmt.Errorf("%s rating invalid", mealType)
		}
		tmp := int16(val)
		rating = &tmp
	}

	var review *string
	if reviewStr != "" {
		review = &reviewStr
	}

	return dto.VoteMealItemDto{
		MealType: mealType,
		Rating:   rating,
		Review:   review,
	}, nil
}

//
// UTIL
//

func firstMealType(req dto.VoteRequestDto) string {
	if len(req.Items) == 0 {
		return ""
	}
	return req.Items[0].MealType
}

func getMealAvailabilities(now time.Time, timezone string) (MealAvailability, MealAvailability, MealAvailability, time.Time, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return MealAvailability{}, MealAvailability{}, MealAvailability{}, time.Time{}, err
	}

	localNow := now.In(location)

	buildAvailability := func(hour int, minute int) MealAvailability {
		target := time.Date(
			localNow.Year(),
			localNow.Month(),
			localNow.Day(),
			hour,
			minute,
			0,
			0,
			location,
		)

		diff := target.Sub(localNow)
		minutesUntilOpen := 0
		if diff > 0 {
			minutesUntilOpen = int(diff.Minutes())
			if diff%time.Minute != 0 {
				minutesUntilOpen++
			}
		}

		return MealAvailability{
			IsAvailable:      !localNow.Before(target),
			OpensAt:          target.Format("15:04"),
			MinutesUntilOpen: minutesUntilOpen,
		}
	}

	breakfast := buildAvailability(6, 0)
	lunch := buildAvailability(12, 0)
	dinner := buildAvailability(17, 0)

	return breakfast, lunch, dinner, localNow, nil
}

func (v *VoteHandler) getVotingClosedState(shiftType string, translations i18n.Translations) (bool, string) {
	location, err := time.LoadLocation(v.cfg.BusinessTimezone)
	if err != nil {
		return false, ""
	}

	SHIFT_DATE := 15
	BORDER_DATE := 16

	now := time.Now().In(location)
	normalizedShiftType := requestedShiftType(shiftType)

	if now.Day() >= 1 && now.Day() <= SHIFT_DATE {
		return false, ""
	}

	if normalizedShiftType == "night" && now.Day() == BORDER_DATE && now.Hour() < v.cfg.NightShiftVoteCutoffHour {
		return false, ""
	}

	return true, buildVotingClosedMessage(now, i18nLanguageFromTranslations(translations))
}

func extractIp(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// ======
// localization
// ======

func translateError(t i18n.Translations, code string) string {
	key := "errors." + code
	if msg, ok := t[key]; ok {
		return msg
	}
	return "unknown error: " + code
}

func translateUI(t i18n.Translations, code string, fallback string) string {
	key := "ui." + code
	if msg, ok := t[key]; ok {
		return msg
	}
	return fallback
}

func (v *VoteHandler) loadTranslations(r *http.Request) i18n.Translations {
	lang := i18n.DetectLanguage(r)

	translations, err := i18n.Load(lang)
	if err != nil {
		translations, err = i18n.Load("ru")
		if err != nil {
			return i18n.Translations{}
		}
	}

	return translations
}

func formErrorCode(err error) string {
	if err == nil {
		return ""
	}

	switch err.Error() {
	case "device_id is required":
		return "device_id_required"
	case "meal_type is required":
		return "meal_type_required"
	case "rating or review is required":
		return "rating_or_review_required"
	default:
		if strings.Contains(err.Error(), "rating invalid") {
			return "invalid_rating"
		}
		return "invalid_form_data"
	}
}

func serviceErrorCode(err error) string {
	if err == nil {
		return ""
	}

	msg := err.Error()

	switch {
	case strings.Contains(msg, "invalid shift type"):
		return "invalid_form_data"
	case strings.Contains(msg, "invalid meal type"):
		return "invalid_meal_type"
	case strings.Contains(msg, "voting is closed"):
		return "voting_closed"
	case strings.Contains(msg, "not available yet"):
		return "meal_not_available_yet"
	case strings.Contains(msg, "load business timezone"):
		return "business_timezone_error"
	case strings.Contains(msg, "get business date"):
		return "business_date_error"
	case strings.Contains(msg, "get vote by day"):
		return "vote_lookup_failed"
	case strings.Contains(msg, "create vote"):
		return "vote_create_failed"
	case strings.Contains(msg, "upsert vote item"):
		return "vote_item_save_failed"
	case strings.Contains(msg, "get vote items"):
		return "vote_items_lookup_failed"
	default:
		return "vote_save_failed"
	}
}

func detectShiftType(vote *dto.VoteResponseDto, requestedShiftType string) string {
	if vote != nil && vote.ShiftType != "" {
		return vote.ShiftType
	}
	if requestedShiftType == "night" {
		return "night"
	}
	return "day"
}

func requestedShiftType(shiftType string) string {
	if shiftType == "night" {
		return "night"
	}
	return "day"
}

func buildVotingClosedMessage(now time.Time, lang string) string {
	switch lang {
	case "en":
		return fmt.Sprintf(
			"Voting for the period from 1 to 16 %s is already closed.",
			englishMonthName(now.Month()),
		)
	case "kk":
		return fmt.Sprintf(
			"%s айының 1-інен 16-сына дейінгі кезең бойынша дауыс беру жабылды.",
			kazakhMonthName(now.Month()),
		)
	default:
		return fmt.Sprintf(
			"Голосование за период с 1 по 16 %s уже закрыто.",
			russianMonthGenitive(now.Month()),
		)
	}
}

func i18nLanguageFromTranslations(translations i18n.Translations) string {
	title := translations["title"]
	switch title {
	case "Quality Control":
		return "en"
	case "Сапаны бақылау":
		return "kk"
	default:
		return "ru"
	}
}

func russianMonthGenitive(month time.Month) string {
	switch month {
	case time.January:
		return "января"
	case time.February:
		return "февраля"
	case time.March:
		return "марта"
	case time.April:
		return "апреля"
	case time.May:
		return "мая"
	case time.June:
		return "июня"
	case time.July:
		return "июля"
	case time.August:
		return "августа"
	case time.September:
		return "сентября"
	case time.October:
		return "октября"
	case time.November:
		return "ноября"
	case time.December:
		return "декабря"
	default:
		return ""
	}
}

func englishMonthName(month time.Month) string {
	return month.String()
}

func kazakhMonthName(month time.Month) string {
	switch month {
	case time.January:
		return "қаңтар"
	case time.February:
		return "ақпан"
	case time.March:
		return "наурыз"
	case time.April:
		return "сәуір"
	case time.May:
		return "мамыр"
	case time.June:
		return "маусым"
	case time.July:
		return "шілде"
	case time.August:
		return "тамыз"
	case time.September:
		return "қыркүйек"
	case time.October:
		return "қазан"
	case time.November:
		return "қараша"
	case time.December:
		return "желтоқсан"
	default:
		return ""
	}
}
