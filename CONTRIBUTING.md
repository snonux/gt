# Contributing to gt

Thank you for your interest in contributing to `gt`! This document provides guidelines for building, testing, and submitting contributions.

## Quick Start

1. Fork the repository
2. Create a branch for your feature or fix
3. Make your changes
4. Run tests and build
5. Submit a pull request

## Building

### Using Mage (Recommended)

Install mage if you haven't already:
```bash
go install github.com/magefile/mage@latest
```

Then build the project:
```bash
mage build
```

This creates the `gt` binary in the project directory.

### Using Go Directly

```bash
go build -o gt ./cmd/gt
```

## Testing

Run all tests:
```bash
mage test
```

Run RPN-specific tests:
```bash
mage testRPN
# or
mage rpn
```

Run tests with coverage:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Code Style

- Follow Go best practices (see [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments))
- Add comprehensive documentation comments for all exported functions and types
- Use descriptive variable and function names
- Keep functions focused and modular
- Write tests for new functionality

## Pull Request Guidelines

### Before Submitting

1. Ensure all tests pass: `mage test`
2. Ensure the code builds: `mage build`
3. Format your code: `go fmt ./...`
4. Check for linting issues: `golangci-lint run ./...`
5. Update documentation if needed (README.md, godoc comments)

### PR Description

Include:
- A clear description of the changes
- Related issues (if any)
- Any breaking changes or migration notes
- Screenshots or examples for UI changes

### Review Process

1. At least one maintainer review is required
2. All tests must pass
3. Code coverage should not decrease
4. Follow the project's coding standards

## Development Workflow

### Typical Workflow

1. Create a branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make changes and commit:
   ```bash
   git add .
   git commit -m "Describe your changes"
   ```

3. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

4. Open a pull request on GitHub

### Versioning

The project uses semantic versioning. Version bumps are handled by the maintainers.

## Questions?

Open an issue or ask in the pull request discussion.
