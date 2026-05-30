package internal

import (
	"testing"
)

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	original := &inspectResult{
		tree: map[string]string{
			"main.go": "package main\nfunc main() {}",
			"foo.go":  "package foo\n",
		},
		paths:        []string{"main.go", "foo.go"},
		totalBytes:   42,
		totalTokens:  11,
		truncated:    true,
		dependencies: []string{"fmt", "os"},
	}

	data, err := marshalCache(original)
	if err != nil {
		t.Fatalf("marshalCache error: %v", err)
	}

	var restored inspectResult
	if err := unmarshalCache(data, &restored); err != nil {
		t.Fatalf("unmarshalCache error: %v", err)
	}

	if len(restored.tree) != len(original.tree) {
		t.Errorf("tree length: got %d, want %d", len(restored.tree), len(original.tree))
	}
	for k, v := range original.tree {
		if restored.tree[k] != v {
			t.Errorf("tree[%q]: got %q, want %q", k, restored.tree[k], v)
		}
	}
	if len(restored.paths) != len(original.paths) {
		t.Errorf("paths length: got %d, want %d", len(restored.paths), len(original.paths))
	}
	for i, p := range original.paths {
		if restored.paths[i] != p {
			t.Errorf("paths[%d]: got %q, want %q", i, restored.paths[i], p)
		}
	}
	if restored.totalBytes != original.totalBytes {
		t.Errorf("totalBytes: got %d, want %d", restored.totalBytes, original.totalBytes)
	}
	if restored.totalTokens != original.totalTokens {
		t.Errorf("totalTokens: got %d, want %d", restored.totalTokens, original.totalTokens)
	}
	if restored.truncated != original.truncated {
		t.Errorf("truncated: got %v, want %v", restored.truncated, original.truncated)
	}
	if len(restored.dependencies) != len(original.dependencies) {
		t.Errorf("dependencies length: got %d, want %d", len(restored.dependencies), len(original.dependencies))
	}
	for i, d := range original.dependencies {
		if restored.dependencies[i] != d {
			t.Errorf("dependencies[%d]: got %q, want %q", i, restored.dependencies[i], d)
		}
	}
}

func TestUnmarshalInvalidData(t *testing.T) {
	var r inspectResult
	if err := unmarshalCache([]byte("not valid json{{{"), &r); err == nil {
		t.Error("expected error for invalid JSON data, got nil")
	}
}
