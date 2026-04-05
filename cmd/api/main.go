package main

import (
	"github.com/ewertonfrnc/social-network/internal/db"
	"github.com/ewertonfrnc/social-network/internal/env"
	"github.com/ewertonfrnc/social-network/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

func main() {
	config := config{
		address: env.GetString("ADDRESS", ":8080"),
		db: dbConfig{
			address:      env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	// Database
	db, err := db.New(
		config.db.address,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	store := store.NewDBStorage(db)

	app := &application{
		config,
		store,
		logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
