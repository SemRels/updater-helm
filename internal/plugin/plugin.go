// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

// Package plugin provides Helm chart versioning and publishing helpers.
package plugin

import (
	"context"
	"fmt"
	"path/filepath"
)

// Result is the outcome of running the standalone Helm plugin.
type Result struct {
	Outputs    map[string]string
	Skipped    bool
	SkipReason string
}

// ReleaseContext contains the release information needed by Execute.
type ReleaseContext struct {
	Version  string
	IsDryRun bool
}

// Plugin updates chart metadata and optionally packages a chart.
type Plugin struct {
	cfg PluginConfig
}

// PluginConfig configures the Helm plugin.
type PluginConfig struct {
	ChartDir         string
	UpdateAppVersion bool
	Publish          bool
	Publisher        *Publisher
}

// NewPlugin creates a Helm plugin.
func NewPlugin(cfg PluginConfig) *Plugin {
	if cfg.ChartDir == "" {
		cfg.ChartDir = "chart"
	}
	return &Plugin{cfg: cfg}
}

// Validate checks that the chart directory contains a valid Chart.yaml file.
func (p *Plugin) Validate() error {
	chartPath := filepath.Join(p.cfg.ChartDir, "Chart.yaml")
	meta, err := ReadChartMeta(chartPath)
	if err != nil {
		return fmt.Errorf("cannot read Chart.yaml at %s: %w", chartPath, err)
	}
	if meta.Name == "" {
		return fmt.Errorf("Chart.yaml is missing the 'name' field")
	}
	if p.cfg.Publish && p.cfg.Publisher == nil {
		return fmt.Errorf("publish is enabled but no Publisher is configured")
	}
	return nil
}

// Execute updates the chart version and optionally packages the chart.
func (p *Plugin) Execute(ctx context.Context, rel ReleaseContext) (*Result, error) {
	if rel.IsDryRun {
		return &Result{Outputs: map[string]string{
			"dry_run": "true",
			"version": rel.Version,
		}}, nil
	}

	chartPath := filepath.Join(p.cfg.ChartDir, "Chart.yaml")
	appVersion := ""
	if p.cfg.UpdateAppVersion {
		appVersion = rel.Version
	}

	meta, err := UpdateChartVersion(chartPath, rel.Version, appVersion)
	if err != nil {
		return nil, fmt.Errorf("update chart version: %w", err)
	}

	outputs := map[string]string{
		"chart_name":    meta.Name,
		"chart_version": meta.Version,
	}
	if meta.AppVersion != "" {
		outputs["app_version"] = meta.AppVersion
	}

	if p.cfg.Publish && p.cfg.Publisher != nil {
		tgz, err := p.cfg.Publisher.Package(ctx, p.cfg.ChartDir, "")
		if err != nil {
			return nil, fmt.Errorf("package chart: %w", err)
		}
		outputs["chart_tgz"] = tgz
	}

	return &Result{Outputs: outputs}, nil
}
