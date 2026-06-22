package main

import (
	"database/sql"
	"fmt"
	"log"
	"module3/config"
	"module3/internal/adapter/httphandlers"
	"module3/internal/adapter/postgres"
	"module3/internal/usecase"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	cfg := config.LoadConfig()
	r := config.NewChi(cfg)

	dbConn, err := config.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatalf("Error while connection to the database: %v", err)
	}
	defer func() {
		log.Println("Closing database connection...")
		dbConn.Close()
	}()

	performDI(r, dbConn)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on port %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func performDI(r *chi.Mux, dbConn *sql.DB) {
	urlRepository := postgres.NewPostgresRepository(dbConn)
	urlService := usecase.NewUrlService(urlRepository)
	urlHandler := httphandlers.NewUrlHandler(urlService)

	r.Route("/urls", func(r chi.Router) {
		r.Post("/", urlHandler.HandleCreateUrl)
		r.Put("/", urlHandler.HandleUpdateUrl)
		r.Delete("/{id}", urlHandler.HandleDeleteUrlById)
		r.Get("/{id}", urlHandler.HandleGetUrlById)
		r.Get("/redirect/{alias}", urlHandler.HandleRedirectByAlias)
	})
}
