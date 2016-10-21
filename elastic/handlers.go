package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeResponse(w http.ResponseWriter, status int, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
}

func translateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var t Term
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("Could not encode query to map: %v", err),
		})
	}

	e, err := Translate(t)
	if err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("Could not translate query: %v", err),
		})
	}
	writeResponse(w, http.StatusOK, e)
}

// TranslateHandler is the default handler for http ES translate requests
func TranslateHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", translateHandler)
	return mux
}
