package main

import (
	"log"

	"github.com/ewertonfrnc/social-network/internal/db"
	"github.com/ewertonfrnc/social-network/internal/env"
	"github.com/ewertonfrnc/social-network/internal/store"
)

func main() {
	address := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable")
	conn, err := db.New(address, 10, 10, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewDBStorage(conn)
	db.Seed(store, conn)
}
