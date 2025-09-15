package main

import (
	"log"
	"wallet-service/internal/app"
	"wallet-service/pkg/config"
)

func main() {
	cfg := config.Load()
	srv := app.NewServer(cfg)
	log.Println("Starting server on port", cfg.Port)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
