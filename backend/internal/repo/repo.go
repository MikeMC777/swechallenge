package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Repo struct{ DB *sql.DB }

func New(db *sql.DB) *Repo { return &Repo{DB: db} }

func (r *Repo) InsertRaw(ctx context.Context, source, pageKey string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, `INSERT INTO stocks_raw (source, page_key, payload) VALUES ($1,$2,$3)`, source, pageKey, b)
	return err
}

var moneyRE = regexp.MustCompile(`[$,]`)

// parseMoney convierte strings como "$1,234.50", "+$1000.00", "-$12.34" → *float64.
// Reglas:
// - Si el string queda sin dígitos luego de limpiar ("" , " ", "$", "+$", "-$"), retorna (nil, nil).
// - Si hay comas, deben estar bien formadas en miles: d{1,3}(,d{3})* (opcional .decimales).
// - Si no hay comas, acepta d+(.decimales)?
// - Acepta signo + o - y símbolo $ en cualquier orden al inicio.
func parseMoney(s string) (*float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	// Regex de validación (antes de normalizar), para detectar comas mal ubicadas.
	// Permitimos dos formas:
	//  A) Con comas de miles correctas:  ^[+-]?\$?\d{1,3}(?:,\d{3})*(?:\.\d+)?$
	//  B) Sin comas:                    ^[+-]?\$?\d+(?:\.\d+)?$
	reThousands := regexp.MustCompile(`^[+-]?\$?\d{1,3}(?:,\d{3})*(?:\.\d+)?$`)
	rePlain := regexp.MustCompile(`^[+-]?\$?\d+(?:\.\d+)?$`)

	// Si contiene coma, debe cumplir la variante con miles; si no, la simple.
	if strings.Contains(s, ",") {
		if !reThousands.MatchString(s) {
			return nil, fmt.Errorf("parseMoney: invalid thousands format")
		}
	} else {
		if !rePlain.MatchString(s) {
			// Puede que sea sólo "$" o "+$" etc. -> tratar como sin dato
			onlySymbols := strings.Trim(s, "+-$ $")
			if onlySymbols == "" {
				return nil, nil
			}
			return nil, fmt.Errorf("parseMoney: invalid format")
		}
	}

	// Extraer signo manualmente (para soportar +$ y -$)
	sign := 1.0
	if strings.HasPrefix(s, "+") {
		s = s[1:]
	} else if strings.HasPrefix(s, "-") {
		sign = -1
		s = s[1:]
	}

	// Quitar símbolo $ si está
	s = strings.TrimPrefix(s, "$")
	// Quitar comas
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)

	// Si después de limpiar no quedan dígitos -> sin dato
	if s == "" {
		return nil, nil
	}

	// Parse final
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("parseMoney: %w", err)
	}
	v *= sign
	return &v, nil
}

func (r *Repo) UpsertAnalystEvent(ctx context.Context, m map[string]any) error {
	ticker := str(m["ticker"])
	if ticker == "" {
		return fmt.Errorf("missing ticker")
	}
	company := optStr(m["company"])
	action := optStr(m["action"])
	brokerage := optStr(m["brokerage"])
	ratingFrom := optStr(m["rating_from"])
	ratingTo := optStr(m["rating_to"])
	tFrom, _ := parseMoney(str(m["target_from"]))
	tTo, _ := parseMoney(str(m["target_to"]))
	var delta *float64
	if tFrom != nil && tTo != nil {
		d := *tTo - *tFrom
		delta = &d
	}
	eventTime := str(m["time"])

	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO analyst_events (ticker, company, action, brokerage, rating_from, rating_to, target_from, target_to, target_delta, event_time)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
    `, ticker, company, action, brokerage, ratingFrom, ratingTo, tFrom, tTo, delta, eventTime)
	if err != nil {
		return err
	}

	_, err = r.DB.ExecContext(ctx, `
		WITH last AS (
			SELECT * FROM analyst_events WHERE ticker=$1 ORDER BY event_time DESC LIMIT 1
		),
		win AS (
			SELECT
			SUM(CASE WHEN action ILIKE '%raised%'     THEN 1 ELSE 0 END) AS raised,
			SUM(CASE WHEN action ILIKE '%lowered%'    THEN 1 ELSE 0 END) AS lowered,
			SUM(CASE WHEN action ILIKE '%reiterated%' THEN 1 ELSE 0 END) AS reiterated,
			SUM(CASE WHEN action ILIKE '%initiated%'  THEN 1 ELSE 0 END) AS initiated
			FROM analyst_events
			WHERE ticker=$1 AND event_time >= now() - INTERVAL '30 days'
		)
		INSERT INTO tickers_summary (
			ticker, company, last_action, last_brokerage, last_rating, last_target, last_time,
			raised_30d, lowered_30d, reiterated_30d, initiated_30d, score, rationale, updated_at
		)
		SELECT
			$1,
			last.company,
			last.action,
			last.brokerage,
			last.rating_to,
			last.target_to,
			last.event_time,
			COALESCE(win.raised, 0),
			COALESCE(win.lowered, 0),
			COALESCE(win.reiterated, 0),
			COALESCE(win.initiated, 0),
			0, '', now()
		FROM last, win
		ON CONFLICT (ticker) DO UPDATE SET
			company = EXCLUDED.company,
			last_action = EXCLUDED.last_action,
			last_brokerage = EXCLUDED.last_brokerage,
			last_rating = EXCLUDED.last_rating,
			last_target = EXCLUDED.last_target,
			last_time = EXCLUDED.last_time,
			raised_30d = EXCLUDED.raised_30d,
			lowered_30d = EXCLUDED.lowered_30d,
			reiterated_30d = EXCLUDED.reiterated_30d,
			initiated_30d = EXCLUDED.initiated_30d,
			updated_at = now();
	`, ticker)

	return err
}

func str(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return strings.Trim(string(b), "\"")
}
func optStr(v any) *string {
	s := str(v)
	if s == "" {
		return nil
	}
	return &s
}
