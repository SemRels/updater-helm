# updater-helm

Updates the version fields in a Helm chart.

This plugin is distributed as the standalone Go binary `semrel-plugin-updater-helm`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

```bash
go install github.com/SemRels/updater-helm/cmd/plugin@latest
```

## Configuration

```yaml
plugins:
  - name: updater-helm
    path: ~/.semrel/plugins/semrel-plugin-updater-helm
    env:
      SEMREL_PLUGIN_FILE: "charts/app/Chart.yaml"
```

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_FILE` | Optional | Path to the Helm chart file to update. | Chart.yaml |

## `SEMREL_*` release context used

| Variable | Description |
| --- | --- |
| `SEMREL_VERSION` | Resolved release version for the current run. |
| `SEMREL_NEXT_VERSION` | Next version computed by semrel for the release. |
| `SEMREL_DRY_RUN` | Whether semrel is running in dry-run mode. |

## Example behavior

The plugin updates the chart version fields in `Chart.yaml` to match the new release version.

## License

Apache-2.0
