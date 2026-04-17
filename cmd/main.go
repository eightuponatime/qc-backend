package main

import (
	"log"
	"net/http"
	"qc/config"
	"qc/internal/handler"
	appMiddleware "qc/internal/middleware"
	"qc/internal/repository/postgres"
	"qc/internal/service/impl"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("connected to database")

	err = loadManifest()
	if err != nil {
		log.Fatalf("failed to load manifest: %v", err)
	}
	initTemplates()

	// repository
	voteRepo := postgres.NewVoteRepository(db)
	txManager := postgres.NewTransactionManager(db)

	// service
	voteService := impl.NewVoteService(voteRepo, txManager, cfg)

	// handler
	voteHandler := handler.NewVoteHandler(voteService)

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
		err := tmpl.ExecuteTemplate(w, "base.html", map[string]any{
			"Title": "Home",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	voteHandler.RegisterRoutes(r, authRequired)

	log.Printf("starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}