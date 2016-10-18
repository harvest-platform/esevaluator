package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func translateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var t Term
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err,
		})
		return
	}
	e, err := Translate(t)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query translation returned with error: %v", err), http.StatusUnprocessableEntity)
	}
	if err := json.NewEncoder(w).Encode(e); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err,
		})
		return
	}
}

// TranslateHandler is the default handler for http ES translate requests
func TranslateHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", translateHandler)
	return mux
}
