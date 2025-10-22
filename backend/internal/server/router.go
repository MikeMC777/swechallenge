// router.go
package server

import "net/http"

func Router(s *Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/api/stocks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.ListStocks(w, r)
			return
		}
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	})
	mux.HandleFunc("/api/stocks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.GetStock(w, r)
			return
		}
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	})
	mux.HandleFunc("/api/recommendations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.Recommend(w, r)
			return
		}
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	})
	return withCORS(mux)
}
