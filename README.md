# All We Are Impact Map

Interactive map showing All We Are's solar installation sites across Uganda, built with Mapbox GL JS and powered by live Google Sheets data.

**Live map:** https://allweare-org.github.io/AWA-Impact/

---

## Architecture

```
┌──────────────────┐       ┌──────────────────────┐       ┌─────────────────┐
│  Google Sheets   │       │   generate_map.go    │       │   index.html    │
│  (4 source tabs) │──────▶│  Joins & writes to   │──────▶│  Reads FinalIME │
│                  │       │  FinalIME sheet       │       │  via published  │
│  • System        │       │                      │       │  CSV URL        │
│  • Project       │       └──────────────────────┘       └─────────────────┘
│  • Location      │                                              │
│  • Population    │                                              ▼
└──────────────────┘                                     Mapbox GL JS Map
```

### Data Flow

1. **`generate_map.go`** reads from 4 source Google Sheets (System, Project, Location, Population)
2. Joins data by Customer ID, extracts coordinates from System sheet column K
3. Writes the joined result to the **FinalIME** tab on a destination spreadsheet
4. **`index.html`** fetches FinalIME as a published CSV and renders clustered map markers

---

## Files

| File | Purpose |
|------|---------|
| `index.html` | Root map page served by GitHub Pages |
| `Impact-Map/index.html` | Source map (same content, kept for development) |
| `Impact-Map/generate_map.go` | Go program that joins sheet data → FinalIME |
| `Impact-Map/generate_map_test.go` | Schema validation tests (CI) |
| `.github/workflows/schema-check.yml` | CI: runs schema tests on PRs |
| `CHANGELOG.md` | Change history |

---

## Running `generate_map.go`

### Prerequisites
- Go 1.21+
- `credentials.json` (Google service account key) in the repo root

### Run
```bash
cd Impact-Map
go run generate_map.go
```

### What it does
- Fetches data from 4 source sheets via Google Sheets API
- Filters systems with status: Active, Nearing EOL, or Offline
- Extracts coordinates from System sheet column K (`"lat, lng"` format)
- Joins with customer name, type, district, and population data
- Writes 14-column output to the FinalIME tab

---

## Google Sheets Setup

### Source Sheets
| Sheet | Tab | ID |
|-------|-----|----|
| Master System | `System` | `1V6UzkyQ6CHRN1RbXi039GJzp00QxmW49zIR6T-xlu1g` |
| Customer Master | `Project` | `1GJy4RzaC8ws8QQ5HMG1kONMO37kW6uFvJGcuBZit9iM` |
| Location Master | `Location` | `15R0AfrWDvkMN1VTfap-NNrH-gca7adRvphF-8O0EOVs` |
| Population Master | `Population` | `1qhinR30z8MdoyBlRgi-MnFegivSDMdwQ8oH2GbqcYFg` |

### Destination Sheet
| Sheet | Tab | ID |
|-------|-----|----|
| Destination | `FinalIME` | `1ASYXQ3Bdt0FWHPqG7nQ1ZVKDCzk0Lk2CsHiTHHP4jaA` |

### Publishing
The FinalIME tab must be **published to the web** (File → Share → Publish to web → FinalIME → CSV). The published URL includes `gid=206506688` to target the correct tab.

---

## CI / Schema Checks

On PRs that modify `Impact-Map/generate_map.go`, the schema-check workflow runs:

- **Go build** — ensures code compiles
- **TestSheetTabNames** — verifies correct tab name references
- **TestOutputHeaders** — validates 14-column FinalIME schema
- **TestStatusFilter** — ensures Active/Nearing EOL/Offline statuses are accepted
- **TestCoordinateParsingFromSystemSheet** — confirms coordinate extraction from column K

---

## Embedding in WordPress

```html
<iframe src="https://allweare-org.github.io/AWA-Impact/?v=2"
  title="AWA Impact Map" loading="lazy" allowfullscreen
  style="display:block;width:100%;height:600px;border:none;border-radius:16px;">
</iframe>
```

---

## Service Account Access

These service accounts need **Editor** access to the destination spreadsheet:
- `all-we-are-master-database@appspot.gserviceaccount.com`
- `awa-sa@all-we-are-496322.iam.gserviceaccount.com`

---

Created by Sean Ryan for All We Are
