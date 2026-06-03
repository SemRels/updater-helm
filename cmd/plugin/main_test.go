// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunUpdatesVersionOnlyByDefault(t *testing.T) {
	t.Parallel()

	file := writeChart(t, "version: 1.0.0\nappVersion: 1.0.0\n")
	env := map[string]string{"SEMREL_VERSION": "v1.1.0", "SEMREL_PLUGIN_FILE": file}

	var stdout, stderr bytes.Buffer
	if code := run(&stdout, &stderr, func(key string) string { return env[key] }); code != 0 {
		t.Fatalf("run() code = %d stderr = %s", code, stderr.String())
	}

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "version: 1.1.0") || !strings.Contains(string(got), "appVersion: 1.0.0") {
		t.Fatalf("updated file = %s", got)
	}
}

func TestRunUpdatesPrefixedAppVersion(t *testing.T) {
	t.Parallel()

	file := writeChart(t, "version: 1.0.0\nappVersion: 1.0.0\n")
	env := map[string]string{
		"SEMREL_VERSION":                   "v1.1.0",
		"SEMREL_PLUGIN_FILE":               file,
		"SEMREL_PLUGIN_UPDATE_APP_VERSION": "true",
		"SEMREL_PLUGIN_APP_VERSION_PREFIX": "v",
	}

	var stdout, stderr bytes.Buffer
	if code := run(&stdout, &stderr, func(key string) string { return env[key] }); code != 0 {
		t.Fatalf("run() code = %d stderr = %s", code, stderr.String())
	}

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "version: 1.1.0") || !strings.Contains(string(got), "appVersion: v1.1.0") {
		t.Fatalf("updated file = %s", got)
	}
}

func TestRunPrefixesExplicitAppVersion(t *testing.T) {
	t.Parallel()

	file := writeChart(t, "version: 1.0.0\nappVersion: 1.0.0\n")
	env := map[string]string{
		"SEMREL_VERSION":                   "1.1.0",
		"SEMREL_PLUGIN_FILE":               file,
		"SEMREL_PLUGIN_APP_VERSION":        "2.0.0",
		"SEMREL_PLUGIN_APP_VERSION_PREFIX": "v",
	}

	var stdout, stderr bytes.Buffer
	if code := run(&stdout, &stderr, func(key string) string { return env[key] }); code != 0 {
		t.Fatalf("run() code = %d stderr = %s", code, stderr.String())
	}

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "version: 1.1.0") || !strings.Contains(string(got), "appVersion: v2.0.0") {
		t.Fatalf("updated file = %s", got)
	}
}

func TestRunDryRun(t *testing.T) {
	t.Parallel()

	env := map[string]string{"SEMREL_VERSION": "1.1.0", "SEMREL_DRY_RUN": "true"}
	var stdout, stderr bytes.Buffer
	if code := run(&stdout, &stderr, func(key string) string { return env[key] }); code != 0 {
		t.Fatalf("run() code = %d", code)
	}
	if !strings.Contains(stdout.String(), "[dry-run]") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunRequiresVersion(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	if code := run(&stdout, &stderr, func(string) string { return "" }); code != 1 {
		t.Fatalf("run() code = %d", code)
	}
	if !strings.Contains(stderr.String(), "SEMREL_VERSION is required") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func writeChart(t *testing.T, content string) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "updater-helm-main-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	file := filepath.Join(dir, "Chart.yaml")
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return file
}
