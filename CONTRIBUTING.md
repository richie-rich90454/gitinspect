# Contributing to gitinspect

First off, thank you for considering contributing to gitinspect! It's people
like you that make gitinspect such a great tool.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)
- [Pull Requests](#pull-requests)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## Code of Conduct

This project and everyone participating in it is governed by our
[Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to
uphold this code. Please report unacceptable behavior.

## How Can I Contribute?

### Report Bugs

Check the [issue tracker](https://github.com/richie-rich90454/gitinspect/issues) to
see if the bug has already been reported. If not, open a new issue using the
**Bug Report** template.

### Suggest Enhancements

Open an issue using the **Feature Request** template. Include as much detail as
possible about the use case and expected behavior.

### Improve Documentation

Documentation is in the `docs/` directory (VitePress). See the
[docs README](docs/) for how to build and preview changes.

### Write Code

See [Development Setup](#development-setup) below.

## Development Setup

### Prerequisites

- Go 1.25+
- Node.js 20+ (for documentation)
- Make (optional)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/richie-rich90454/gitinspect.git
cd gitinspect

# Build
make build

# Run tests
make test

# Run linter
make lint
```

### Documentation

```bash
cd docs
npm install
npm run docs:dev    # Start dev server
npm run docs:build  # Build for production
```

## Development Workflow

1. **Fork** the repository
2. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/my-feature
   ```
3. **Make your changes** and add tests
4. **Run tests and linter**:
   ```bash
   make test
   make lint
   ```
5. **Commit** with a descriptive message
6. **Push** to your fork
7. **Open a Pull Request** against `main`

## Coding Standards

- Follow idiomatic Go conventions
- Run `golangci-lint` before submitting
- Add tests for new functionality
- Keep functions focused and small
- Use proper error handling (wrap errors with context)
- Use context for cancellation where appropriate

## Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move to..." not "Moves to...")
- Limit the first line to 72 characters
- Reference issues and pull requests liberally after the first line

## Pull Requests

- Fill in the [PR template](.github/PULL_REQUEST_TEMPLATE.md)
- Include tests for any new functionality
- Update documentation if needed
- Ensure CI passes (test + lint)
- Keep PRs focused — one feature or fix per PR

## Reporting Bugs

When filing a bug report, please include:

- **OS and version** (e.g., Ubuntu 22.04, macOS 14, Windows 11)
- **gitinspect version** (`gitinspect --version` or commit hash)
- **Steps to reproduce** the issue
- **Expected behavior**
- **Actual behavior**
- **Logs or error output** (use `--quiet=false` for verbose output)

## Suggesting Features

When suggesting a feature, please include:

- **Use case** — what problem does this solve?
- **Proposed solution** — how should it work?
- **Alternatives considered** — what other approaches did you consider?
- **Additional context** — screenshots, examples, etc.

---

Thank you for contributing! 🎉
