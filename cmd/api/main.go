package main

import (
	"context"
	"log"

	"github.com/NirajDonga/todo/internal/config"
	"github.com/NirajDonga/todo/internal/database"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	dbPool, err := database.Connect(ctx, cfg.DB_URI)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()


	

}
