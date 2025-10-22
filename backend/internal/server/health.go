package server

import "net/http"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"ok": "true"})
}
