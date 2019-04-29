package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/chodyo/hafenhause"
	"github.com/joho/godotenv"
)

type config struct {
	PublicPort string `env:"PUBLIC_PORT" envDefault:":8080"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, using defaults")
	}

	var config config
	if err := env.Parse(&config); err != nil {
		log.Fatalln("failed to parse ENV")
	}

	log.Printf("Starting server on port %s\n", config.PublicPort)

	log.Fatal(http.ListenAndServe(config.PublicPort, handler()))
}

func handler() http.Handler {
	return http.HandlerFunc(hafenhause.Bedtime)
}
