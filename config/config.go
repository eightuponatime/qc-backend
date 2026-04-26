package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL              string
	Port                     string
	Env                      string
	StaticExternalIp         string
	GeoLongitude             string
	GeoLatitude              string
	BusinessTimezone         string
	SmtpHost                 string
	SmtpPort                 string
	SmtpUsername             string
	SmtpPassword             string
	SmtpFrom                 string
	ReportTo                 string
	AnalyticsURL             string
	NightShiftVoteCutoffHour int
	ReportSendHour           int
}

func Load() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file")
	}

	return &Config{
		DatabaseURL:              getEnv("DATABASE_URL", ""),
		Port:                     getEnv("PORT", "8080"),
		Env:                      getEnv("ENV", "development"),
		StaticExternalIp:         getEnv("STATIC_EXTERNAL_IP", ""),
		GeoLongitude:             getEnv("GEO_LONGITUDE", ""),
		GeoLatitude:              getEnv("GEO_LATITUDE", ""),
		BusinessTimezone:         getEnv("BUSINESS_TIMEZONE", "Asia/Almaty"),
		SmtpHost:                 getEnv("SMTP_HOST", "smtp.gmail.com"),
		SmtpPort:                 getEnv("SMTP_PORT", "587"),
		SmtpUsername:             getEnv("SMTP_USERNAME", ""),
		SmtpPassword:             getEnv("SMTP_PASSWORD", ""),
		SmtpFrom:                 getEnv("SMTP_FROM", ""),
		ReportTo:                 getEnv("REPORT_TO", ""),
		AnalyticsURL:             getEnv("ANALYTICS_URL", "http://localhost:8080/analytics"),
		NightShiftVoteCutoffHour: getEnvInt("NIGHT_SHIFT_VOTE_CUTOFF_HOUR", 8),
		ReportSendHour:           getEnvInt("REPORT_SEND_HOUR", 11),
	}, nil
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
