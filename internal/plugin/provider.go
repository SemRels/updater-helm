// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package plugin

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// Release contains the SemRel release data consumed by this plugin.
type Release struct {
	Version         string
	PreviousVersion string
	TagName         string
	Repository      string
	Changelog       string
	CommitSHA       string
	DryRun          bool
	Metadata        map[string]string
	Commits         []string
}

// Result captures the outcome of a plugin execution.
type Result struct {
	Name       string
	Outputs    map[string]string
	Skipped    bool
	SkipReason string
}

// Provider is the contract exposed by this plugin implementation.
type Provider interface {
	Name() string
	HealthCheck(context.Context) error
	Validate(map[string]interface{}) error
	Execute(context.Context, *Release) (*Result, error)
	ReleaseContext() []string
}

// HelmUpdater updates version fields inside a Helm Chart.yaml file.
type HelmUpdater struct {
	ChartPath string
}

// NewHelmUpdater constructs a Helm updater with an optional Chart.yaml path.
func NewHelmUpdater(chartPath string) *HelmUpdater {
	if chartPath == "" {
		chartPath = "Chart.yaml"
	}
	return &HelmUpdater{ChartPath: chartPath}
}

func (h *HelmUpdater) Name() string { return "updater-helm" }

func (h *HelmUpdater) HealthCheck(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (h *HelmUpdater) Validate(map[string]interface{}) error {
	if strings.TrimSpace(h.ChartPath) == "" {
		return fmt.Errorf("helm: chart path must not be empty")
	}
	return nil
}

func (h *HelmUpdater) ReleaseContext() []string {
	return []string{"version"}
}

func (h *HelmUpdater) Execute(ctx context.Context, rel *Release) (*Result, error) {
	if err := h.HealthCheck(ctx); err != nil {
		return nil, err
	}
	if err := h.Validate(nil); err != nil {
		return nil, err
	}
	if rel == nil {
		return nil, fmt.Errorf("helm: release is required")
	}
	if strings.TrimSpace(rel.Version) == "" {
		return nil, fmt.Errorf("helm: release version is required")
	}
	if rel.DryRun {
		return &Result{Name: h.Name(), Outputs: map[string]string{"chart": h.ChartPath, "version": rel.Version, "dry_run": "true"}}, nil
	}
	if err := UpdateChartVersion(h.ChartPath, rel.Version); err != nil {
		return nil, err
	}
	return &Result{Name: h.Name(), Outputs: map[string]string{"chart": h.ChartPath, "version": rel.Version}}, nil
}

// UpdateChartVersion rewrites version and appVersion fields in Chart.yaml.
func UpdateChartVersion(path, version string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("helm: read %s: %w", path, err)
	}

	var (
		lines         []string
		sawVersion    bool
		sawAppVersion bool
		scanner       = bufio.NewScanner(strings.NewReader(string(data)))
	)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]

		switch {
		case strings.HasPrefix(trimmed, "version:"):
			line = indent + "version: " + version
			sawVersion = true
		case strings.HasPrefix(trimmed, "appVersion:"):
			line = indent + fmt.Sprintf("appVersion: %q", version)
			sawAppVersion = true
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("helm: scan %s: %w", path, err)
	}
	if !sawVersion {
		lines = append(lines, "version: "+version)
	}
	if !sawAppVersion {
		lines = append(lines, fmt.Sprintf("appVersion: %q", version))
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
