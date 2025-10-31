package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

// Verifica que ListStocks responde 200, usa LIMIT/OFFSET correctos cuando no hay filtro
// y serializa el JSON esperado con una fila.
func TestListStocks_OK_NoQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Columnas en el mismo orden del SELECT del handler
	cols := []string{
		"ticker", "company", "last_action", "last_brokerage", "last_rating",
		"last_target", "last_time",
		"raised_30d", "lowered_30d", "reiterated_30d", "initiated_30d",
		"score", "rationale", "updated_at",
	}
	rows := sqlmock.NewRows(cols).
		AddRow(
			"CECO", "CECO Environmental", "target raised by", "Needham & Company LLC", "Buy",
			52.0, "2025-08-22T00:30:05Z",
			1, 0, 0, 0,
			0.92, "demo rationale", "2025-08-22T00:31:00Z",
		)

	// Cuando no hay q, los placeholders quedan: LIMIT $1 OFFSET $2
	mock.ExpectQuery("SELECT ticker, company").
		WithArgs(20, 0).
		WillReturnRows(rows)

	srv := &Server{DB: db}
	req := httptest.NewRequest(http.MethodGet, "/api/stocks?limit=20&offset=0&sort=score&dir=DESC", nil)
	w := httptest.NewRecorder()

	srv.ListStocks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Items  []map[string]any `json:"items"`
		Limit  int              `json:"limit"`
		Offset int              `json:"offset"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v; body=%s", err, w.Body.String())
	}
	if len(body.Items) != 1 {
		t.Fatalf("want 1 item, got %d", len(body.Items))
	}
	if body.Limit != 20 || body.Offset != 0 {
		t.Fatalf("limit/offset mismatch: got (%d,%d)", body.Limit, body.Offset)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

// Verifica que al pasar q se numeran bien los placeholders: WHERE usa $1 y LIMIT/OFFSET $2/$3
func TestListStocks_WithQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	cols := []string{
		"ticker", "company", "last_action", "last_brokerage", "last_rating",
		"last_target", "last_time",
		"raised_30d", "lowered_30d", "reiterated_30d", "initiated_30d",
		"score", "rationale", "updated_at",
	}
	rows := sqlmock.NewRows(cols).
		AddRow(
			"RY", "Royal Bank Of Canada", "target raised by", "Argus", "Buy",
			162.0, "2025-08-31T00:30:05Z",
			1, 0, 0, 0,
			1.10, "demo", "2025-08-31T00:31:00Z",
		)

	// Con q, esperamos 3 args: $1 = %RY%, $2 = limit, $3 = offset
	mock.ExpectQuery("SELECT ticker, company").
		WithArgs("%RY%", 10, 20).
		WillReturnRows(rows)

	srv := &Server{DB: db}
	req := httptest.NewRequest(http.MethodGet, "/api/stocks?q=RY&limit=10&offset=20", nil)
	w := httptest.NewRecorder()

	srv.ListStocks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}
