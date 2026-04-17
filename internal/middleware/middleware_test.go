package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"qc/config"
	"qc/internal/dto"
	"qc/internal/middleware"
	"testing"
)

func TestAuthRequired_ValidExternalIp(t *testing.T) {
	cfg := &config.Config{
		StaticExternalIp: "1.2.3.4",
		GeoLatitude:      "43.0",
		GeoLongitude:     "76.0",
	}

	body, _ := json.Marshal(dto.VoteRequestDto{
		DeviceId: "test-device",
	})

	r := httptest.NewRequest(http.MethodPost, "/api/vote", bytes.NewReader(body))
	r.Header.Set("X-Real-IP", "1.2.3.4")
	w := httptest.NewRecorder()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware.AuthRequired(cfg)(next).ServeHTTP(w, r)

	if !called {
		t.Error("next should be called for valid IP")
	}
}

func TestAuthRequired_InvalidExternalIp(t *testing.T) {
	cfg := &config.Config{
		StaticExternalIp: "1.2.3.4",
		GeoLatitude:      "43.0",
		GeoLongitude:     "76.0",
	}

	body, _ := json.Marshal(dto.VoteRequestDto{
		DeviceId: "test-device",
	})

	r := httptest.NewRequest(http.MethodPost, "/api/vote", bytes.NewReader(body))
	r.Header.Set("X-Real-IP", "9.9.9.9")
	w := httptest.NewRecorder()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware.AuthRequired(cfg)(next).ServeHTTP(w, r)

	if called {
		t.Error("next should not be called for invalid IP")
	}
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_ValidGeoLocation(t *testing.T) {
	cfg := &config.Config{
		StaticExternalIp: "1.2.3.4",
		GeoLatitude:      "43.0",
		GeoLongitude:     "76.0",
	}

	lat := "43.001"
	lon := "76.001"
	body, _ := json.Marshal(dto.VoteRequestDto{
		DeviceId:  "test-device",
		Latitude:  &lat,
		Longitude: &lon,
	})

	r := httptest.NewRequest(http.MethodPost, "/api/vote", bytes.NewReader(body))
	r.Header.Set("X-Real-IP", "9.9.9.9")
	w := httptest.NewRecorder()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware.AuthRequired(cfg)(next).ServeHTTP(w, r)

	if !called {
		t.Error("next should be called for valid geolocation")
	}
}

func TestAuthRequired_OutOfArea(t *testing.T) {
	cfg := &config.Config{
		StaticExternalIp: "1.2.3.4",
		GeoLatitude:      "43.0",
		GeoLongitude:     "76.0",
	}

	lat := "55.75"
	lon := "37.61"
	body, _ := json.Marshal(dto.VoteRequestDto{
		DeviceId:  "test-device",
		Latitude:  &lat,
		Longitude: &lon,
	})

	r := httptest.NewRequest(http.MethodPost, "/api/vote", bytes.NewReader(body))
	r.Header.Set("X-Real-IP", "9.9.9.9")
	w := httptest.NewRecorder()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware.AuthRequired(cfg)(next).ServeHTTP(w, r)

	if called {
		t.Error("next should not be called for out of area")
	}
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
