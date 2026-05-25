package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdaterUpdateChart(t *testing.T) {
	t.Parallel()

	dir, err := os.MkdirTemp("", "updater-helm-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "Chart.yaml")
	original := "apiVersion: v2\nname: demo\nversion: 1.2.3\nappVersion: 1.2.3\n"
	if err := os.WriteFile(file, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := NewUpdater().Update(file, "1.3.0", "1.3.1", false); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "version: 1.3.0") || !strings.Contains(string(got), "appVersion: 1.3.1") {
		t.Fatalf("updated file = %s", got)
	}
}

func TestUpdaterAppVersionOnly(t *testing.T) {
	t.Parallel()

	updated, err := updateContent("version: 1.2.3\n", "1.3.0", "2.0.0", true)
	if err != nil {
		t.Fatalf("updateContent() error = %v", err)
	}
	if !strings.Contains(updated, "version: 1.2.3") || !strings.Contains(updated, "appVersion: 2.0.0") {
		t.Fatalf("updated content = %s", updated)
	}
}

func TestUpdaterMissingFile(t *testing.T) {
	t.Parallel()

	err := NewUpdater().Update(filepath.Join(t.TempDir(), "Chart.yaml"), "1.3.0", "", false)
	if err == nil || !strings.Contains(err.Error(), "read") {
		t.Fatalf("expected read error, got %v", err)
	}
}

func TestUpdaterMissingVersion(t *testing.T) {
	t.Parallel()

	_, err := updateContent("name: demo\n", "1.3.0", "", false)
	if err == nil || !strings.Contains(err.Error(), "version field not found") {
		t.Fatalf("expected version error, got %v", err)
	}
}
