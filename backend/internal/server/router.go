package server

import (
	"net/http"
	"strings"
)

// Aquí se registran TODOS los endpoints públicos de la API.
//
// - GET /health → quick healthcheck JSON
// - GET /api/stocks → lista (q, sort, dir, limit, offset)
// - GET /api/stocks/{symbol} → detalle/timeline del ticker
// - GET /api/recommendations → top N por score
// - POST /api/score/recompute → recomputa scores (stub por ahora)
//
// CORS se aplica con withCORS.
func Router(s *Server) http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"ok": "true"})
	})

	// Lista de tickers (tabla)
	mux.HandleFunc("/api/stocks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.ListStocks(w, r)
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})

	// Detalle de ticker (timeline). Coincide /api/stocks/{symbol}
	mux.HandleFunc("/api/stocks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Validar que hay un símbolo después del prefijo
			if sym := strings.TrimPrefix(r.URL.Path, "/api/stocks/"); sym == "" || strings.Contains(sym, "/") {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
				return
			}
			s.GetStock(w, r)
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})

	// Recomendaciones
	mux.HandleFunc("/api/recommendations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.Recommend(w, r)
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})

	// Recompute scores (POST)
	mux.HandleFunc("/api/score/recompute", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.RecomputeScores(w, r)
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})

	return withCORS(mux)
}
