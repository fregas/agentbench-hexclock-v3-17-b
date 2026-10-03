package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var fixed = time.Date(2026, 10, 3, 12, 34, 56, 0, time.FixedZone("EST", -5*3600))

func serve(t *testing.T, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	newHandler(func() time.Time { return fixed }).ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func TestRootReturnsUTCTime(t *testing.T) {
	rec := serve(t, http.MethodGet, "/")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var got timeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	want := timeResponse{UTC: "2026-10-03T17:34:56Z", Unix: fixed.Unix(), Hex: "0x6ac13cc0"}
	if got != want {
		t.Errorf("body = %+v, want %+v", got, want)
	}
}

func TestHealthz(t *testing.T) {
	rec := serve(t, http.MethodGet, "/healthz")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got["status"] != "ok" {
		t.Errorf("status field = %q, want ok", got["status"])
	}
}

func TestUnknownPathIs404(t *testing.T) {
	if rec := serve(t, http.MethodGet, "/nope"); rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestWrongMethodIs405(t *testing.T) {
	if rec := serve(t, http.MethodPost, "/"); rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}
