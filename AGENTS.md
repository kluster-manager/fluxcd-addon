# AGENTS.md

This file provides guidance to coding agents (e.g. Claude Code, claude.ai/code) when working with code in this repository.

## Repository purpose

Go module `github.com/kluster-manager/fluxcd-addon` — an OCM (Open Cluster Management) **addon** that installs and manages [FluxCD v2](https://fluxcd.io/) (`flux-system`) on managed clusters. It defines a `FluxCDConfig` CRD on the hub that captures per-cluster Flux install knobs (image registry, components, RBAC), and the addon manager renders the embedded `flux2` manifest tree into each managed cluster via the OCM `addon-framework`.

The produced binary is the single root-package `main.go` (`go build .` builds it). Drives Flux installs on spoke clusters by way of OCM `ManagedClusterAddOn` resources.

The local filesystem path is `open-cluster-management.io/fluxcd-addon`; the **actual Go module is `github.com/kluster-manager/fluxcd-addon`**.

## Architecture

- `main.go` — entry point (at the module root, not under `cmd/`). Builds a Cobra command and registers it via `open-cluster-management.io/addon-framework/pkg/version`.
- `pkg/manager/`:
  - `manager.go` — the addon manager runtime (uses `addon-framework` to register the addon and template manifests per-cluster).
  - `helpers.go` — value/registry-FQDN helpers used when templating.
  - `agent-manifests/flux2/` — embedded chart/manifest tree that gets pushed to each managed cluster. Includes `values.yaml`. Loaded via `go:embed`.
- `apis/fluxcd/v1alpha1/` — Kubebuilder API types:
  - `fluxcdconfig_types.go` — `FluxCDConfig`.
  - `groupversion_info.go`, `doc.go`, `types_test.go`, generated `zz_generated.deepcopy.go` / `openapi_generated.go`.
- `crds/` — generated CRD YAML (`fluxcd.open-cluster-management.io_fluxcdconfigs.yaml`) plus `lib.go` exposing it via `go:embed`.
- `Dockerfile.in` (PROD, distroless) + `Dockerfile.dbg` (debian) — two image variants (no UBI for this one).
- `hack/`, `Makefile` — AppsCode build harness (runs everything inside `ghcr.io/appscode/golang-dev`).
- `vendor/` — checked-in deps.

CRD API group is `fluxcd.open-cluster-management.io` (sits under the OCM domain, not AppsCode's, despite this being an AppsCode-maintained addon).

## Common commands

All Make targets run inside `ghcr.io/appscode/golang-dev` — Docker must be running.

- `make build` / `make all-build` — build host or all-platform binaries.
- `make gen` — regenerate clientset + manifests + openapi (`clientset manifests openapi`). Run after any change to `apis/fluxcd/v1alpha1/*_types.go`.
- `make manifests` — regenerate CRDs only.
- `make clientset` — regenerate client code.
- `make openapi` — regenerate OpenAPI definitions.
- `make fmt` — gofmt + goimports.
- `make lint` — golangci-lint.
- `make unit-tests` / `make test` — Go unit tests.
- `make verify` — `verify-gen verify-modules`; `go mod tidy && go mod vendor` must leave the tree clean.
- `make container` — build PROD and DBG images.
- `make push` — push both; `make docker-manifest` writes multi-arch manifests; `make release` is the full publish flow.
- `make push-to-kind` / `make deploy-to-kind` — load into Kind and Helm-install.
- `make install` / `make uninstall` / `make purge` — Helm install lifecycle.
- `make add-license` / `make check-license` — manage license headers.

Run a single Go test (requires a local Go toolchain):

```
go test ./apis/fluxcd/v1alpha1/... -run TestName -v
```

The README's quick-start runs the addon end-to-end via `clusteradm`:

```
kubectl apply -f api/config/samples/fluxcd_v1alpha1_fluxcdconfig.yaml
clusteradm addon enable --names fluxcd-addon --namespace flux-system --clusters c1
```

## Conventions

- Module path is `github.com/kluster-manager/fluxcd-addon`. Filesystem location under `open-cluster-management.io/` is for layout only — imports use the GitHub path.
- License: Apache-2.0 (`LICENSE`); new files need the standard "Copyright AppsCode Inc. and Contributors." header (`make add-license`).
- Sign off commits (`git commit -s`); contributions follow the DCO (`DCO` file).
- Vendor directory is checked in — `go mod tidy && go mod vendor` must leave the tree clean (enforced by `verify-modules`).
- Do not hand-edit `zz_generated.*.go`, `openapi_generated.go`, or anything under `crds/` — change `apis/fluxcd/v1alpha1/*_types.go` and re-run `make gen`.
- The `pkg/manager/agent-manifests/flux2/` tree is the **manifest bundle pushed to every managed cluster**. Adding/removing files there changes the cluster-side install; keep `values.yaml` keys in sync with the templating logic in `manager.go` / `helpers.go`.
- Two Dockerfiles, one binary — keep `Dockerfile.in` and `Dockerfile.dbg` in sync.
- `main.go` lives at the module root (not under `cmd/`) — match that layout if you add another binary, or move everyone under `cmd/` in a single change.
