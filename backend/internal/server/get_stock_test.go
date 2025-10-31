package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

// 200 OK con al menos un evento
func TestGetStock_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	cols := []string{"company", "action", "brokerage", "rating_from", "rating_to", "target_from", "target_to", "target_delta", "event_time"}
	rows := sqlmock.NewRows(cols).
		AddRow("CECO Environmental", "target raised by", "Needham & Company LLC", "Buy", "Buy", 44.0, 52.0, 8.0, "2025-08-22T00:30:05Z")

	mock.ExpectQuery("SELECT company, action").
		WithArgs("CECO").
		WillReturnRows(rows)

	srv := &Server{DB: db}
	req := httptest.NewRequest(http.MethodGet, "/api/stocks/CECO", nil)
	w := httptest.NewRecorder()

	srv.GetStock(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d body=%s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// 404 cuando no hay eventos para el ticker
func TestGetStock_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT company, action").
		WithArgs("XYZ").
		WillReturnRows(sqlmock.NewRows([]string{"company", "action", "brokerage", "rating_from", "rating_to", "target_from", "target_to", "target_delta", "event_time"}))

	srv := &Server{DB: db}
	req := httptest.NewRequest(http.MethodGet, "/api/stocks/XYZ", nil)
	w := httptest.NewRecorder()

	srv.GetStock(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d body=%s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
