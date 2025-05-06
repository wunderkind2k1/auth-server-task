# GitHub Configuration

This directory contains GitHub-specific configurations and tools.

## Workflows

### Basic Branch Build (`workflows/basic-branch-build.yml`)

This workflow handles our CI/CD pipeline with a focus on code quality, testing, and build verification.

#### Jobs

1. **Lint**
   - Uses golangci-lint 2.0.2 for code quality checks
   - Runs on all pushes and PRs
   - Configuration:
     - Uses `.golangci.yml` for linting rules
     - Leverages `make lint` for consistent behavior

2. **Test**
   - Runs tests and generates coverage reports
   - Calculates test-to-code ratio
   - Features:
     - Generates HTML coverage reports
     - Calculates test-to-code ratio
     - Shows warnings at 0.7 (70%)
     - Targets 1.0 (100%)
   - Uploads coverage reports as artifacts

3. **Build**
   - Builds the application
   - Runs after successful linting and testing
   - Verifies:
     - Keytool build
     - Server build
   - Displays test ratio summary

#### Behavior

The workflow has different requirements depending on the context:

- **When Targeting Main**:
  - All jobs (lint, test, build) must pass
  - Test ratio must meet minimum threshold (0.5)
  - Acts as a merge gate
  - Prevents merging if any check fails
  - Applies to:
    - Direct pushes to main
    - Pull requests targeting main

- **On Feature Branches**:
  - All jobs run but failures don't block
  - Results are reported but not enforced
  - Allows pushing regardless of check results
  - Encourages good practices without blocking development

## Scripts

### Test Ratio Calculator (`scripts/calculate_test_ratio.sh`)

A tool to calculate the test-to-code ratio, a standard metric that measures the relative size of test code compared to production code.

- Uses `cloc` for accurate code counting
- Ignores comments and blank lines
- Provides component-specific ratios
- Shows warnings at 0.7
- Targets 1.0
- Only enforces minimum ratio (0.5) when targeting main

## Usage

The workflow is automatically triggered on:
- Push to any branch
- Pull request creation/updates
- Manual workflow dispatch

## Configuration

- Linting rules are defined in `.golangci.yml`
- Build commands are defined in `Makefile`
- Go version is set to 1.24
- Test ratio thresholds are defined in `scripts/calculate_test_ratio.sh`
