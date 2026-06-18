package main

import (
	"log"
	"net/http"

	"github.com/ijry/sanjie/internal/app"
	"github.com/ijry/sanjie/internal/config"
	"github.com/ijry/sanjie/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Database.Driver != "sqlite" {
		log.Fatalf("unsupported database driver: %s", cfg.Database.Driver)
	}

	database, err := db.Open(cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}
	if cfg.Demo.Seed {
		if err := db.Seed(database); err != nil {
			log.Fatal(err)
		}
	}

	application := app.New(database)
	log.Printf("sanjie server listening on %s", cfg.Server.Addr)
	log.Fatal(http.ListenAndServe(cfg.Server.Addr, application.Router))
}
