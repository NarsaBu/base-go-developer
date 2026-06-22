package config

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewChi(cfg *Config) *chi.Mux {
	r := chi.NewRouter()
	creds := map[string]string{
		cfg.Authorization.Username: cfg.Authorization.Password,
	}

	r.Use(middleware.BasicAuth("Restricted Area", creds))

	return r
}
