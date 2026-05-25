# updater-helm

`updater-helm` is a SemRel plugin that rewrites `version:` and `appVersion:` in a Helm `Chart.yaml` file.

## Behavior

- Defaults to `Chart.yaml` in the current working directory
- Updates `version:` to the SemRel release version
- Updates `appVersion:` to a quoted copy of the SemRel release version
- Supports dry-run execution through the plugin API

## Development

```bash
go mod tidy
go build ./...
go test ./...
```
