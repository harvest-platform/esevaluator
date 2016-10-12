package elastic

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func tHandler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	bodyStr := buf.String()
	body := []byte(bodyStr)
	encoded := EncodeQuery(body)
	b, err := json.Marshal(encoded)
	if err != nil {
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))
}

// TranslateHandler is the default handler for http ES translate requests
func TranslateHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", tHandler)
	return mux
}
