# AGENTS.md

## Project Overview

BrowserKube is a Kubernetes-native platform for running browser sessions (WebDriver, Playwright, DevTools) in pods. It consists of a Go backend, a Kubernetes operator, a React frontend, and several sidecar/helper containers, all deployed via Helm + Skaffold.

## Repository Structure

```
backend/          Go module — main control plane and sidecar binaries
operator/         Go module — Kubernetes operator (controller-runtime)
frontend/app/     React + TypeScript UI (Webpack 5, MUI)
cmds/             Helper binaries and containers
  clipboard/        Go — clipboard HTTP service (chi)
  extension-installer/ Go — browser extension installer CLI (urfave/cli)
  recorder/         Go — X11 screen recorder via ffmpeg (urfave/cli)
  vnc-server/       Shell/Dockerfile — x11vnc container
  x-server/         Shell/Dockerfile — Xvfb + Openbox container
helm/charts/browserkube/  Helm chart (deploys all components)
.taskfiles/       Shared Taskfile templates for Go modules
.github/workflows/ CI: backend lint+test, frontend lint+test, image builds
```

## Go Workspace

The project uses `go.work` with modules: `backend`, `operator`, `cmds/clipboard`, `cmds/extension-installer`, `cmds/recorder`, `example/webdriver-go`. Go version is **1.26**.

The `backend` module depends on `operator` via a `replace` directive (`../operator`).

## Build System

**[Taskfile](https://taskfile.dev/) is the primary build tool.** Always use `task` to lint, test, format, and build. Do not call `make` or `go` directly unless a specific Makefile-only target is needed.

### Taskfile layout

- `Taskfile.yaml` — root orchestrator, includes per-module taskfiles
- `Taskfile_{darwin,linux,windows}.yaml` — OS-specific prerequisite installation
- `.taskfiles/go.yaml` — shared Go tasks reused by `backend`, `clipboard`, `extension-installer`
- `operator/Taskfile.yaml` — includes shared tasks but overrides `test` (proxies to `make test` for envtest)

### Common commands

```shell
task all:fmt              # format all Go modules
task all:lint             # lint all Go modules
task all:test             # test all Go modules

task backend:lint         # lint backend only
task backend:test         # test backend only
task operator:lint        # lint operator only
task operator:test        # test operator (envtest via Makefile)

task skaffold:dev         # local dev with Skaffold + port forwarding
```

### Makefiles (secondary, used by Taskfile or for specialized targets)

- `backend/Makefile` — `build component=<name>`, `swagger`, `generate-mocks`
- `operator/Makefile` — Kubebuilder-standard: `manifests`, `generate`, `test`, `build`, `docker-build`, `install`/`deploy` (kustomize)

### Skaffold

`skaffold.yaml` orchestrates Docker builds for all images and Helm deployment. Profiles: `dev` (live-reload), `arm64`, `remote-dev`, `github-build`.

## Linting

All Go modules share `.golangci.yml` (golangci-lint **v2.10.1**). Key rules:
- Formatters: `gofumpt`, `goimports` (local prefix `github.com/browserkube/browserkube`), `gci`
- Import ordering: standard → third-party → `github.com/browserkube/browserkube`
- Linter line length limit: 160 chars
- Enabled linters include: `govet` (shadow), `gosec`, `errcheck`, `bodyclose`, `dupl`, `staticcheck`, `tagalign`

## Architecture

### Backend (`backend/`)

Four binaries built from `cmd/`:

| Binary | Role |
|--------|------|
| `browserkube` | Main API server: WebDriver/Playwright proxy, REST API, session management |
| `sidecar` | Runs inside browser pods: proxies to local Selenium, serves downloads/videos |
| `session-archiver` | CronJob: archives session results between blob storage buckets |
| `browser-updater` | Job: syncs BrowserSet CRs with container registry image tags |

**DI pattern:** `browserkube` and `sidecar` use **uber/fx**. Modules are `fx.Options(...)` with `fx.Provide`/`fx.Invoke`. The app entrypoint is `pkg/app.Run(...)`. `session-archiver` and `browser-updater` use **urfave/cli** without fx.

**Key packages:**
- `pkg/http` — Chi router setup, health endpoint, CORS, server lifecycle
- `pkg/wd` — WebDriver reverse proxy with plugin pipeline
- `pkg/storage` — gocloud.dev blob abstraction (S3/GCS/file)
- `pkg/opentelemetry` — OTel metrics, Prometheus exporter
- `pkg/session`, `pkg/sessionresult` — domain models
- `pkg/websocketproxy` — WebSocket reverse proxy (gorilla/websocket)

Business logic lives in `cmd/browserkube/internal/` (API handlers, K8s provisioner, WD plugins, Playwright handler).

### Operator (`operator/`)

Kubebuilder-based operator using **controller-runtime v0.20.0**. API group: `api.browserkube.io/v1`.

**CRDs:**
- `Browser` — represents a single browser session (spec: browser name/version/type, status: phase/pod/URLs)
- `BrowserSet` — configuration catalog of available browsers per WebDriver/Playwright type
- `SessionResult` — archived session metadata (files, timestamps)

**Controller:** `BrowserReconciler` watches `Browser` CRs, reads `BrowserSet` for config, creates multi-container pods (browser + sidecar + VNC + recorder + clipboard + extension-installer), manages lifecycle (Pending → Running → Terminated).

**Typed client:** `pkg/client/v1/` — hand-written REST client for Browser, BrowserSet, SessionResult (used by the backend module).

### Frontend (`frontend/app/`)

React 18 + TypeScript SPA. Webpack 5 build. MUI v5 component library.

Key features: live session list, VNC viewer (react-vnc), terminal (xterm), session recording playback.

State: Redux Toolkit with WebSocket middleware. Routing: HashRouter.

## Testing

Run all tests via Taskfile: `task all:test`. Per-module: `task backend:test`, `task operator:test`, etc.

- **Backend:** `gotestsum ./...`. Mock generation with `mockery`.
- **Operator:** Ginkgo/Gomega + envtest (real API server, no cluster). E2e tests use Kind. Invoked via `make test` (proxied by `task operator:test`).
- **Frontend:** Jest + React Testing Library. ESLint + Stylelint.

## CI/CD

GitHub Actions workflows on push/PR to `main` and `develop`:
- `build-backend.yml` — Go lint (golangci-lint) + build + test
- `build-frontend.yml` — npm install + webpack build + eslint + stylelint + jest
- `build-images.yml` — manual dispatch, Skaffold build to ECR

## Conventions

- Go error handling: always handle errors, use `pkg/errors` for wrapping
- HTTP: Chi router, JSON responses via `pkg/http` helpers
- K8s resources: always set controller references and finalizers
- Struct tags: keep aligned (`tagalign` linter enforced)
- No OpenCensus — use OpenTelemetry exclusively (`depguard` enforced)
- Swagger docs generated from annotations (`swag init`)
