package unmarshalcsv

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestUnmarshalCsv_UnmarshalCsv(t *testing.T) {
	p := filepath.Join("testdata", "sample.csv")
	u, err := NewUnmarshalCsv(p, 0)
	if err != nil {
		t.Fatal(err)
	}

	n := []UnmarshalledData{}
	if err := u.UnmarshalCsv(&n); err != nil {
		t.Fatal(err)
	}
	if len(n) == 0 {
		t.Fatalf("expected rows, got 0")
	}
}

func TestUnmarshal_Generic_XLSX(t *testing.T) {
	// create a temporary xlsx file with the same header/one row
	dir := t.TempDir()
	file := filepath.Join(dir, "sample.xlsx")
	xf := excelize.NewFile()
	sheet := xf.GetSheetName(0)
	header := []string{"direction", "source_specifier", "destination_namespace", "destination_selector", "destination_protocol", "destination_ports", "source_namespace", "source_selector", "node_role", "destination_specifier", "comment", "network_policy_name"}
	row := []string{"egress", "", "ns-b", "app=backend", "TCP", "80", "ns-a", "app=frontend", "", "10.0.0.0/24", "", "frontend-to-backend"}
	for i, v := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := xf.SetCellStr(sheet, cell, v); err != nil {
			t.Fatal(err)
		}
	}
	for i, v := range row {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		if err := xf.SetCellStr(sheet, cell, v); err != nil {
			t.Fatal(err)
		}
	}
	if err := xf.SaveAs(file); err != nil {
		t.Fatal(err)
	}
	_ = xf.Close()

	var out []UnmarshalledData
	if err := Unmarshal(&out, file, 0); err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 row, got %d", len(out))
	}
	if out[0].Direction != "egress" || out[0].DestinationNamespace != "ns-b" {
		t.Fatalf("unexpected data: %+v", out[0])
	}

	// also test generic CSV path via Unmarshal
	csvPath := filepath.Join("testdata", "sample.csv")
	out = nil
	if err := Unmarshal(&out, csvPath, 0); err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Fatalf("expected csv rows")
	}
	_ = os.Remove(file)
}
