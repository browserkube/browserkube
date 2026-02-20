# BrowserKube – Copilot Coding Agent Instructions

## What This Repository Is
BrowserKube is a Kubernetes-native browser farm. It provides a Kubernetes operator and backend server that manage browser pods (Chrome, Firefox, etc.) on demand, exposing WebDriver, Playwright, and DevTools Protocol endpoints. It includes a React/TypeScript frontend UI, a Go backend server (`browserkube`), a Go Kubernetes operator, auxiliary Go commands, and Helm charts for deployment.

**Trust these instructions first. Search the codebase only when information here is incomplete or appears incorrect.**

---

## Repository Layout

```
browserkube/
├── backend/            # Main Go backend server (module: github.com/browserkube/browserkube)
│   ├── browserkube/    # Main server binary (browserkube/main.go is entry point)
│   ├── cmd/            # Additional binaries: sidecar, browser-updater, session-archiver
│   ├── pkg/            # Shared packages: api, http, provision, session, storage, wd, etc.
│   ├── Makefile        # Backend build/test/lint targets
│   ├── .golangci.yml   # Linter config (golangci-lint v1.64.6, line-length 160, tests excluded)
│   └── go.mod          # Go 1.26, module github.com/browserkube/browserkube
├── operator/           # Kubernetes operator (module: github.com/browserkube/browserkube/operator)
│   ├── api/v1/         # CRD types (Browser, BrowserSet custom resources)
│   ├── cmd/main.go     # Operator entry point
│   ├── internal/controller/ # Reconcilers
│   ├── config/crd/     # Generated CRD YAML (required by controller tests)
│   ├── Makefile        # Operator build/test targets (uses controller-gen, kustomize)
│   ├── Taskfile.yaml   # Proxies to Makefile: lint → make lint, test → make test
│   └── go.mod          # Go 1.26, Kubernetes 1.32, controller-runtime 0.20
├── frontend/app/       # React+TypeScript SPA (webpack, Node 14+, npm)
├── cmds/               # Standalone Go commands: clipboard, extension-installer, recorder
├── helm/charts/browserkube/  # Helm chart for deployment
├── .taskfiles/go.yaml  # Shared Task targets for Go modules (lint, fmt, test, build)
├── Taskfile.yaml       # Root task runner; includes operator, clipboard, extension-installer
├── go.work             # Go workspace: backend, operator, cmds/clipboard, cmds/extension-installer, cmds/recorder
├── skaffold.yaml       # Skaffold config for K8s dev/prod deployment
└── .github/workflows/  # CI: build-backend.yml, build-frontend.yml, build-images.yml
```

---

## Go Workspace
The repo uses a **Go workspace** (`go.work`, Go 1.26). All Go modules are declared in `go.work`. When adding a new module or dependency, update `go.work` accordingly.

- Backend module: `github.com/browserkube/browserkube` (in `backend/`)
- Operator module: `github.com/browserkube/browserkube/operator` (in `operator/`)
- Commands are in `cmds/clipboard`, `cmds/extension-installer`, `cmds/recorder`

---

## Build Instructions

### Go Backend

```bash
# Build a specific component (from repo root or backend/):
cd backend
make build component=browserkube   # produces backend/bin/browserkube/app.bin
make build component=sidecar       # produces backend/bin/sidecar/app.bin
# Cross-compiled: CGO_ENABLED=0 GOOS=linux GOARCH=amd64
```

Or use `go build ./...` from `backend/` for all packages (native arch, useful for validation).

### Go Operator

```bash
cd operator
make build        # runs manifests → generate → fmt → vet → go build -o bin/manager cmd/main.go
# NOTE: 'make build' also runs controller-gen code generation; always use make, not plain go build
```

### Frontend

```bash
cd frontend/app
npm install       # always run before building (prebuild also does this)
npm run build     # runs: webpack build
npm run lint      # eslint src/
npm run lint:styles  # stylelint **/*.scss
npm test          # or: npm run test:all (cross-env CI=true jest --passWithNoTests)
```

---

## Testing

### Backend Tests

```bash
cd backend
go test ./...
# NOTE: pkg/storage tests have a known pre-existing failure (mock context mismatch).
# The CI workflow runs: task all:test (from repo root)
```

### Operator Tests

The controller tests require `setup-envtest` and pre-generated CRDs. Run via Makefile:

```bash
cd operator
make test   # runs: manifests generate fmt vet setup-envtest, then go test with KUBEBUILDER_ASSETS set
# Running 'go test ./...' directly WILL FAIL because envtest binaries won't be set up
# and config/crd/bases must be populated by 'make manifests' first
```

E2E tests require a running Kind cluster: `make test-e2e`.

---

## Linting

### Backend

```bash
cd backend
make install-lint          # installs golangci-lint v1.64.6 into backend/bin/
make lint                  # runs bin/golangci-lint run --timeout=10m -v ./...
# Config: backend/.golangci.yml — enable-all with specific disables; line-length: 160; tests excluded
```

### Operator

```bash
cd operator
make lint    # downloads golangci-lint v1.63.4 into operator/bin/, then runs it
```

### Via Task (from repo root)

```bash
task all:lint    # lints clipboard, extension-installer, operator
task all:test    # tests clipboard, extension-installer, operator
task all:fmt     # formats clipboard, extension-installer, operator
```

---

## Code Generation (Operator)

After modifying any type in `operator/api/v1/`, **always regenerate**:

```bash
cd operator
make generate   # controller-gen DeepCopy methods
make manifests  # CRD YAML, RBAC, webhooks
```

These are prerequisites for `make build` and `make test` as well.

---

## CI Pipelines (`.github/workflows/`)

| File | Trigger | What it checks |
|---|---|---|
| `build-backend.yml` | push/PR to main/develop | golangci-lint on `backend/`, `make build` for `sidecar` and `browserkube`, `task all:test` |
| `build-frontend.yml` | push/PR to main/develop | `npm install`, webpack build, eslint, stylelint, jest |
| `build-images.yml` | manual (`workflow_dispatch`) | `skaffold build -p github-build` to ECR |

**To replicate CI locally before pushing:**

```bash
# Backend CI:
cd backend && make install-lint && make lint && make build component=browserkube && make build component=sidecar
cd <repo-root> && task all:test

# Frontend CI:
cd frontend/app && npm install && npm run build && npm run lint && npm run lint:styles && npm run test:all
```

---

## Key Architecture Facts

- **Backend** uses `go.uber.org/fx` for dependency injection. Entry point: `backend/browserkube/main.go`.
- **Operator** is Kubebuilder-scaffolded. Custom resources: `Browser` and `BrowserSet` (in `operator/api/v1/`). Reconcilers in `operator/internal/controller/`.
- **Sidecar** (`backend/cmd/sidecar`) runs inside each browser pod, proxying WebDriver/CDP.
- **Provision layer** (`backend/browserkube/internal/provision/k8s/`) creates K8s pods for browser sessions.
- **Swagger docs** are generated via `make swagger` in `backend/` (uses swag; output: `backend/browserkube/docs/`).
- **Mocks** live in `**/mocks/` subdirectories and are generated with `mockery`.
- Import prefix for `goimports`/`gci`: `github.com/browserkube/browserkube`.
- Backend `.golangci.yml`: `run.tests: false` — linter skips test files. Key disabled linters: `wsl`, `err113`, `exhaustruct`, `revive`, `mnd`.
