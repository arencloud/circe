package netpol_test

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    "circe/pkg/netpol"
    "circe/pkg/unmarshalcsv"
)

// TestEgressGeneratorEndToEnd validates the end-to-end flow:
//  CSV -> generic policies (egress-only) -> rendered YAML file with expected content.
func TestEgressGeneratorEndToEnd(t *testing.T) {
    // 1) Unmarshal CSV fixture
    csvPath := filepath.Join("..", "unmarshalcsv", "testdata", "sample.csv")
    var rows []unmarshalcsv.UnmarshalledData
    if err := unmarshalcsv.Unmarshal(&rows, csvPath, 0); err != nil {
        t.Fatalf("unmarshal csv: %v", err)
    }
    if len(rows) == 0 {
        t.Fatalf("expected sample CSV to contain rows")
    }

    // 2) Build generic egress policies and render to a temp dir
    outDir := t.TempDir()
    gp := netpol.NewGenericPoliciesForDirection(rows, outDir, "Egress")
    if err := gp.RenderGeneric(); err != nil {
        t.Fatalf("render generic egress: %v", err)
    }

    // 3) Expect a file for the egress row from the sample: frontend-to-backend
    outFile := filepath.Join(outDir, "frontend-to-backend.yaml")
    b, err := os.ReadFile(outFile)
    if err != nil {
        t.Fatalf("reading rendered file: %v", err)
    }
    s := string(b)

    // 4) Assert key content is present (basic sanity of output)
    // Note: We assert substrings rather than parsing YAML to keep the test focused
    // on the generator output shape without enforcing strict YAML schema here.
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
