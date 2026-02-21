# Contributing to gofilter

Thank you for your interest in contributing to gofilter! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git

### Setup

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/gofilter.git
   cd gofilter
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/sidneip/gofilter.git
   ```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test ./... -race

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Development Workflow

### Creating a Branch

```bash
git checkout -b feat/my-feature
# or
git checkout -b fix/my-bugfix
```

### Branch Naming Convention

- `feat/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Adding or updating tests
- `chore/` - Maintenance tasks

### Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <description>

[optional body]
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation only
- `refactor` - Code change that neither fixes a bug nor adds a feature
- `test` - Adding or correcting tests
- `chore` - Changes to the build process or auxiliary tools

**Examples:**
```
feat: add StartsWith filter operator
fix: handle nil pointers in nested field access
docs: add examples for geospatial filters
```

### Code Style

- Follow standard Go conventions (`go fmt`)
- Run `go vet ./...` before committing
- Keep functions focused and small
- Add comments for exported functions (GoDoc style)
- Write tests for new functionality

### Testing Guidelines

1. **Write tests first** when possible (TDD)
2. **Test edge cases**: nil values, empty slices, invalid inputs
3. **Use table-driven tests** for multiple scenarios
4. **Aim for >80% coverage** on new code

Example test structure:
```go
func TestFilterName(t *testing.T) {
    tests := []struct {
        name     string
        input    []MyType
        filter   Filter[MyType]
        expected []MyType
    }{
        {
            name:     "description of test case",
            input:    []MyType{...},
            filter:   Eq[MyType]("Field", "value"),
            expected: []MyType{...},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Apply(tt.input, tt.filter)
            // assertions
        })
    }
}
```

## Pull Request Process

1. **Update your branch** with the latest upstream changes:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Ensure all tests pass**:
   ```bash
   go test ./... -race
   ```

3. **Push your branch**:
   ```bash
   git push origin feat/my-feature
   ```

4. **Open a Pull Request** on GitHub with:
   - Clear description of changes
   - Link to related issues (if any)
   - Screenshots/examples if applicable

5. **Address review feedback** promptly

## What to Contribute

### Good First Issues

- Add new filter operators (e.g., `StartsWith`, `EndsWith`)
- Improve test coverage
- Add examples to documentation
- Fix typos or clarify documentation

### Feature Ideas

Check the [Roadmap](README.md#roadmap) for planned features, or open an issue to discuss new ideas.

### Reporting Bugs

When reporting bugs, please include:

1. Go version (`go version`)
2. Operating system
3. Minimal code to reproduce
4. Expected vs actual behavior
5. Error messages (if any)

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help others learn and grow

## Questions?

- Open an issue for questions about contributing
- Check existing issues before creating new ones

Thank you for contributing!
