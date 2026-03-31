package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"email-dispatcher/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	DB, err = sql.Open("postgres", config.GetEnv("DB_URL"))
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = DB.PingContext(ctx)
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	log.Println("✅ Connected to PostgreSQL")
}
