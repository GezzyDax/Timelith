# Development Guide

This guide covers the development workflow, CI/CD pipeline, and best practices for contributing to Timelith.

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Development Workflow](#development-workflow)
3. [Testing](#testing)
4. [CI/CD Pipeline](#cicd-pipeline)
5. [Versioning](#versioning)
6. [Code Quality](#code-quality)

## Development Environment Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- Make (optional, but recommended)

### Quick Setup

```bash
# 1. Clone the repository
git clone https://github.com/GezzyDax/Timelith.git
cd Timelith

# 2. Set up environment
cp .env.example .env
# Edit .env with your credentials

# 3. Start infrastructure
make quick-start

# 4. Install dependencies
cd go-backend && go mod download && cd ..
cd web-ui && npm ci && cd ..
```

## Development Workflow

### Daily Workflow

1. **Start infrastructure**
   ```bash
   make quick-start
   ```

2. **Make your changes**
   - Edit code in `go-backend/` or `web-ui/`
   - Test locally as you develop

3. **Run checks before committing**
   ```bash
   make dev-check
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

### Commit Message Convention

We follow semantic commit messages:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Test changes
- `chore:` - Build process or tooling changes
- `ci:` - CI/CD changes

Examples:
```bash
git commit -m "feat: add user authentication"
git commit -m "fix: resolve database connection issue"
git commit -m "docs: update API documentation"
```

### Available Scripts

All scripts are located in `scripts/` directory:

#### quick-start.sh
Starts PostgreSQL and Redis containers for local development.

```bash
./scripts/quick-start.sh
# Or
make quick-start
```

#### dev-check.sh
Runs comprehensive pre-commit checks:
- Go formatting and linting
- Go tests with race detection
- TypeScript linting
- TypeScript type checking
- Build verification
- Docker Compose config validation

```bash
./scripts/dev-check.sh
# Or
make dev-check
```

#### test-all.sh
Runs full test suite:
- Go unit tests with coverage
- Docker image builds
- Integration tests

```bash
./scripts/test-all.sh
# Or
make test-all
```

#### clean-all.sh
Cleans all build artifacts and caches:
- Go build artifacts
- Node modules
- Docker containers and volumes

```bash
./scripts/clean-all.sh
# Or
make clean
```

#### bump-version.sh
Interactive version bumping tool.

```bash
./scripts/bump-version.sh
# Or
make bump-version
```

## Testing

### Backend Tests

```bash
# Run all tests
cd go-backend
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detection
go test -race ./...
```

### Frontend Type Checking

```bash
cd web-ui
npx tsc --noEmit
```

### Integration Tests

Integration tests are run automatically in CI, but you can run them locally:

```bash
./scripts/test-all.sh
```

## CI/CD Pipeline

### Overview

Our CI/CD pipeline consists of three main workflows:

1. **CI Pipeline** (`.github/workflows/ci.yml`)
2. **PR Checks** (`.github/workflows/pr-checks.yml`)
3. **Release** (`.github/workflows/release.yml`)
4. **CodeQL Security** (`.github/workflows/codeql.yml`)

### CI Pipeline

Runs on every push and pull request to `main` and `claude/**` branches.

**Jobs:**
- `go-lint` - Go code linting with golangci-lint
- `go-test` - Go unit tests with coverage
- `go-build` - Go binary build
- `web-lint` - ESLint and TypeScript checking
- `web-build` - Next.js production build
- `docker-backend` - Backend Docker image build
- `docker-web` - Web UI Docker image build
- `integration-test` - Full stack integration tests

**Features:**
- Parallel job execution
- Dependency caching (Go modules, npm packages, Docker layers)
- Automatic cancellation of outdated runs
- Artifact uploads (binaries, Docker images)

### PR Checks

Additional checks for pull requests:

- Security scanning with Trivy
- Dependency vulnerability checks
- Commit message validation
- PR statistics and changed files summary

### Release Workflow

Automatically runs on pushes to `main` branch:

1. **Version Detection**
   - Gets latest tag
   - Analyzes commits since last tag
   - Determines version bump type (major/minor/patch)

2. **Version Bump Logic**
   - `BREAKING CHANGE:` or `major:` → major bump
   - `feat:` or `feature:` → minor bump
   - Other commits → patch bump

3. **Changelog Generation**
   - Groups commits by type
   - Adds emojis for visual clarity
   - Creates GitHub release with notes

4. **Version File Updates**
   - Updates `go-backend/internal/version/version.go`
   - Updates `web-ui/package.json`
   - Commits changes with `[skip ci]`

5. **Release Creation**
   - Creates and pushes git tag
   - Creates GitHub release
   - Builds Docker images with version tags

### CodeQL Security Analysis

Runs weekly and on security-relevant changes:
- Static code analysis
- Security vulnerability detection
- Code quality checks

## Versioning

### Semantic Versioning

We follow [Semantic Versioning](https://semver.org/):

- **Major** (X.0.0) - Breaking changes
- **Minor** (0.X.0) - New features (backward compatible)
- **Patch** (0.0.X) - Bug fixes

### Automatic Versioning

Version bumps happen automatically on `main` branch based on commit messages:

```bash
# Patch bump (0.0.X)
git commit -m "fix: resolve login issue"

# Minor bump (0.X.0)
git commit -m "feat: add dark mode"

# Major bump (X.0.0)
git commit -m "feat: redesign API

BREAKING CHANGE: API endpoints have changed"
```

### Manual Versioning

For local development or testing:

```bash
make bump-version
```

This will:
1. Ask you to select bump type
2. Update version files
3. Generate changelog
4. Create git tag (optional)

## Code Quality

### Go Code Quality

We use golangci-lint with the following configuration:

**Enabled Linters:**
- errcheck - Check for unchecked errors
- gosimple - Simplify code
- govet - Vet examines Go source code
- ineffassign - Detect ineffectual assignments
- staticcheck - Static analysis
- unused - Find unused code
- gofmt - Format check
- goimports - Imports check
- gosec - Security check

Configuration file: `.golangci.yml`

```bash
# Run linter
cd go-backend
golangci-lint run

# Auto-fix issues
golangci-lint run --fix
```

### TypeScript Code Quality

**Tools:**
- ESLint for linting
- TypeScript compiler for type checking
- Prettier (via Next.js)

```bash
cd web-ui

# Lint
npm run lint

# Type check
npx tsc --noEmit
```

### Pre-commit Checks

Always run before committing:

```bash
make dev-check
```

This ensures:
- Code is properly formatted
- All tests pass
- No type errors
- Builds successfully

## Best Practices

### 1. Branch Strategy

- `main` - Production-ready code
- `claude/*` - Feature branches
- Pull requests for all changes

### 2. Commit Often

Make small, focused commits with clear messages.

### 3. Test Locally

Always run `make dev-check` before pushing.

### 4. Security

- Never commit secrets or credentials
- Use environment variables
- Keep dependencies updated

### 5. Documentation

- Update README for user-facing changes
- Update DEVELOPMENT.md for developer changes
- Add code comments for complex logic

## Troubleshooting

### CI Fails on Linting

```bash
# Run linter locally
cd go-backend
golangci-lint run
```

### CI Fails on Tests

```bash
# Run tests locally
make test-all
```

### Docker Build Fails

```bash
# Clean and rebuild
make clean
make build
```

### Version File Not Created

The release workflow creates `go-backend/internal/version/version.go` automatically. Don't create it manually.

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Next.js Documentation](https://nextjs.org/docs)
- [Docker Documentation](https://docs.docker.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
