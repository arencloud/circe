package unmarshalcsv

import "testing"

func TestNormalize_Egress(t *testing.T) {
	row := UnmarshalledData{
		Direction:            "egress",
		SourceNamespace:      "ns-a",
		SourceSelector:       "app=frontend",
		DestinationSpecifier: "10.0.0.0/24,10.0.1.0/24",
		DestinationProtocol:  "TCP,UDP",
		DestinationPorts:     "80,443",
		NodeRole:             "worker",
		NetworkPolicyName:    "frontend-egress",
	}
	row.Normalize()

	if row.PolicyName != "frontend-egress" {
		t.Fatalf("PolicyName mismatch: %s", row.PolicyName)
	}
	if row.SubjectNamespace != "ns-a" {
		t.Fatalf("SubjectNamespace mismatch: %s", row.SubjectNamespace)
	}
	if row.SubjectSelector != "app=frontend" {
		t.Fatalf("SubjectSelector mismatch: %s", row.SubjectSelector)
	}
	if row.PeerSpecifier != "10.0.0.0/24,10.0.1.0/24" {
		t.Fatalf("PeerSpecifier mismatch: %s", row.PeerSpecifier)
	}
	if row.Protocols != "TCP,UDP" {
		t.Fatalf("Protocols mismatch: %s", row.Protocols)
	}
	if row.Ports != "80,443" {
		t.Fatalf("Ports mismatch: %s", row.Ports)
	}
	if row.Role != "worker" {
		t.Fatalf("Role mismatch: %s", row.Role)
	}
}

func TestNormalize_Ingress(t *testing.T) {
	row := UnmarshalledData{
		Direction:            "ingress",
		DestinationNamespace: "ns-b",
		DestinationSelector:  "app=backend",
		SourceSpecifier:      "10.1.0.0/24",
		DestinationProtocol:  "TCP",
		DestinationPorts:     "8080",
		NodeRole:             "",
		NetworkPolicyName:    "backend-ingress",
	}
	row.Normalize()

	if row.PolicyName != "backend-ingress" {
		t.Fatalf("PolicyName mismatch: %s", row.PolicyName)
	}
	if row.SubjectNamespace != "ns-b" {
		t.Fatalf("SubjectNamespace mismatch: %s", row.SubjectNamespace)
	}
	if row.SubjectSelector != "app=backend" {
		t.Fatalf("SubjectSelector mismatch: %s", row.SubjectSelector)
	}
	if row.PeerSpecifier != "10.1.0.0/24" {
		t.Fatalf("PeerSpecifier mismatch: %s", row.PeerSpecifier)
	}
	if row.Protocols != "TCP" {
		t.Fatalf("Protocols mismatch: %s", row.Protocols)
	}
	if row.Ports != "8080" {
		t.Fatalf("Ports mismatch: %s", row.Ports)
	}
}
