package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/harvest-platform/esevaluator/transport"
)

func main() {
	log.SetFlags(0)

	var (
		httpaddr string
		tlscert  string
		tlskey   string
	)

	flag.Parse()

	flag.StringVar(&httpaddr, "http", "127.0.0.1:8080", "Address for HTTP transport.")
	flag.StringVar(&tlscert, "tlscert", "", "Path to TLS certificate.")
	flag.StringVar(&tlskey, "tlskey", "", "Path to TLS key.")

	http.Handle("/", transport.PingHandler())
	http.Handle("/elastic", transport.TranslateHandler())

	fmt.Printf("Listening on %s\n", httpaddr)
	if tlscert == "" {
		log.Fatal(http.ListenAndServe(httpaddr, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(httpaddr, tlscert, tlskey, nil))
	}
}
