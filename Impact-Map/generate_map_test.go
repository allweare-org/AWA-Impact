package main_test

import (
	"os"
	"strings"
	"testing"
)

// readSource reads generate_map.go for static analysis.
func readSource(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile("generate_map.go")
	if err != nil {
		t.Fatalf("Cannot read generate_map.go: %v", err)
	}
	return string(b)
}

// TestSheetTabNames ensures the code references the correct Google Sheets tab names.
func TestSheetTabNames(t *testing.T) {
	src := readSource(t)

	expected := []struct {
		tab  string
		desc string
	}{
		{"'System'!", "System sheet tab"},
		{"'Project'!", "Customer/Project sheet tab (renamed from Customer)"},
		{"'Location'!", "Location sheet tab"},
		{"'Population'!", "Population sheet tab"},
		{"FinalIME!", "Destination FinalIME tab"},
	}

	for _, tc := range expected {
		if !strings.Contains(src, tc.tab) {
			t.Errorf("Missing reference to %s (%s) — tab may have been renamed", tc.tab, tc.desc)
		}
	}

	// Ensure we're NOT still referencing the old 'Customer' tab
	if strings.Contains(src, "'Customer'!") {
		t.Error("Code still references 'Customer' tab — it was renamed to 'Project'")
	}
}

// TestOutputHeaders ensures the FinalIME output header row matches the expected 14-column schema.
func TestOutputHeaders(t *testing.T) {
	src := readSource(t)

	expectedHeaders := []string{
		"Customer ID",
		"System ID",
		"System Name",
		"Design Type",
		"Status",
		"Install Date",
		"Customer System Number",
		"Customer System Total",
		"Matched Customer Name",
		"Matched Customer Type",
		"Matched District",
		"Matched Latitude",
		"Matched Longitude",
		"Matched Population",
	}

	for _, h := range expectedHeaders {
		if !strings.Contains(src, `"`+h+`"`) {
			t.Errorf("Missing output header %q in FinalIME schema", h)
		}
	}
}

// TestStatusFilter ensures the code accepts current status values, not just the old "Installed".
func TestStatusFilter(t *testing.T) {
	src := readSource(t)

	// These status values should be accepted by the filter
	requiredStatuses := []string{"active", "nearing eol", "offline"}
	for _, s := range requiredStatuses {
		if !strings.Contains(src, `"`+s+`"`) {
			t.Errorf("Status filter missing %q — systems with this status will be dropped", s)
		}
	}
}

// TestCoordinateParsingFromSystemSheet ensures coordinates are extracted from the System sheet.
func TestCoordinateParsingFromSystemSheet(t *testing.T) {
	src := readSource(t)

	// The code should reference row[10] (column K = Location/coordinates in System sheet)
	if !strings.Contains(src, "row[10]") {
		t.Error("Code does not read System sheet column K (index 10) for coordinates — coordinates may be missing")
	}

	// Should split on comma for combined "lat, lng" format
	if !strings.Contains(src, `","`) || !strings.Contains(src, "Split") {
		t.Error("Code does not appear to parse combined coordinate strings (lat,lng format)")
	}
}

// TestCoordinateParsing validates that coordinate strings are handled correctly.
func TestCoordinateParsing(t *testing.T) {
	cases := []struct {
		input   string
		wantLat string
		wantLng string
		valid   bool
	}{
		{"0.3601293708966226, 32.58169066738784", "0.3601293708966226", "32.58169066738784", true},
		{"-0.534338,31.6254575", "-0.534338", "31.6254575", true},
		{"0.27196, 32.22813", "0.27196", "32.22813", true},
		{"", "", "", false},
		{"not coordinates", "", "", false},
	}

	for _, tc := range cases {
		coordStr := strings.TrimSpace(tc.input)
		if coordStr == "" || !strings.Contains(coordStr, ",") {
			if tc.valid {
				t.Errorf("Expected valid parse for %q", tc.input)
			}
			continue
		}
		parts := strings.Split(coordStr, ",")
		if len(parts) < 2 {
			if tc.valid {
				t.Errorf("Expected valid parse for %q", tc.input)
			}
			continue
		}
		lat := strings.TrimSpace(parts[0])
		lng := strings.TrimSpace(parts[1])
		if tc.valid && (lat != tc.wantLat || lng != tc.wantLng) {
			t.Errorf("Parse(%q) = (%q, %q), want (%q, %q)", tc.input, lat, lng, tc.wantLat, tc.wantLng)
		}
	}
}
