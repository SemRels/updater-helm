// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateChartVersion(t *testing.T) {
	t.Parallel()

	chartPath := filepath.Join(t.TempDir(), "Chart.yaml")
	err := os.WriteFile(chartPath, []byte("apiVersion: v2\nname: demo\nversion: 0.1.0\nappVersion: \"0.1.0\"\n"), 0644)
	require.NoError(t, err)

	err = UpdateChartVersion(chartPath, "1.2.3")
	require.NoError(t, err)

	data, err := os.ReadFile(chartPath)
	require.NoError(t, err)
	require.Equal(t, "apiVersion: v2\nname: demo\nversion: 1.2.3\nappVersion: \"1.2.3\"\n", string(data))
}
