package handlers

import (
	"encoding/json"
	"net/http"

	"bookmanagement/internal/database"
)

// HandleHealth corresponds to GET /health.
// Hosting platforms use this to confirm the app is alive AND its DB is reachable.
// Returns 200 + {"status":"ok"} on success, 503 + {"status":"db unreachable"} otherwise.
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := database.Pool.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "db unreachable"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
