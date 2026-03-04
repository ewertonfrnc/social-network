package main

import (
	"log"

	"github.com/ewertonfrnc/social-network/internal/env"
)

func main() {
	config := config{
		address: env.GetString("ADDRESS", ":8080"),
	}

	app := &application{
		config,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
