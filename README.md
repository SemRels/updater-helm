# updater-helm

[![Latest Release](https://img.shields.io/github/v/release/SemRels/updater-helm?label=version&color=blue)](https://github.com/SemRels/updater-helm/releases/latest)

Updates the version fields in a Helm chart.

This plugin is distributed as the standalone Go binary `semrel-plugin-updater-helm`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

### Binary

```bash
go install github.com/SemRels/updater-helm/cmd/plugin@latest
```

### Docker

Pre-built, multi-platform images (linux/amd64, linux/arm64) are published to the GitHub Container Registry on every release:

```bash
docker pull ghcr.io/semrels/updater-helm:latest
```

Images are signed with [cosign](https://github.com/sigstore/cosign) and include a full SBOM attestation. Verify the signature:

```bash
cosign verify ghcr.io/semrels/updater-helm:latest \
  --certificate-identity-regexp 'https://github.com/SemRels/updater-helm/.github/workflows/release.yml.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```


## Configuration

```yaml
plugins:
  - name: updater-helm
    path: ~/.semrel/plugins/semrel-plugin-updater-helm
    env:
      SEMREL_PLUGIN_FILE: "charts/app/Chart.yaml"
      SEMREL_PLUGIN_UPDATE_APP_VERSION: "true"
      SEMREL_PLUGIN_APP_VERSION_PREFIX: "v"
```

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_FILE` | Optional | Path to the Helm chart file to update. | `Chart.yaml` |
| `SEMREL_PLUGIN_UPDATE_APP_VERSION` | Optional | When `true`, update `appVersion` to match the resolved release version. | `false` |
| `SEMREL_PLUGIN_APP_VERSION_ONLY` | Optional | When `true`, only update `appVersion` and leave `version` unchanged. | `false` |
| `SEMREL_PLUGIN_APP_VERSION` | Optional | Explicit `appVersion` value to write instead of using the release version. | unset |
| `SEMREL_PLUGIN_APP_VERSION_PREFIX` | Optional | String to prepend to the computed `appVersion` value. | empty |

`SEMREL_PLUGIN_APP_VERSION_ONLY` also works with `SEMREL_PLUGIN_APP_VERSION` and `SEMREL_PLUGIN_APP_VERSION_PREFIX` when you want to update only the application version field.

## `SEMREL_*` release context used

| Variable | Description |
| --- | --- |
| `SEMREL_VERSION` | Resolved release version for the current run. |
| `SEMREL_NEXT_VERSION` | Next version computed by semrel for the release. |
| `SEMREL_DRY_RUN` | Whether semrel is running in dry-run mode. |

## Example behavior

By default, the plugin updates only `version` in `Chart.yaml`.

To keep `appVersion` in sync with the release version as well:

```yaml
plugins:
  - name: updater-helm
    env:
      SEMREL_PLUGIN_UPDATE_APP_VERSION: "true"
```

To prefix only `appVersion` while leaving `version` unchanged:

```yaml
plugins:
  - name: updater-helm
    env:
      SEMREL_PLUGIN_UPDATE_APP_VERSION: "true"
      SEMREL_PLUGIN_APP_VERSION_PREFIX: "v"
```

With a release version of `1.2.3`, this produces:

```yaml
version: 1.2.3
appVersion: v1.2.3
```

To set a completely custom `appVersion`:

```yaml
plugins:
  - name: updater-helm
    env:
      SEMREL_PLUGIN_APP_VERSION: "2.0.0"
      SEMREL_PLUGIN_APP_VERSION_PREFIX: "v"
```

## License

Apache-2.0
