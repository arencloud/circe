package command

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

// TestIngressCommand_Run exercises the ingress CLI command end-to-end using the CSV fixture.
// It ensures the command reads input CSV, generates ingress network policy YAML, and writes it to the output dir.
func TestIngressCommand_Run(t *testing.T) {
    // Use the shared CSV fixture from pkg/unmarshalcsv/testdata/sample.csv
    csvPath := filepath.Join("..", "..", "pkg", "unmarshalcsv", "testdata", "sample.csv")
    if _, err := os.Stat(csvPath); err != nil {
        t.Fatalf("fixture not found at %s: %v", csvPath, err)
    }

    outDir := t.TempDir()

    cmd := NewIngressCommand()
    cmd.input = csvPath
    cmd.output = outDir
    cmd.headerStart = 0

    // Run the command (it panics on error; the test will fail in that case)
    cmd.Run(nil, nil)

    // The ingress row in the fixture should produce allow-ingress-https.yaml
    outFile := filepath.Join(outDir, "allow-ingress-https.yaml")
    b, err := os.ReadFile(outFile)
    if err != nil {
        t.Fatalf("reading rendered file: %v", err)
    }
    s := string(b)

    // Basic content assertions (keep this aligned with the generator template)
    wantSubs := []string{
        "apiVersion: networking.k8s.io/v1",
        "kind: NetworkPolicy",
        "name: allow-ingress-https",
        "namespace: ns-b",
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
}
