package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/MikeMC777/swechallenge/internal/config"
	"github.com/MikeMC777/swechallenge/internal/db"
	"github.com/MikeMC777/swechallenge/internal/ingest"
	"github.com/MikeMC777/swechallenge/internal/repo"
	"github.com/MikeMC777/swechallenge/internal/server"
)

func main() {
	cfg := config.Load()

	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	// Ejecuta migraciones desde el archivo incluido en la imagen
	if b, err := os.ReadFile("internal/db/migrate.sql"); err == nil {
		if _, err := database.Exec(string(b)); err != nil {
			log.Fatal(err)
		}
	}

	// Ingesta on-boot (idempotente)
	cl := ingest.New(cfg.SourceAPIURL, cfg.SourceAPIToken)
	r := repo.New(database)
	svc := &ingest.Service{Client: cl, Repo: r}
	if n, err := svc.Run(context.Background()); err != nil {
		log.Printf("ingest error: %v", err)
	} else {
		log.Printf("ingested %d items", n)
	}

	api := &server.Server{DB: database}
	log.Printf("listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, server.Router(api)))
}
