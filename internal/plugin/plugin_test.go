// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package plugin_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	helm "github.com/SemRels/updater-helm/internal/plugin"
)

func writePluginChartYAML(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, "Chart.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write Chart.yaml: %v", err)
	}
	return p
}

const sampleChart = `apiVersion: v2
name: my-app
description: A sample Helm chart
version: 0.1.0
appVersion: 0.1.0
`

func TestHelmPluginValidateNoChart(t *testing.T) {
	t.Parallel()
	p := helm.NewPlugin(helm.PluginConfig{ChartDir: filepath.Join(t.TempDir(), "missing")})
	if err := p.Validate(); err == nil {
		t.Fatal("Validate() should fail when Chart.yaml does not exist")
	}
}

func TestHelmPluginValidateOK(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writePluginChartYAML(t, dir, sampleChart)
	p := helm.NewPlugin(helm.PluginConfig{ChartDir: dir})
	if err := p.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

func TestHelmPluginExecuteDryRun(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writePluginChartYAML(t, dir, sampleChart)
	p := helm.NewPlugin(helm.PluginConfig{ChartDir: dir})

	result, err := p.Execute(context.Background(), helm.ReleaseContext{Version: "1.2.3", IsDryRun: true})
	if err != nil {
		t.Fatalf("Execute() dry run error: %v", err)
	}
	if result.Outputs["dry_run"] != "true" {
		t.Fatalf("expected dry_run output, got %#v", result.Outputs)
	}
}

func TestHelmPluginExecuteUpdatesVersion(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writePluginChartYAML(t, dir, sampleChart)
	p := helm.NewPlugin(helm.PluginConfig{ChartDir: dir, UpdateAppVersion: true})

	result, err := p.Execute(context.Background(), helm.ReleaseContext{Version: "2.0.0"})
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if result.Outputs["chart_version"] != "2.0.0" {
		t.Fatalf("chart_version = %q", result.Outputs["chart_version"])
	}
	if result.Outputs["app_version"] != "2.0.0" {
		t.Fatalf("app_version = %q", result.Outputs["app_version"])
	}

	meta, err := helm.ReadChartMeta(filepath.Join(dir, "Chart.yaml"))
	if err != nil {
		t.Fatalf("ReadChartMeta error: %v", err)
	}
	if meta.Version != "2.0.0" {
		t.Fatalf("Chart.yaml version = %q", meta.Version)
	}
}
