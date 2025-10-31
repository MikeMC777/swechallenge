package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

// Si el cliente manda sort inválido, el handler cae a 'score'.
// No verificamos el ORDER BY literal, pero sí que la consulta se ejecuta con LIMIT/OFFSET
// y responde 200 (sanity check de la ruta "fallback").
func TestListStocks_SortFallbackToScore(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	cols := []string{"ticker", "company", "last_action", "last_brokerage", "last_rating", "last_target", "last_time", "raised_30d", "lowered_30d", "reiterated_30d", "initiated_30d", "score", "rationale", "updated_at"}
	rows := sqlmock.NewRows(cols).
		AddRow("CECO", "CECO Environmental", "reiterated by", "Canaccord", "Buy", 52.0, "2025-08-25T00:30:04Z", 0, 0, 1, 0, 0.80, "demo", "2025-08-25T00:31:00Z")

	mock.ExpectQuery("SELECT ticker, company").
		WithArgs(5, 0).
		WillReturnRows(rows)

	srv := &Server{DB: db}
	req := httptest.NewRequest(http.MethodGet, "/api/stocks?sort=__h4x__&limit=5&offset=0", nil)
	w := httptest.NewRecorder()

	srv.ListStocks(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d body=%s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
