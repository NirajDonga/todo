package main

import (
	"context"
	"log"

	"github.com/NirajDonga/todo/internal/config"
	"github.com/NirajDonga/todo/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	ctx := context.Background()
	dbPool, err := database.Connect(ctx, cfg.DB_URI)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	defer dbPool.Close()

}
