package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"qc/config"
	"qc/internal/handler"
	appLogger "qc/internal/logger"
	appMiddleware "qc/internal/middleware"
	"qc/internal/repository/postgres"
	"qc/internal/service/impl"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"qc/internal/i18n"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := appLogger.Setup(cfg)

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connected")

	err = loadManifest()
	if err != nil {
		logger.Error("manifest load failed", slog.Any("error", err))
		os.Exit(1)
	}
	initTemplates()

	// repository
	voteRepo := postgres.NewVoteRepository(db)
	reportRepo := postgres.NewReportRepository(db)
	sentReportRepo := postgres.NewSentReportRepository(db)
	analyticsAccessRepo := postgres.NewAnalyticsAccessRepository(db)
	txManager := postgres.NewTransactionManager(db)

	// service
	voteService := impl.NewVoteService(voteRepo, txManager, cfg)
	reportService := impl.NewReportService(reportRepo, cfg)
	analyticsAccessService := impl.NewAnalyticsAccessService(analyticsAccessRepo, cfg)
	emailService := impl.NewEmailService(reportService, analyticsAccessService, cfg)
	reportScheduler := impl.NewReportScheduler(emailService, sentReportRepo, cfg, time.Hour)

	// handler
	voteHandler := handler.NewVoteHandler(voteService, tmpl, cfg)
	reportHandler := handler.NewReportHander(reportService)
	analyticsHandler := handler.NewAnalyticsHandler(
		analyticsAccessService,
		reportService,
		sentReportRepo,
		tmpl,
		cfg,
	)

	// middleware
	authRequired := appMiddleware.AuthRequired(cfg)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Authorization"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
	}))

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		lang := i18n.DetectLanguage(r)

		http.SetCookie(w, &http.Cookie{
			Name:     "lang",
			Value:    lang,
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 365,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		})

		translations, err := i18n.Load(lang)
		if err != nil {
			lang = "ru"
			translations, err = i18n.Load(lang)
			if err != nil {
				http.Error(w, "failed to load translations", http.StatusInternalServerError)
				return
			}
		}

		err = tmpl.ExecuteTemplate(w, "base.html", map[string]any{
			"Title": "Home",
			"Lang":  lang,
			"T":     translations,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	voteHandler.RegisterRoutes(r, authRequired)
	reportHandler.RegisterRoutes(r)
	analyticsHandler.RegisterRoutes(r)

	go reportScheduler.Start(context.Background())

	logger.Info("server starting", slog.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Error("server stopped", slog.Any("error", err))
		os.Exit(1)
	}
}
