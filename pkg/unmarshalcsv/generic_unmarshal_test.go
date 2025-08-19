package unmarshalcsv

import (
    "path/filepath"
    "testing"

    "github.com/xuri/excelize/v2"
)

// TestGenericCSV_UnmarshalAndNormalize ensures we can unmarshal CSV via the generic
// entrypoint and populate the generic alias fields using NormalizeAll.
func TestGenericCSV_UnmarshalAndNormalize(t *testing.T) {
    csvPath := filepath.Join("testdata", "sample.csv")

    var rows []UnmarshalledData
    if err := Unmarshal(&rows, csvPath, 0); err != nil {
        t.Fatalf("unmarshal csv: %v", err)
    }
    if len(rows) < 2 {
        t.Fatalf("expected at least 2 rows from sample.csv, got %d", len(rows))
    }

    NormalizeAll(rows)

    // Row 0 is egress example per fixture
    eg := rows[0]
    if eg.Direction != "egress" {
        t.Fatalf("expected first row direction egress, got %s", eg.Direction)
    }
    if eg.PolicyName != "frontend-to-backend" {
        t.Fatalf("PolicyName: expected frontend-to-backend, got %s", eg.PolicyName)
    }
    if eg.SubjectNamespace != "ns-a" {
        t.Fatalf("SubjectNamespace: expected ns-a, got %s", eg.SubjectNamespace)
    }
    if eg.SubjectSelector != "app=frontend" {
        t.Fatalf("SubjectSelector: expected app=frontend, got %s", eg.SubjectSelector)
    }
    if eg.PeerSpecifier != "10.0.0.0/24" {
        t.Fatalf("PeerSpecifier: expected 10.0.0.0/24, got %s", eg.PeerSpecifier)
    }
    if eg.Protocols != "TCP" {
        t.Fatalf("Protocols: expected TCP, got %s", eg.Protocols)
    }
    if eg.Ports != "80" {
        t.Fatalf("Ports: expected 80, got %s", eg.Ports)
    }

    // Row 1 is ingress example per fixture
    in := rows[1]
    if in.Direction != "ingress" {
        t.Fatalf("expected second row direction ingress, got %s", in.Direction)
    }
    if in.PolicyName != "allow-ingress-https" {
        t.Fatalf("PolicyName: expected allow-ingress-https, got %s", in.PolicyName)
    }
    if in.SubjectNamespace != "ns-b" {
        t.Fatalf("SubjectNamespace: expected ns-b, got %s", in.SubjectNamespace)
    }
    if in.SubjectSelector != "app=backend" {
        t.Fatalf("SubjectSelector: expected app=backend, got %s", in.SubjectSelector)
    }
    if in.PeerSpecifier != "10.1.0.0/24" {
        t.Fatalf("PeerSpecifier: expected 10.1.0.0/24, got %s", in.PeerSpecifier)
    }
    if in.Protocols != "TCP" {
        t.Fatalf("Protocols: expected TCP, got %s", in.Protocols)
    }
    if in.Ports != "443" {
        t.Fatalf("Ports: expected 443, got %s", in.Ports)
    }
}

// TestGenericXLSX_UnmarshalAndNormalize creates a temporary XLSX with the canonical
// header and one egress row, then validates generic alias fields post-normalization.
func TestGenericXLSX_UnmarshalAndNormalize(t *testing.T) {
    dir := t.TempDir()
    xlsx := filepath.Join(dir, "sample.xlsx")

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
    if err := xf.SaveAs(xlsx); err != nil {
        t.Fatal(err)
    }
    _ = xf.Close()

    var rows []UnmarshalledData
    if err := Unmarshal(&rows, xlsx, 0); err != nil {
        t.Fatalf("unmarshal xlsx: %v", err)
    }
    if len(rows) != 1 {
        t.Fatalf("expected 1 row, got %d", len(rows))
    }

    NormalizeAll(rows)
    eg := rows[0]
    if eg.PolicyName != "frontend-to-backend" || eg.SubjectNamespace != "ns-a" || eg.SubjectSelector != "app=frontend" || eg.PeerSpecifier != "10.0.0.0/24" || eg.Protocols != "TCP" || eg.Ports != "80" {
        t.Fatalf("unexpected normalized data: %+v", eg)
    }
}
