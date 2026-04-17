package middleware

import (
	"context"
	"encoding/json"
	"math"
	"net"
	"net/http"
	"qc/config"
	"qc/internal/dto"
	"strconv"
)

type contextKey string

const VoteRequestKey contextKey = "vote_request"
const ExternalIpKey contextKey = "external_ip"

/*
if user uses factory Wi-Fi, it's static external ip will be equal to the static ip in .env file.
but if user's external ip doesn't equal, we ask permission to geolocation and get user's longitude
and latitude, then compare with the ones stored in .env file (real coord of the factory)
*/
func AuthRequired(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var req dto.VoteRequestDto

			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "json parsing error", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			if req.Latitude == nil && req.Longitude == nil {
				externalIp := extractExternalIp(r)
				var compareIp bool = compareIps(cfg, externalIp)
				if !compareIp {
					http.Error(w, `{"error": "invalid external ip"`, http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), VoteRequestKey, req)
				ctx = context.WithValue(ctx, ExternalIpKey, extractExternalIp(r))
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				dist := getDistance(
					cfg.GeoLatitude,
					cfg.GeoLongitude,
					req.Latitude,
					req.Longitude,
				)

				if dist <= 2000 {
					ctx := context.WithValue(r.Context(), VoteRequestKey, req)
					next.ServeHTTP(w, r.WithContext(ctx))
				} else {
					http.Error(w, `{"error": "out of the allowed area"}`, http.StatusUnauthorized)
				}
			}
		})
	}
}

func getDistance(factoryLatStr, factoryLonStr string, userLatStr, userLonStr *string) float64 {
	const R = 6371000

	coords := convertCoordsToFloat64(factoryLatStr, factoryLonStr, userLatStr, userLonStr)

	lat1 := coords.factoryLat * math.Pi / 180
	lon1 := coords.factoryLon * math.Pi / 180
	lat2 := coords.userLat * math.Pi / 180
	lon2 := coords.userLon * math.Pi / 180

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return distance
}

func convertCoordsToFloat64(
	factoryLatStr string,
	factoryLonStr string,
	userLatStr *string,
	userLonStr *string,
) Coords {
	var coords Coords
	temp, _ := strconv.ParseFloat(factoryLatStr, 64)
	coords.factoryLat = float64(temp)
	temp, _ = strconv.ParseFloat(factoryLonStr, 64)
	coords.factoryLon = float64(temp)
	temp, _ = strconv.ParseFloat(*userLatStr, 64)
	coords.userLat = float64(temp)
	temp, _ = strconv.ParseFloat(*userLonStr, 64)
	coords.userLon = float64(temp)

	return coords
}

type Coords struct {
	factoryLat float64
	factoryLon float64
	userLat    float64
	userLon    float64
}

func compareIps(cfg *config.Config, externalIp string) bool {
	staticIp := cfg.StaticExternalIp
	var compareIp bool = staticIp == externalIp
	return compareIp
}

func extractExternalIp(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")

	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		ip = host
	}

	return ip
}
