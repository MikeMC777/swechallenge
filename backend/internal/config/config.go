package config

import (
	"log"
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	SourceAPIURL   string
	SourceAPIToken string
}

func Load() Config {
	return Config{
		Port:           getenv("PORT", "8081"),
		DatabaseURL:    must("DATABASE_URL"),
		SourceAPIURL:   getenv("SOURCE_API_URL", "https://api.karenai.click/swechallenge/list"),
		SourceAPIToken: must("SOURCE_API_TOKEN"),
	}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}
