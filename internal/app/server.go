package app

import (
	"log"
	"net/http"

	"github.com/Niiazgulov/bwg.git/internal/config"
	"github.com/Niiazgulov/bwg.git/internal/handlers"
	"github.com/Niiazgulov/bwg.git/internal/storage"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Start() {
	if err := env.Parse(&config.Cfg); err != nil {
		log.Fatal(err)
	}
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	config.Cfg = *cfg
	repo, err := storage.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	router := Route(repo)
	log.Fatal(http.ListenAndServe(config.Cfg.ServerAddress, router))
}

func Route(repo storage.Transaction) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetBalanceHandler(repo))
		r.Post("/invoice", handlers.PostInvoiceHandler(repo))
		r.Post("/withdraw", handlers.PostWithdrawHandler(repo))
	})
	return r
}
