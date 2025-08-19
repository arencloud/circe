# Circe

Circe is a small CLI and library to convert tabular policy definitions (CSV/XLSX) into Kubernetes NetworkPolicy YAML files.

It supports:
- Reading policy rows from CSV or XLSX files (library), and CSV via CLI.
- A generic, direction‑agnostic data model with normalization helpers.
- Rendering both Egress and Ingress NetworkPolicies from the same unified template.
- A helper command to generate sample CSV/XLSX files for quick starts.


## Table of Contents
- Overview
- Prerequisites
- Install / Build
- Quick Start
- CLI Usage
  - network-policy egress
  - network-policy ingress
- Input Schema (CSV/XLSX)
- Examples
- Troubleshooting / FAQ
- Version & License


## Overview
Given a simple spreadsheet that describes network policies (direction, pod selectors, peer CIDRs, protocols, and ports), Circe renders Kubernetes NetworkPolicy YAML files per policy. You can use it in two ways:
- CLI: Point it to a CSV file and an output directory, and it renders YAML files.
- Library: Import the package to unmarshal CSV/XLSX, normalize rows, and render using generic APIs.


## Prerequisites
- Go (1.22+ recommended). The module is set to `go 1.24`, but does not use bleeding-edge language features.
- A UNIX-like shell for the examples (Linux/macOS); Windows works too, just adjust paths.


## Install / Build
Because this repository uses a local module path (`module circe`), install by building from source:

- Clone the repository:
  git clone <your-fork-or-origin-url>
  cd circe

- Build the CLI binary:
  go build -o bin/circe ./cmd/main

- Optionally, run tests:
  go test ./...

The built binary will be at `bin/circe`.


## Quick Start
1) Generate a sample CSV with the correct schema:
   bin/circe generate --csv ./sample.csv

2) Render egress policies from the CSV to ./out:
   mkdir -p out
   bin/circe network-policy egress -i ./sample.csv -o ./out

3) Render ingress policies from the CSV to ./out:
   bin/circe network-policy ingress -i ./sample.csv -o ./out

4) Inspect results:
   ls ./out
   # Example files:
   #   frontend-to-backend.yaml
   #   allow-ingress-https.yaml


## CLI Usage

Root command:
- circe --help
- circe version

Subcommands are grouped under `network-policy`.

### network-policy egress
Generates Egress NetworkPolicy YAML from a CSV file.

Flags:
- -i, --input string        Path to the input CSV file (required)
- -o, --output string       Output directory for YAML files (default: current directory)
-     --header int          Header row index (0-based) in the CSV/XLSX; default 0

Example:
- bin/circe network-policy egress -i ./policies.csv -o ./out

Notes:
- Currently, the CLI path uses CSV input. XLSX is supported by the library APIs.

### network-policy ingress
Generates Ingress NetworkPolicy YAML from a CSV file.

Flags:
- -i, --input string        Path to the input CSV file (required)
- -o, --output string       Output directory for YAML files (default: current directory)
-     --header int          Header row index (0-based) in the CSV/XLSX; default 0

Example:
- bin/circe network-policy ingress -i ./policies.csv -o ./out



## Input Schema (CSV/XLSX)
Circe expects the following header row (order matters):

- direction
- source_specifier
- destination_namespace
- destination_selector
- destination_protocol
- destination_ports
- source_namespace
- source_selector
- node_role
- destination_specifier
- comment
- network_policy_name

You can inspect an example at `pkg/unmarshalcsv/testdata/sample.csv`. Sample rows:

- Egress example:
  egress,,ns-b,app=backend,TCP,80,ns-a,app=frontend,,10.0.0.0/24,,frontend-to-backend
- Ingress example:
  ingress,10.1.0.0/24,ns-b,app=backend,TCP,443,,,,,,allow-ingress-https

Semantics:
- direction: egress|ingress (case-insensitive)
- For egress: subject is the “source_*” namespace/selector; peers are from destination_specifier (CIDRs).
- For ingress: subject is the “destination_*” namespace/selector; peers are from source_specifier (CIDRs).
- destination_protocol: TCP, UDP, or both (comma-separated). Unknown protocols are ignored; TCP is the default if none provided.
- destination_ports: Comma-separated numeric ports.

Header row index: by default 0, use --header to change if your sheet has preamble rows.

## Examples
End-to-end (CSV to YAML):

- Input CSV (minimal):
  direction,source_specifier,destination_namespace,destination_selector,destination_protocol,destination_ports,source_namespace,source_selector,node_role,destination_specifier,comment,network_policy_name
  egress,,ns-b,app=backend,TCP,80,ns-a,app=frontend,,10.0.0.0/24,,frontend-to-backend

- Command:
  bin/circe network-policy egress -i ./policies.csv -o ./out

- Output (./out/frontend-to-backend.yaml), abbreviated:
  apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    name: frontend-to-backend
    namespace: ns-a
  spec:
    podSelector:
      matchLabels:
        app: frontend
    policyTypes:
    - Egress
    egress:
    - to:
      - ipBlock:
          cidr: 10.0.0.0/24
      ports:
      - protocol: TCP
        port: 80


## Troubleshooting / FAQ
- The command panics with a file error.
  - Ensure you pass -i/--input with a readable CSV file. Example: bin/circe network-policy egress -i ./file.csv -o ./out
- “unsupported file extension” error when using library Unmarshal.
  - Only .csv and .xlsx are supported. The CLI network-policy subcommands currently consume CSV; XLSX is supported in the library APIs.
- Ports or protocols look wrong in output.
  - Ensure destination_protocol is TCP and/or UDP (comma‑separated) and destination_ports are integers (comma‑separated). Unknown protocols are ignored. If none provided, TCP is assumed.
- Header not detected (empty output).
  - Check the header row index (--header). Default is 0; set it if your sheet starts later.


## Automated Releases (GitHub Actions)
This repository includes an automated release workflow using GitHub Actions. When you push a tag matching `v*` (e.g., `v0.2.0`) to the repository, the workflow will:
- Build cross-platform binaries for Linux, macOS (darwin), and Windows; architectures amd64 and arm64.
- Inject version metadata (Version, Commit short SHA, Date) into the binary via ldflags.
- Create a GitHub Release for the tag and attach the built binaries named like:
  - `circe_v0.2.0_linux_amd64`
  - `circe_v0.2.0_darwin_arm64`
  - `circe_v0.2.0_windows_amd64.exe`

How to cut a release:
- Ensure your changes are merged to main.
- Create and push a tag:
  git tag v0.2.0
  git push origin v0.2.0
- Wait for the "Release" workflow to complete; the release will appear under GitHub Releases with attached assets.

Note: The binaries embed version info; running `circe --version` will print, for example: `circe version v0.2.0 (commit abc1234, date 2025-08-19)`.

## Version & License
- Versioning: The CLI prints version information with `circe version` or `circe --version`.
  - By default it shows `dev`.
  - You can inject version metadata at build time using ldflags:
    go build -o bin/circe -ldflags "-X 'circe/internal/command.Version=v0.2.0' -X 'circe/internal/command.Commit=$(git rev-parse --short HEAD)' -X 'circe/internal/command.Date=$(date -u +%Y-%m-%d)'" ./cmd/main
  - Output example:
    circe version v0.2.0 (commit abc1234, date 2025-08-19)
- License: See [LICENSE](./LICENSE).


## Security Scanning (GitHub)
This repository enables automated vulnerability scanning using GitHub CodeQL.

- Workflow: .github/workflows/codeql.yml
- Triggers: on push/PR to main/master and on a weekly schedule.
- Where to see results: GitHub repository ➜ Security ➜ Code scanning alerts.

Notes:
- CodeQL analyzes Go code and reports potential vulnerabilities and security issues.
- You can fine-tune rules or add queries as needed in the CodeQL workflow.

### Secret Scanning guard (gitleaks)
To minimize the risk of committing sensitive data (secrets, passwords, API tokens), the repository also runs gitleaks on pushes and pull requests.

- Workflow: .github/workflows/gitleaks.yml
- Configuration: .gitleaks.toml (extends default rules and allowlists test data)
- Behavior: Fails CI if secrets are detected; output is redacted.

Recommended developer practices:
- Avoid committing generated artifacts or environment-specific files. The .gitignore excludes out/ and sample.csv.
- Use environment variables or CI secrets for credentials, never hard-code them in code or YAML.
