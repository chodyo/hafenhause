package main

import (
	"log"
	"net/http"

	"github.com/chodyo/hafenhause"
)

func main() {
	http.HandleFunc("/Bedtime", hafenhause.Bedtime)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
