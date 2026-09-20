# Repository audit — 2026-09-20

## Changes

Rooted output prevents path/symlink escape; reject nonlocal keys; create missing output parents.

## Validation

go test -race ./... passes with traversal, symlink and nested output regressions.

## Limits

Requires Go 1.25+ (previously 1.23.1). Existing in-root outputs may still be overwritten as before.
