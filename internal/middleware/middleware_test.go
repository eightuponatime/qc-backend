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

func TestAuthRequired_InvalidExternalIpStillAllowed(t *testing.T) {
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

	if !called {
		t.Error("next should be called for invalid IP when IP restriction is disabled")
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

func TestAuthRequired_OutOfAreaStillAllowed(t *testing.T) {
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

	if !called {
		t.Error("next should be called for out of area when IP restriction is disabled")
	}
}
