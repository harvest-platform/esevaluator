package elastic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func translateHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	var t Term
	json.Unmarshal(b, &t)
	e, err := Translate(t)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query translation returned with error: %v", err), 422)
	}
	q, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(q))
}

// TranslateHandler is the default handler for http ES translate requests
func TranslateHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", translateHandler)
	return mux
}
