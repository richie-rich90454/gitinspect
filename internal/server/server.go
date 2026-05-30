// Package server provides an HTTP API for running repository inspections.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	internal "github.com/richie-rich90454/gitinspect/internal"
	"github.com/richie-rich90454/gitinspect/internal/repo"
)

const maxBodySize = 10 << 20

// InspectRequest represents the JSON body of an /inspect request.
type InspectRequest struct {
	Repo      string   `json:"repo"`
	Format    string   `json:"format,omitempty"`
	MaxTokens int      `json:"max_tokens,omitempty"`
	MaxFiles  int      `json:"max_files,omitempty"`
	Include   []string `json:"include,omitempty"`
	Exclude   []string `json:"exclude,omitempty"`
	Strip     bool     `json:"strip,omitempty"`
	NoCache   bool     `json:"no_cache,omitempty"`
}

// Server is an HTTP server that handles repository inspection requests.
type Server struct {
	apiKey string
	sem    chan struct{}
}

// Serve starts the inspection HTTP server on the given port.
func Serve(port int, apiKey string, maxConcurrent int) error {
	if maxConcurrent <= 0 {
		maxConcurrent = 4
	}

	s := &Server{
		apiKey: apiKey,
		sem:    make(chan struct{}, maxConcurrent),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/inspect", s.inspectHandler)

	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("gitinspect server listening on %s", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Printf("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
		return err
	}

	cleanupTempDirs()
	log.Printf("server stopped gracefully")
	return nil
}

func cleanupTempDirs() {
	tmpDir := os.TempDir()
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "gitinspect-") {
			_ = os.RemoveAll(filepath.Join(tmpDir, entry.Name()))
		}
	}
}

func (s *Server) inspectHandler(w http.ResponseWriter, r *http.Request) {
	if s.apiKey != "" {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != s.apiKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var req InspectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Repo == "" {
		http.Error(w, "repo field is required", http.StatusBadRequest)
		return
	}

	if req.Format == "" {
		req.Format = "json"
	}

	validFormats := map[string]bool{"json": true, "text": true, "yaml": true}
	if !validFormats[req.Format] {
		http.Error(w, "format must be one of: json, text, yaml", http.StatusBadRequest)
		return
	}

	if req.MaxTokens != 0 && (req.MaxTokens < 1 || req.MaxTokens > 100000) {
		http.Error(w, "max_tokens must be between 1 and 100000", http.StatusBadRequest)
		return
	}

	if req.MaxFiles < 0 || req.MaxFiles > 100000 {
		http.Error(w, "max_files must be between 0 and 100000", http.StatusBadRequest)
		return
	}

	info, statErr := os.Stat(req.Repo)
	isLocal := statErr == nil && info.IsDir()

	if isLocal {
		if err := repo.ValidateLocalPath(req.Repo, true); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if err := repo.ValidateRemoteURL(req.Repo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if !isLocal {
		select {
		case s.sem <- struct{}{}:
			defer func() { <-s.sem }()
		default:
			http.Error(w, "too many concurrent requests", http.StatusServiceUnavailable)
			return
		}
	}

	if req.MaxTokens == 0 {
		req.MaxTokens = 6000
	}

	opts := internal.Options{
		Format:    req.Format,
		MaxTokens: req.MaxTokens,
		MaxFiles:  req.MaxFiles,
		Include:   req.Include,
		Exclude:   req.Exclude,
		Strip:     req.Strip,
		NoCache:   req.NoCache,
	}

	result, err := internal.RunInspect(req.Repo, opts)
	if err != nil {
		http.Error(w, "inspection failed", http.StatusInternalServerError)
		return
	}

	contentType := "application/json"
	switch req.Format {
	case "text":
		contentType = "text/plain"
	case "yaml":
		contentType = "application/yaml"
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)
}
