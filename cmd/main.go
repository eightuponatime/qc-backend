package main

import (
	"log"
	"net/http"
	"qc/config"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	err = loadManifest()
	if err != nil {
		panic(err)
	}

	initTemplates()

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "index.html", map[string]any{
			"Title": "Home",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) 
		}
	})

	log.Printf("starting server api on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
