package app

import (
	"fmt"
	"log"
	"net/http"
	"wallet-service/internal/api"
	"wallet-service/internal/repo"
	"wallet-service/internal/service"
	"wallet-service/pkg/config"
	"wallet-service/pkg/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Server struct {
	cfg    *config.Config
	router *http.ServeMux
}

func NewServer(cfg *config.Config) *Server {
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	repo := repo.NewWalletRepo(db)
	svc := service.NewWalletService(repo)
	handler := api.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/wallets/", handler.GetBalance)
	mux.HandleFunc("/api/v1/wallet", handler.HandleOperation)

	return &Server{
		cfg:    cfg,
		router: mux,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.cfg.Port)
	logger.Info("server listening on", addr)
	return http.ListenAndServe(addr, s.router)
}
