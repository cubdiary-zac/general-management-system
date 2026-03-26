package handlers

import (
	"net/http"
	"testing"
)

func TestHealth_ReturnsOKStatus(t *testing.T) {
	router, _, _ := setupTestRouter(t)

	rr := doJSONRequest(t, router, http.MethodGet, "/api/health", nil, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}

	resp := decodeJSON[map[string]any](t, rr)
	if resp["status"] != "ok" {
		t.Fatalf("expected status=ok, got %v", resp["status"])
	}
}
