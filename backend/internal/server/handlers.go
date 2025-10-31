package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Server struct{ DB *sql.DB }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// GET /api/stocks
// Lista de tickers con resumen (para la tabla). Soporta q, sort, dir, limit, offset
func (s *Server) ListStocks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	// Whitelist de columnas válidas para ORDER BY
	allowed := map[string]bool{
		"ticker": true, "company": true,
		"last_action": true, "last_brokerage": true, "last_rating": true,
		"last_target": true, "last_time": true,
		"raised_30d": true, "lowered_30d": true, "reiterated_30d": true, "initiated_30d": true,
		"score": true, "updated_at": true,
	}

	sort := r.URL.Query().Get("sort")
	if !allowed[sort] {
		sort = "score"
	}

	dir := strings.ToUpper(r.URL.Query().Get("dir"))
	if dir != "ASC" {
		dir = "DESC"
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	base := `SELECT ticker, company, last_action, last_brokerage, last_rating, last_target, last_time,
                    raised_30d, lowered_30d, reiterated_30d, initiated_30d, score, rationale, updated_at
             FROM tickers_summary`

	args := []any{}
	where := ""
	if q != "" {
		where = " WHERE ticker ILIKE $1 OR company ILIKE $1 OR last_brokerage ILIKE $1"
		args = append(args, "%"+q+"%")
	}

	// Numerar LIMIT/OFFSET según cantidad de args acumulados
	limIdx := len(args) + 1
	offIdx := len(args) + 2

	sqlStr := base + where +
		" ORDER BY " + sort + " " + dir +
		fmt.Sprintf(" LIMIT $%d OFFSET $%d", limIdx, offIdx)

	args = append(args, limit, offset)

	rows, err := s.DB.Query(sqlStr, args...)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	type item struct {
		Ticker        string   `json:"ticker"`
		Company       *string  `json:"company,omitempty"`
		LastAction    *string  `json:"last_action,omitempty"`
		LastBrokerage *string  `json:"last_brokerage,omitempty"`
		LastRating    *string  `json:"last_rating,omitempty"`
		LastTarget    *float64 `json:"last_target,omitempty"`
		LastTime      string   `json:"last_time"`
		Raised30      int      `json:"raised_30d"`
		Lowered30     int      `json:"lowered_30d"`
		Reiterated30  int      `json:"reiterated_30d"`
		Initiated30   int      `json:"initiated_30d"`
		Score         float64  `json:"score"`
		Rationale     string   `json:"rationale"`
		UpdatedAt     string   `json:"updated_at"`
	}

	var out []item
	for rows.Next() {
		var it item
		if err := rows.Scan(
			&it.Ticker, &it.Company, &it.LastAction, &it.LastBrokerage, &it.LastRating,
			&it.LastTarget, &it.LastTime, &it.Raised30, &it.Lowered30, &it.Reiterated30, &it.Initiated30,
			&it.Score, &it.Rationale, &it.UpdatedAt,
		); err == nil {
			out = append(out, it)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": out, "limit": limit, "offset": offset})
}

// GET /api/stocks/{symbol}
// Devuelve timeline de eventos del ticker
func (s *Server) GetStock(w http.ResponseWriter, r *http.Request) {
	sym := strings.TrimPrefix(r.URL.Path, "/api/stocks/")
	rows, err := s.DB.Query("SELECT company, action, brokerage, rating_from, rating_to, target_from, target_to, target_delta, event_time FROM analyst_events WHERE ticker=$1 ORDER BY event_time DESC", sym)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()
	type ev struct {
		Company     *string  `json:"company,omitempty"`
		Action      *string  `json:"action,omitempty"`
		Brokerage   *string  `json:"brokerage,omitempty"`
		RatingFrom  *string  `json:"rating_from,omitempty"`
		RatingTo    *string  `json:"rating_to,omitempty"`
		TargetFrom  *float64 `json:"target_from,omitempty"`
		TargetTo    *float64 `json:"target_to,omitempty"`
		TargetDelta *float64 `json:"target_delta,omitempty"`
		Time        string   `json:"time"`
	}
	var list []ev
	for rows.Next() {
		var e ev
		rows.Scan(&e.Company, &e.Action, &e.Brokerage, &e.RatingFrom, &e.RatingTo, &e.TargetFrom, &e.TargetTo, &e.TargetDelta, &e.Time)
		list = append(list, e)
	}
	if len(list) == 0 {
		writeJSON(w, 404, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, 200, map[string]any{"ticker": sym, "events": list})
}

// GET /api/recommendations
// Devuelve top-N por score (N fijo por simplicidad)
func (s *Server) Recommend(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.Query(`SELECT ticker, company, score, rationale FROM tickers_summary ORDER BY score DESC LIMIT 5`)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()
	type out struct {
		Ticker    string  `json:"symbol"`
		Company   *string `json:"name,omitempty"`
		Score     float64 `json:"score"`
		Rationale string  `json:"rationale"`
	}
	var list []out
	for rows.Next() {
		var o out
		rows.Scan(&o.Ticker, &o.Company, &o.Score, &o.Rationale)
		list = append(list, o)
	}
	writeJSON(w, 200, map[string]any{"items": list})
}

// POST /api/score/recompute
// Dispara recomputo de score (stub: solo toca updated_at). Se integrará con score.Scorer.
func (s *Server) RecomputeScores(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if _, err := s.DB.Exec(`UPDATE tickers_summary SET score=score, updated_at=now()`); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
