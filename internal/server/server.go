// Package server provides an HTTP endpoint for repository inspection.
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	internal "github.com/richie-rich90454/gitinspect/internal"
)

// InspectRequest is the JSON body for POST /inspect.
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

// Serve starts the HTTP server on the given port.
func Serve(port int) error {
	http.HandleFunc("/inspect", inspectHandler)
	addr := fmt.Sprintf(":%d", port)
	log.Printf("gitinspect server listening on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func inspectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	const maxBodySize = 10 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var req InspectRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, maxBodySize)).Decode(&req); err != nil {
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
