package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func newServer(apiKey string) *Server {
	return &Server{
		apiKey: apiKey,
		sem:    make(chan struct{}, 4),
	}
}

func postJSON(s *Server, body any, apiKey string) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/inspect", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	rec := httptest.NewRecorder()
	s.inspectHandler(rec, req)
	return rec
}

func TestMethodNotAllowed(t *testing.T) {
	s := newServer("")
	req := httptest.NewRequest(http.MethodGet, "/inspect", nil)
	rec := httptest.NewRecorder()
	s.inspectHandler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET /inspect: got %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestMissingRepoField(t *testing.T) {
	s := newServer("")
	rec := postJSON(s, map[string]any{}, "")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("missing repo: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestInvalidFormat(t *testing.T) {
	s := newServer("")
	rec := postJSON(s, InspectRequest{Repo: "/tmp", Format: "xml"}, "")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("invalid format: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestMaxTokensOutOfRange(t *testing.T) {
	s := newServer("")
	rec := postJSON(s, InspectRequest{Repo: "/tmp", MaxTokens: -1}, "")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("negative max_tokens: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
	rec = postJSON(s, InspectRequest{Repo: "/tmp", MaxTokens: 200000}, "")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("too large max_tokens: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestMaxFilesOutOfRange(t *testing.T) {
	s := newServer("")
	rec := postJSON(s, InspectRequest{Repo: "/tmp", MaxFiles: -1}, "")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("negative max_files: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
	rec = postJSON(s, InspectRequest{Repo: "/tmp", MaxFiles: 200000}, "")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("too large max_files: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func sampleRepoPath(t *testing.T) string {
	t.Helper()
	repoPath := filepath.Join("..", "..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}
	return absPath
}

func TestAuthValidAPIKey(t *testing.T) {
	s := newServer("secret123")
	absPath := sampleRepoPath(t)
	rec := postJSON(s, InspectRequest{Repo: absPath, MaxTokens: 6000}, "secret123")
	if rec.Code != http.StatusOK {
		t.Errorf("valid API key: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAuthInvalidAPIKey(t *testing.T) {
	s := newServer("secret123")
	rec := postJSON(s, InspectRequest{Repo: "/tmp"}, "wrongkey")
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("invalid API key: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthMissingAPIKey(t *testing.T) {
	s := newServer("secret123")
	rec := postJSON(s, InspectRequest{Repo: "/tmp"}, "")
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("missing API key: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestNoAuthWhenNoAPIKeyConfigured(t *testing.T) {
	s := newServer("")
	absPath := sampleRepoPath(t)
	rec := postJSON(s, InspectRequest{Repo: absPath, MaxTokens: 6000}, "")
	if rec.Code != http.StatusOK {
		t.Errorf("no auth required: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestValidLocalRepoInspection(t *testing.T) {
	s := newServer("")
	absPath := sampleRepoPath(t)
	rec := postJSON(s, InspectRequest{Repo: absPath, MaxTokens: 6000}, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("valid local repo: got %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected application/json content type, got %s", contentType)
	}
}
