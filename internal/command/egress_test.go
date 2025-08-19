package command

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

// TestEgressCommand_Run exercises the egress CLI command end-to-end using the CSV fixture.
// It ensures the command reads input CSV, generates egress network policy YAML, and writes it to the output dir.
func TestEgressCommand_Run(t *testing.T) {
    // Use the shared CSV fixture from pkg/unmarshalcsv/testdata/sample.csv
    csvPath := filepath.Join("..", "..", "pkg", "unmarshalcsv", "testdata", "sample.csv")
    if _, err := os.Stat(csvPath); err != nil {
        t.Fatalf("fixture not found at %s: %v", csvPath, err)
    }

    outDir := t.TempDir()

    cmd := NewEgressCommand()
    cmd.input = csvPath
    cmd.output = outDir
    cmd.headerStart = 0

    // Run the command (it panics on error; the test will fail in that case)
    cmd.Run(nil, nil)

    // The egress row in the fixture should produce frontend-to-backend.yaml
    outFile := filepath.Join(outDir, "frontend-to-backend.yaml")
    b, err := os.ReadFile(outFile)
    if err != nil {
        t.Fatalf("reading rendered file: %v", err)
    }
    s := string(b)

    // Basic content assertions (keep this aligned with the generator template)
    wantSubs := []string{
        "apiVersion: networking.k8s.io/v1",
        "kind: NetworkPolicy",
        "name: frontend-to-backend",
        "namespace: ns-a",
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
}
