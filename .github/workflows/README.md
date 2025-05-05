# GitHub Actions Workflows

This directory contains our GitHub Actions workflow configurations.

## Available Workflows

### Basic Branch Build (`basic-branch-build.yml`)

This workflow handles our basic CI/CD pipeline with a focus on code quality and build verification.

#### Jobs

1. **Lint**
   - Uses golangci-lint 2.0.2 for code quality checks
   - Runs on all pushes and PRs
   - Acts as a merge gate for the main branch
   - Configuration:
     - Uses `.golangci.yml` for linting rules
     - Leverages `make lint` for consistent behavior
     - Fails only when merging to main, not on feature branch pushes

2. **Build**
   - Builds and tests the application
   - Runs after successful linting
   - Verifies:
     - Keytool build and tests
     - Server build and tests

#### Behavior

- **On Merge to Main**:
  - Runs both lint and build jobs
  - Must pass all checks
  - Acts as a merge gate

- **On Pull Requests**:
  - Runs both lint and build jobs
  - Must pass all checks to be mergeable
  - Prevents merging if checks fail

- **On Push to Feature Branches**:
  - Runs both jobs
  - Linting failures don't block pushes
  - Build failures still block pushes

## Usage

The workflow is automatically triggered on:
- Push to any branch
- Pull request creation/updates
- Manual workflow dispatch

## Configuration

- Linting rules are defined in `.golangci.yml`
- Build commands are defined in `Makefile`
- Go version is set to 1.24
