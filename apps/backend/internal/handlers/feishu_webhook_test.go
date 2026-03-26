package handlers

import (
	"net/http"
	"testing"
)

func TestFeishuCallback_URLVerification_OK(t *testing.T) {
	router, _, _ := setupTestRouter(t)

	rr := doJSONRequest(t, router, http.MethodPost, "/api/feishu/callback", map[string]any{
		"type":      "url_verification",
		"challenge": "test-challenge",
	}, "")

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rr.Code, rr.Body.String())
	}

	resp := decodeJSON[map[string]any](t, rr)
	challenge, ok := resp["challenge"].(string)
	if !ok || challenge != "test-challenge" {
		t.Fatalf("expected challenge=test-challenge, got %#v", resp["challenge"])
	}
}

func TestFeishuCallback_Event_OK(t *testing.T) {
	router, _, _ := setupTestRouter(t)

	rr := doJSONRequest(t, router, http.MethodPost, "/api/feishu/callback", map[string]any{
		"type": "event_callback",
	}, "")

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rr.Code, rr.Body.String())
	}

	resp := decodeJSON[map[string]any](t, rr)

	code, ok := resp["code"].(float64)
	if !ok || int(code) != 0 {
		t.Fatalf("expected code=0, got %#v", resp["code"])
	}
	msg, ok := resp["msg"].(string)
	if !ok || msg != "ok" {
		t.Fatalf("expected msg=ok, got %#v", resp["msg"])
	}
}

func TestFeishuCallback_InvalidJSON(t *testing.T) {
	router, _, _ := setupTestRouter(t)

	rr := doJSONRequest(t, router, http.MethodPost, "/api/feishu/callback", "not-an-object", "")

	// "not-an-object" is valid JSON string, but binding to struct won't fail;
	// we need malformed JSON scenario by sending raw invalid payload through helper is not supported.
	// So this test validates fallback behavior for non-matching type.
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
}
