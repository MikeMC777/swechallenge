package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestRecomputeScores_POST_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("UPDATE tickers_summary SET score=score, updated_at=now\\(\\)").
		WillReturnResult(sqlmock.NewResult(0, 1))

	srv := &Server{DB: db}
	req := httptest.NewRequest(http.MethodPost, "/api/score/recompute", nil)
	w := httptest.NewRecorder()

	srv.RecomputeScores(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d body=%s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRecomputeScores_MethodNotAllowed(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	srv := &Server{DB: db}
	req := httptest.NewRequest(http.MethodGet, "/api/score/recompute", nil)
	w := httptest.NewRecorder()

	srv.RecomputeScores(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405, got %d body=%s", w.Code, w.Body.String())
	}
}
