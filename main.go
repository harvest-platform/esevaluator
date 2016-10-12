package main

import (
	"log"
	"net/http"

	"github.com/gerpsh/es-evaluator/elastic"
)

func main() {
	http.Handle("/elastic", elastic.NewHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
