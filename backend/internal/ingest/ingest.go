package ingest

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MikeMC777/swechallenge/internal/repo"
)

type Service struct {
	Client *Client
	Repo   *repo.Repo
}

func (s *Service) Run(ctx context.Context) (int, error) {
	next := ""
	total := 0
	for {
		page, err := s.Client.Fetch(next)
		if err != nil {
			return total, err
		}
		if page.Items == nil {
			break
		}

		if err := s.Repo.InsertRaw(ctx, "karenai", next, page.Items); err != nil {
			return total, err
		}

		switch t := page.Items.(type) {
		case []any:
			for _, row := range t {
				b, _ := json.Marshal(row)
				var m map[string]any
				json.Unmarshal(b, &m)
				if err := s.Repo.UpsertAnalystEvent(ctx, m); err != nil {
					log.Printf("skip row: %v", err)
				} else {
					total++
				}
			}
		case map[string]any:
			if err := s.Repo.UpsertAnalystEvent(ctx, t); err == nil {
				total++
			}
		}

		if page.Next == "" {
			break
		}
		next = page.Next
		time.Sleep(150 * time.Millisecond)
	}
	return total, nil
}
