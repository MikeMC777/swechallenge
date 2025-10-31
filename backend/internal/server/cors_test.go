package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_PreflightOPTIONS(t *testing.T) {
	// Dummy handler para probar el wrapper CORS
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	h := withCORS(next)

	req := httptest.NewRequest(http.MethodOptions, "/api/stocks", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", w.Code)
	}
	// Debe incluir cabeceras CORS
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatal("missing Access-Control-Allow-Origin")
	}
}

func TestCORS_AllowsGET(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true}`))
	})
	h := withCORS(next)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatal("missing Access-Control-Allow-Origin")
	}
	if w.Body.Len() == 0 {
		t.Fatal("expected body")
	}
}
