// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

// Package plugin updates Helm Chart.yaml files in-place.
package plugin

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Updater updates Helm chart metadata.
type Updater struct{}

// NewUpdater creates an updater.
func NewUpdater() *Updater {
	return &Updater{}
}

// Update rewrites version and appVersion fields.
func (u *Updater) Update(path, version, appVersion string, appVersionOnly bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	updated, err := updateContent(string(data), version, appVersion, appVersionOnly)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func updateContent(content, version, appVersion string, appVersionOnly bool) (string, error) {
	if appVersion == "" {
		appVersion = version
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	lines := make([]string, 0)
	versionUpdated := false
	appUpdated := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]

		if strings.HasPrefix(trimmed, "version:") && !appVersionOnly {
			line = indent + "version: " + version
			versionUpdated = true
		}
		if strings.HasPrefix(trimmed, "appVersion:") {
			line = indent + "appVersion: " + appVersion
			appUpdated = true
		}

		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan Chart.yaml: %w", err)
	}
	if !appVersionOnly && !versionUpdated {
		return "", fmt.Errorf("version field not found in Chart.yaml")
	}
	if !appUpdated {
		lines = append(lines, "appVersion: "+appVersion)
	}
	return strings.Join(lines, "\n"), nil
}
