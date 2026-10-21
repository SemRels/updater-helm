# updater-helm

[![Latest Release](https://img.shields.io/github/v/release/SemRels/updater-helm?label=version&color=blue)](https://github.com/SemRels/updater-helm/releases/latest)
[![CI](https://github.com/SemRels/updater-helm/actions/workflows/ci.yml/badge.svg)](https://github.com/SemRels/updater-helm/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/SemRels/updater-helm)](LICENSE)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/SemRels/updater-helm/badge)](https://scorecard.dev/viewer/?uri=github.com/SemRels/updater-helm)
[![REUSE status](https://api.reuse.software/badge/github.com/SemRels/updater-helm)](https://api.reuse.software/info/github.com/SemRels/updater-helm)

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
  --certificate-identity-regexp 'https://github.com/SemRels/updater-helm/.github/workflows/container.yml.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```


## Configuration

```yaml
plugins:
  - uses: "@semrel/updater-helm"
    phase: pre-tag
    args:
      file: "charts/app/Chart.yaml"
      update_app_version: "true"
      app_version_prefix: "v"
```

## Options

Every option can be set as an entry under `args` in `.semrel.yaml`. semrel passes each entry to the plugin as the environment variable shown in the table, so you can also set that variable directly (for example in your CI environment) instead of using `args`.

| Argument | Environment variable | Required | Description | Default |
| --- | --- | --- | --- | --- |
| `file` | `SEMREL_PLUGIN_FILE` | Optional | Path to the Helm chart file to update. | `Chart.yaml` |
| `update_app_version` | `SEMREL_PLUGIN_UPDATE_APP_VERSION` | Optional | When `true`, update `appVersion` to match the resolved release version. | `false` |
| `app_version_only` | `SEMREL_PLUGIN_APP_VERSION_ONLY` | Optional | When `true`, only update `appVersion` and leave `version` unchanged. | `false` |
| `app_version` | `SEMREL_PLUGIN_APP_VERSION` | Optional | Explicit `appVersion` value to write instead of using the release version. | unset |
| `app_version_prefix` | `SEMREL_PLUGIN_APP_VERSION_PREFIX` | Optional | String to prepend to the computed `appVersion` value. | empty |

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
  - uses: "@semrel/updater-helm"
    phase: pre-tag
    args:
      update_app_version: "true"
```

To prefix only `appVersion` while leaving `version` unchanged:

```yaml
plugins:
  - uses: "@semrel/updater-helm"
    phase: pre-tag
    args:
      update_app_version: "true"
      app_version_prefix: "v"
```

With a release version of `1.2.3`, this produces:

```yaml
version: 1.2.3
appVersion: v1.2.3
```

To set a completely custom `appVersion`:

```yaml
plugins:
  - uses: "@semrel/updater-helm"
    phase: pre-tag
    args:
      app_version: "2.0.0"
      app_version_prefix: "v"
```

## License

Apache-2.0
