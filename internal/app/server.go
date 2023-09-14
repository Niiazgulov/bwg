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
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	jobCh := make(chan storage.InvoiceJob, 200)
	for i := 0; i < cfg.WorkerCount; i++ {
		go func() {
			for job := range jobCh {
				if err := repo.Invoice(job); err != nil {
					log.Println("Error while invoice", err)
				}
			}
		}()
	}
	jobCh2 := make(chan storage.InvoiceJob, 200)
	for i := 0; i < cfg.WorkerCount; i++ {
		go func() {
			for job := range jobCh2 {
				if err := repo.Withdraw(job); err != nil {
					log.Println("Error while Withdraw", err)
				}
			}
		}()
	}
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetBalanceHandler(repo))
		r.Post("/invoice", handlers.PostInvoiceHandler(repo, jobCh))
		r.Post("/withdraw", handlers.PostWithdrawHandler(repo, jobCh2))
	})
	log.Fatal(http.ListenAndServe(config.Cfg.ServerAddress, r))
}
