# Output helpers (tables + paging)

Goal: kill copy/paste; keep output consistent.

## Tables

Use `internal/cmd/output_helpers.go:tableWriter(ctx)`:

- `--plain`: `os.Stdout` (no alignment, TSV-friendly)
- default: `tabwriter.Writer` (aligned columns)

Call pattern:

- `w, flush := tableWriter(cmd.Context())`
- `defer flush()`

## Pagination hint

Use `internal/cmd/output_helpers.go:printNextPageHint(u, token)`:

- prints to stderr
- exact format (tests depend on it): `# Next page: --page <token>`

