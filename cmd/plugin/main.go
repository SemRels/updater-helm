// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package main

import (
	"log"

	plugin "github.com/SemRels/updater-helm/internal/plugin"
)

func main() {
	helmPlugin := plugin.NewPlugin(plugin.PluginConfig{})
	log.Printf("updater-helm plugin ready: updates Helm chart metadata and optionally packages charts (%T)", helmPlugin)
}
