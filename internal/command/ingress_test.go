package command

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

// TestIngressCommand_Run exercises the ingress CLI command end-to-end using both CSV and XLSX inputs.
func TestIngressCommand_Run(t *testing.T) {
    // Use the shared CSV fixture from pkg/unmarshalcsv/testdata/sample.csv
    csvPath := filepath.Join("..", "..", "pkg", "unmarshalcsv", "testdata", "sample.csv")
    if _, err := os.Stat(csvPath); err != nil {
        t.Fatalf("fixture not found at %s: %v", csvPath, err)
    }

    // Subtest: CSV
    t.Run("csv", func(t *testing.T) {
        outDir := t.TempDir()
        cmd := NewIngressCommand()
        cmd.input = csvPath
        cmd.output = outDir
        cmd.headerStart = 0
        cmd.Run(nil, nil)

        outFile := filepath.Join(outDir, "allow-ingress-https.yaml")
        b, err := os.ReadFile(outFile)
        if err != nil {
            t.Fatalf("reading rendered file: %v", err)
        }
        s := string(b)
        wantSubs := []string{
            "apiVersion: networking.k8s.io/v1",
            "kind: NetworkPolicy",
            "name: allow-ingress-https",
            "namespace: ns-b",
            "app: backend",
            "policyTypes:",
            "- Ingress",
            "ingress:",
            "ipBlock:",
            "cidr: 10.1.0.0/24",
            "protocol: TCP",
            "port: 443",
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
        cmd := NewIngressCommand()
        cmd.input = xlsxPath
        cmd.output = outDir
        cmd.headerStart = 0
        cmd.Run(nil, nil)

        outFile := filepath.Join(outDir, "allow-ingress-https.yaml")
        b, err := os.ReadFile(outFile)
        if err != nil {
            t.Fatalf("reading rendered file: %v", err)
        }
        s := string(b)
        wantSubs := []string{
            "apiVersion: networking.k8s.io/v1",
            "kind: NetworkPolicy",
            "name: allow-ingress-https",
            "namespace: ns-b",
            "app: backend",
            "policyTypes:",
            "- Ingress",
            "ingress:",
            "ipBlock:",
            "cidr: 10.1.0.0/24",
            "protocol: TCP",
            "port: 443",
        }
        for _, sub := range wantSubs {
            if !strings.Contains(s, sub) {
                t.Fatalf("rendered YAML missing substring %q. Content:\n%s", sub, s)
            }
        }
    })
}
