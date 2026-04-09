# 03 — Config Reference

> Field-by-field reference for `powershell.json`.

## Location

`powershell.json` must be in the repository root. If missing, `run.ps1` uses built-in defaults.

## Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `deployPath` | string | `E:\bin-run` (Win) / `/usr/local/bin` (Unix) | Directory where the compiled binary is deployed |
| `buildOutput` | string | `./bin` | Local directory for build artifacts |
| `binaryName` | string | `mahin.exe` | Name of the compiled binary file |
| `copyData` | boolean | `false` | If `true`, copies the `data/` directory alongside the binary during deploy |

## Example

```json
{
  "deployPath": "E:\\bin-run",
  "buildOutput": "./bin",
  "binaryName": "mahin.exe",
  "copyData": false
}
```

## Notes

- On macOS/Linux, set `binaryName` to `mahin` (no `.exe`).
- `deployPath` is created automatically if it doesn't exist.
- `copyData` is useful when the binary needs runtime data files (e.g., templates, seed configs).
