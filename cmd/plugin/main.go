// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	plugin "github.com/SemRels/updater-helm/internal/plugin"
)

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Getenv))
}

func run(stdout, stderr io.Writer, getenv func(string) string) int {
	version := getenv("SEMREL_VERSION")
	if version == "" {
		version = getenv("SEMREL_NEXT_VERSION")
	}
	if version == "" {
		fmt.Fprintln(stderr, "updater-helm: SEMREL_VERSION is required")
		return 1
	}
	version = strings.TrimPrefix(version, "v")

	file := getenv("SEMREL_PLUGIN_FILE")
	if file == "" {
		file = "Chart.yaml"
	}
	appOnly := strings.EqualFold(getenv("SEMREL_PLUGIN_APP_VERSION_ONLY"), "true")
	updateAppVersion := strings.EqualFold(getenv("SEMREL_PLUGIN_UPDATE_APP_VERSION"), "true")
	appVersionPrefix := getenv("SEMREL_PLUGIN_APP_VERSION_PREFIX")
	appVersion := getenv("SEMREL_PLUGIN_APP_VERSION")
	hasAppVersion := appVersion != ""
	if hasAppVersion {
		appVersion = strings.TrimPrefix(appVersion, "v")
	}
	shouldUpdateAppVersion := appOnly || updateAppVersion || hasAppVersion || appVersionPrefix != ""
	if shouldUpdateAppVersion {
		if appVersion == "" {
			appVersion = version
		}
		appVersion = appVersionPrefix + appVersion
	}

	if getenv("SEMREL_DRY_RUN") == "true" {
		fmt.Fprintf(stdout, "updater-helm: [dry-run] would update %s to version %s\n", file, version)
		return 0
	}

	if err := plugin.NewUpdater().Update(file, version, appVersion, appOnly, shouldUpdateAppVersion); err != nil {
		fmt.Fprintln(stderr, "updater-helm:", err)
		return 1
	}

	fmt.Fprintf(stdout, "updater-helm: updated %s to version %s\n", file, version)
	return 0
}
