package netpol

import (
	"circe/pkg/unmarshalcsv"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type NetworkPolicy struct {
	generic []GenericPolicy
	output  string
}

// GenericPolicy is a unified representation for both Ingress and Egress policies
// Direction must be either "Ingress" or "Egress" (title case to match K8s spec)
type GenericPolicy struct {
	Name        string
	Namespace   string
	Selector    string
	SelectorMap map[string]string
	Direction   string
	PeerCIDRs   []string
	Ports       []string
	Protocols   []string // e.g., ["TCP"], ["UDP"], or ["TCP","UDP"]
}

// NewGenericPolicies builds a unified slice from CSV inputs for both directions
func NewGenericPolicies(input []unmarshalcsv.UnmarshalledData, output string) *NetworkPolicy {
	var gp []GenericPolicy
	for _, d := range input {
		name := d.NetworkPolicyName
		if name == "" {
			continue
		}

		protocols := normalizeProtocols(d.DestinationProtocol)
		ports := splitAndTrim(d.DestinationPorts)

		if strings.EqualFold(d.Direction, "egress") && d.SourceNamespace != "" && d.SourceSelector != "" {
   gp = append(gp, GenericPolicy{
				Name:        name,
				Namespace:   d.SourceNamespace,
				Selector:    d.SourceSelector,
				SelectorMap: parseSelector(d.SourceSelector),
				Direction:   "Egress",
				PeerCIDRs:   appendSlash(splitAndTrim(d.DestinationSpecifier)),
				Ports:       ports,
				Protocols:   protocols,
			})
		} else if strings.EqualFold(d.Direction, "ingress") && d.DestinationNamespace != "" && d.DestinationSelector != "" {
   gp = append(gp, GenericPolicy{
				Name:        name,
				Namespace:   d.DestinationNamespace,
				Selector:    d.DestinationSelector,
				SelectorMap: parseSelector(d.DestinationSelector),
				Direction:   "Ingress",
				PeerCIDRs:   appendSlash(splitAndTrim(d.SourceSpecifier)),
				Ports:       ports,
				Protocols:   protocols,
			})
		}
	}
	return &NetworkPolicy{generic: gp, output: output}
}

// NewGenericPoliciesForDirection is like NewGenericPolicies but filters to a single direction ("Egress" or "Ingress")
func NewGenericPoliciesForDirection(input []unmarshalcsv.UnmarshalledData, output string, direction string) *NetworkPolicy {
	var filtered []unmarshalcsv.UnmarshalledData
	for _, d := range input {
		if strings.EqualFold(direction, "egress") && strings.EqualFold(d.Direction, "egress") {
			filtered = append(filtered, d)
		} else if strings.EqualFold(direction, "ingress") && strings.EqualFold(d.Direction, "ingress") {
			filtered = append(filtered, d)
		}
	}
	return NewGenericPolicies(filtered, output)
}

// RenderGeneric renders the generic policies using the unified template
func (netpol *NetworkPolicy) RenderGeneric() error {
	tmpl := template.Must(template.New("generic").Parse(NetworkPolicyGeneric))
	if len(netpol.generic) == 0 {
		return fmt.Errorf("no generic policies defined")
	}
	for _, p := range netpol.generic {
		f, err := os.Create(fmt.Sprintf("%s/%s.yaml", netpol.output, p.Name))
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		if err := tmpl.Execute(f, p); err != nil {
			_ = f.Close()
			return fmt.Errorf("error executing template: %w", err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}
	}
	return nil
}

func appendSlash(in []string) []string {
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if strings.Contains(s, "/") {
			out = append(out, s)
		} else {
			out = append(out, fmt.Sprintf("%s/32", s))
		}
	}
	return out
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func normalizeProtocols(p string) []string {
	ps := splitAndTrim(p)
	seen := map[string]struct{}{}
	var out []string
	for _, v := range ps {
		u := strings.ToUpper(v)
		if u == "TCP" || u == "UDP" {
			if _, ok := seen[u]; !ok {
				seen[u] = struct{}{}
				out = append(out, u)
			}
		}
	}
	if len(out) == 0 {
		return []string{"TCP"} // default sensible fallback
	}
	return out
}

func parseSelector(s string) map[string]string {
	m := map[string]string{}
	if s == "" {
		return m
	}
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		var k, v string
		if idx := strings.Index(p, "="); idx >= 0 {
			k = strings.TrimSpace(p[:idx])
			v = strings.TrimSpace(p[idx+1:])
		} else if idx := strings.Index(p, ":"); idx >= 0 { // tolerate already "k: v" style
			k = strings.TrimSpace(p[:idx])
			v = strings.TrimSpace(p[idx+1:])
		} else {
			k = p
			v = ""
		}
		if k != "" {
			m[k] = v
		}
	}
	return m
}
