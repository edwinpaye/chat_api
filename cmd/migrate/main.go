// migrate applies migrations/*.sql in order. Usage: go run ./cmd/migrate
package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"chat_api/internal/config"
	"chat_api/internal/platform"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	db, err := platform.OpenDB(cfg.Database.URL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(files)
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("%s: %v", f, err)
		}
		// Split on statement boundaries (naive, fine for generated DDL).
		for _, stmt := range strings.Split(string(raw), ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "--") && !strings.Contains(stmt, "\n") {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				// idempotency: ignore "already exists" errors
				if strings.Contains(err.Error(), "already exists") {
					continue
				}
				log.Fatalf("%s: %v", f, err)
			}
		}
		log.Printf("applied %s", f)
	}
	log.Println("migrations complete")
}