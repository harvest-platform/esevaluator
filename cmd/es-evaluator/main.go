package main

import (
	"log"
	"net/http"

	"github.com/gerpsh/esevaluator/transport"
)

func main() {
	http.Handle("/elastic", transport.TranslateHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
