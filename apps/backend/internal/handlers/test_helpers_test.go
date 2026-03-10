package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gms/backend/internal/auth"
	"gms/backend/internal/config"
	"gms/backend/internal/models"
)

func setupTestRouter(t *testing.T) (http.Handler, *gorm.DB, config.Config) {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Project{}, &models.Task{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	hash, err := auth.HashPassword("admin123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	seed := models.User{Name: "Owner", Email: "admin@gms.local", PasswordHash: hash, Role: models.RoleOwner}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	cfg := config.Config{JWTSecret: "test-secret", JWTTTLHours: 24}
	engine := SetupRouter(db, cfg)
	return engine, db, cfg
}

func doJSONRequest(t *testing.T, h http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()

	var payload []byte
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		payload = data
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func decodeJSON[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode response: %v\nbody=%s", err, rr.Body.String())
	}
	return out
}
