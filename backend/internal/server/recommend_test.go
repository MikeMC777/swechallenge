package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestRecommend_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	cols := []string{"ticker", "company", "score", "rationale"}
	rows := sqlmock.NewRows(cols).
		AddRow("RY", "Royal Bank Of Canada", 1.10, "raised target + recent buys")

	mock.ExpectQuery("SELECT ticker, company, score, rationale FROM tickers_summary").
		WillReturnRows(rows)

	srv := &Server{DB: db}
	req := httptest.NewRequest(http.MethodGet, "/api/recommendations", nil)
	w := httptest.NewRecorder()

	srv.Recommend(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d body=%s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
