package netpol_test

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    "circe/pkg/netpol"
    "circe/pkg/unmarshalcsv"
)

// TestIngressGeneratorEndToEnd validates the end-to-end flow:
//  CSV -> generic policies (ingress-only) -> rendered YAML file with expected content.
func TestIngressGeneratorEndToEnd(t *testing.T) {
    // 1) Unmarshal CSV fixture
    csvPath := filepath.Join("..", "unmarshalcsv", "testdata", "sample.csv")
    var rows []unmarshalcsv.UnmarshalledData
    if err := unmarshalcsv.Unmarshal(&rows, csvPath, 0); err != nil {
        t.Fatalf("unmarshal csv: %v", err)
    }
    if len(rows) == 0 {
        t.Fatalf("expected sample CSV to contain rows")
    }

    // 2) Build generic ingress policies and render to a temp dir
    outDir := t.TempDir()
    gp := netpol.NewGenericPoliciesForDirection(rows, outDir, "Ingress")
    if err := gp.RenderGeneric(); err != nil {
        t.Fatalf("render generic ingress: %v", err)
    }

    // 3) Expect a file for the ingress row from the sample: allow-ingress-https
    outFile := filepath.Join(outDir, "allow-ingress-https.yaml")
    b, err := os.ReadFile(outFile)
    if err != nil {
        t.Fatalf("reading rendered file: %v", err)
    }
    s := string(b)

    // 4) Assert key content is present (basic sanity of output)
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
