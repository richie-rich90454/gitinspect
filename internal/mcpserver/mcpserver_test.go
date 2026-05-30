package mcpserver

import (
	"reflect"
	"testing"
)

func TestSplitPatternsCommaSeparated(t *testing.T) {
	got := splitPatterns("*.go, *.ts, *.js")
	want := []string{"*.go", "*.ts", "*.js"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("splitPatterns(%q) = %v, want %v", "*.go, *.ts, *.js", got, want)
	}
}

func TestSplitPatternsEmptyStrings(t *testing.T) {
	got := splitPatterns("*.go,,*.ts,")
	want := []string{"*.go", "*.ts"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("splitPatterns(%q) = %v, want %v", "*.go,,*.ts,", got, want)
	}
}

func TestSplitPatternsWhitespace(t *testing.T) {
	got := splitPatterns("  *.go  ,  *.ts  ")
	want := []string{"*.go", "*.ts"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("splitPatterns(%q) = %v, want %v", "  *.go  ,  *.ts  ", got, want)
	}
}

func TestSplitPatternsSinglePattern(t *testing.T) {
	got := splitPatterns("*.go")
	want := []string{"*.go"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("splitPatterns(%q) = %v, want %v", "*.go", got, want)
	}
}

func TestSplitPatternsNoCommas(t *testing.T) {
	got := splitPatterns("**/*.go")
	want := []string{"**/*.go"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("splitPatterns(%q) = %v, want %v", "**/*.go", got, want)
	}
}
