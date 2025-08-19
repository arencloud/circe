package command

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

// TestEgressCommand_Run exercises the egress CLI command end-to-end using both CSV and XLSX inputs.
func TestEgressCommand_Run(t *testing.T) {
    // Shared CSV fixture
    csvPath := filepath.Join("..", "..", "pkg", "unmarshalcsv", "testdata", "sample.csv")
    if _, err := os.Stat(csvPath); err != nil {
        t.Fatalf("fixture not found at %s: %v", csvPath, err)
    }

    // Subtest: CSV
    t.Run("csv", func(t *testing.T) {
        outDir := t.TempDir()
        cmd := NewEgressCommand()
        cmd.input = csvPath
        cmd.output = outDir
        cmd.headerStart = 0
        cmd.Run(nil, nil)

        outFile := filepath.Join(outDir, "frontend-to-backend.yaml")
        b, err := os.ReadFile(outFile)
        if err != nil {
            t.Fatalf("reading rendered file: %v", err)
        }
        s := string(b)
        wantSubs := []string{
            "apiVersion: networking.k8s.io/v1",
            "kind: NetworkPolicy",
            "name: frontend-to-backend",
            "namespace: ns-a",
            "app: frontend",
            "policyTypes:",
            "- Egress",
            "egress:",
            "ipBlock:",
            "cidr: 10.0.0.0/24",
            "protocol: TCP",
            "port: 80",
        }
        for _, sub := range wantSubs {
            if !strings.Contains(s, sub) {
                t.Fatalf("rendered YAML missing substring %q. Content:\n%s", sub, s)
            }
        }
    })

    // Subtest: XLSX generated from the CSV fixture
    t.Run("xlsx", func(t *testing.T) {
        xlsxPath := filepath.Join(t.TempDir(), "sample.xlsx")
        if err := writeXLSXFromCSV(csvPath, xlsxPath); err != nil {
            t.Fatalf("failed to generate xlsx from csv: %v", err)
        }

        outDir := t.TempDir()
        cmd := NewEgressCommand()
        cmd.input = xlsxPath
        cmd.output = outDir
        cmd.headerStart = 0
        cmd.Run(nil, nil)

        outFile := filepath.Join(outDir, "frontend-to-backend.yaml")
        b, err := os.ReadFile(outFile)
        if err != nil {
            t.Fatalf("reading rendered file: %v", err)
        }
        s := string(b)
        wantSubs := []string{
            "apiVersion: networking.k8s.io/v1",
            "kind: NetworkPolicy",
            "name: frontend-to-backend",
            "namespace: ns-a",
            "app: frontend",
            "policyTypes:",
            "- Egress",
            "egress:",
            "ipBlock:",
            "cidr: 10.0.0.0/24",
            "protocol: TCP",
            "port: 80",
        }
        for _, sub := range wantSubs {
            if !strings.Contains(s, sub) {
                t.Fatalf("rendered YAML missing substring %q. Content:\n%s", sub, s)
            }
        }
    })
}
