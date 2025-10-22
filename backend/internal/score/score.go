package score

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
)

type Scorer struct{ DB *sql.DB }

func weightAction(a string) float64 {
	switch {
	case contains(a, "raised"):
		return 1.0
	case contains(a, "lowered"):
		return -1.0
	case contains(a, "initiated"):
		return 0.3
	case contains(a, "reiterated"):
		return 0.2
	default:
		return 0
	}
}

func contains(s, sub string) bool       { return len(s) >= len(sub) && (stringIndexFold(s, sub) >= 0) }
func stringIndexFold(s, sub string) int { return indexFold(s, sub) }
func indexFold(s, sep string) int {
	ls := []rune(strings.ToLower(s))
	lp := []rune(strings.ToLower(sep))
	for i := 0; i+len(lp) <= len(ls); i++ {
		ok := true
		for j := range lp {
			if ls[i+j] != lp[j] {
				ok = false
				break
			}
		}
		if ok {
			return i
		}
	}
	return -1
}

func (s *Scorer) Compute(ctx context.Context) error {
	rows, err := s.DB.QueryContext(ctx, `
      SELECT ticker, action, COALESCE(target_delta,0), EXTRACT(EPOCH FROM (now() - event_time))/86400.0 AS days
      FROM analyst_events
    `)
	if err != nil {
		return err
	}
	type ev struct {
		a    string
		d    float64
		days float64
	}
	m := map[string][]ev{}
	for rows.Next() {
		var t, a string
		var d, days float64
		rows.Scan(&t, &a, &d, &days)
		m[t] = append(m[t], ev{a: a, d: d, days: days})
	}
	for ticker, list := range m {
		var sum float64
		for _, e := range list {
			w := weightAction(e.a)
			decay := math.Pow(0.5, e.days/14.0)
			intensity := math.Tanh(e.d / 5.0)
			sum += w * decay * (1.0 + 0.5*intensity)
		}
		rationale := fmt.Sprintf("sum=%.2f over %d events (w=action, decay=HL14d, intensity=tgtΔ)", sum, len(list))
		if _, err := s.DB.ExecContext(ctx, `UPDATE tickers_summary SET score=$1, rationale=$2, updated_at=now() WHERE ticker=$3`, sum, rationale, ticker); err != nil {
			return err
		}
	}
	return nil
}
