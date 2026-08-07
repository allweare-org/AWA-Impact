# Changelog

## 2026-08-07 — Schema Migration & CI Validation

### Bug Fixes (`generate_map.go`)

- **Fixed status filter**: The System sheet now uses `Active`, `Nearing EOL`, and `Offline` instead of the old `Installed` value. Updated the filter to accept all three current statuses. Previously, every row was rejected → empty FinalIME → empty map.

- **Fixed coordinate source**: Coordinates moved from the Location sheet to the System sheet's column K (`Location`) as a combined `"lat, lng"` string. Added a pre-pass over System sheet rows to extract coordinates into a lookup map keyed by Customer ID. System-sheet coordinates take priority, with Location-sheet coordinates as fallback.

- **Fixed Customer tab name**: The Customer Master spreadsheet tab was renamed from `Customer` to `Project`. Updated the sheet range reference accordingly.

### New: CI Schema Validation

- **`Impact-Map/generate_map_test.go`** — 5 static analysis tests that run without Google credentials:
  - `TestSheetTabNames` — verifies correct tab names (`System`, `Project`, `Location`, `Population`, `FinalIME`) and flags stale references
  - `TestOutputHeaders` — ensures the 14-column FinalIME output schema is intact
  - `TestStatusFilter` — ensures `active`, `nearing eol`, `offline` are accepted by the filter
  - `TestCoordinateParsingFromSystemSheet` — ensures coordinates are read from System sheet column K (index 10)
  - `TestCoordinateParsing` — validates the `"lat, lng"` combined string splitting logic

- **`.github/workflows/schema-check.yml`** — GitHub Actions workflow that runs `go build` and the schema tests on any PR that modifies `Impact-Map/generate_map.go`
