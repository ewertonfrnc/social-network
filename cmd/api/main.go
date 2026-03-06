package main

import (
	"log"

	"github.com/ewertonfrnc/social-network/internal/db"
	"github.com/ewertonfrnc/social-network/internal/env"
	"github.com/ewertonfrnc/social-network/internal/store"
)

func main() {
	config := config{
		address: env.GetString("ADDRESS", ":8080"),
		db: dbConfig{
			address:      env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		config.db.address,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("Database connection pool established")

	store := store.NewDBStorage(db)

	app := &application{
		config,
		store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
