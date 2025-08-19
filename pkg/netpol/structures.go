package netpol

// Generic template capable of rendering both Ingress and Egress network policies
const NetworkPolicyGeneric = `
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  podSelector:
    matchLabels:
    {{- range $k, $v := .SelectorMap }}
      {{$k}}: {{$v}}
    {{- end }}
  policyTypes:
  - {{ .Direction }}
  {{- if eq .Direction "Egress" }}
  egress:
  - to:
    {{- range .PeerCIDRs }}
    - ipBlock:
        cidr: {{ . }}
    {{- end }}
    {{- if .Ports }}
    ports:
    {{- if eq (len .Protocols) 1 }}
    {{- $proto := index .Protocols 0 }}
    {{- range .Ports }}
    - protocol: {{$proto}}
      port: {{ . }}
    {{- end }}
    {{- else }}
    {{- range .Ports }}
    - protocol: TCP
      port: {{ . }}
    {{- end }}
    {{- range .Ports }}
    - protocol: UDP
      port: {{ . }}
    {{- end }}
    {{- end }}
    {{- end }}
  {{- else }}
  ingress:
  - from:
    {{- range .PeerCIDRs }}
    - ipBlock:
        cidr: {{ . }}
    {{- end }}
    {{- if .Ports }}
    ports:
    {{- if eq (len .Protocols) 1 }}
    {{- $proto := index .Protocols 0 }}
    {{- range .Ports }}
    - protocol: {{$proto}}
      port: {{ . }}
    {{- end }}
    {{- else }}
    {{- range .Ports }}
    - protocol: TCP
      port: {{ . }}
    {{- end }}
    {{- range .Ports }}
    - protocol: UDP
      port: {{ . }}
    {{- end }}
    {{- end }}
    {{- end }}
  {{- end }}`
