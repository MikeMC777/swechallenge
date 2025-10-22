package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
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

func parseMoney(s string) (*float64, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	clean := moneyRE.ReplaceAllString(s, "")
	var f float64
	_, err := fmt.Sscanf(clean, "%f", &f)
	if err != nil {
		return nil, err
	}
	return &f, nil
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
