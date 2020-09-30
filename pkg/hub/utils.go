package hub

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/gommon/log"
)

// json http response helper
func httpWriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("could not encode json; got %v", err)
		return
	}
}
