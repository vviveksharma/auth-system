# Contributing to GuardRail

Thank you for your interest in contributing to GuardRail! We welcome contributions from the community.

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:
- A clear description of the problem
- Steps to reproduce the issue
- Expected vs actual behavior
- Your Go version and environment details

### Suggesting Features

Feature requests are welcome! Please open an issue with:
- A clear description of the feature
- Use cases and examples
- Any relevant code snippets or mockups

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Write clear, documented code** following Go best practices
3. **Add tests** for any new functionality
4. **Update documentation** if you're changing functionality
5. **Ensure all tests pass** with `go test ./...`
6. **Run `go fmt`** and `go vet` before committing
7. **Write meaningful commit messages** following conventional commits

### Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and concise

### Testing

- Write unit tests for new features
- Maintain or improve code coverage
- Test edge cases and error conditions

### Documentation

- Update README.md for user-facing changes
- Add godoc comments for exported types and functions
- Include examples in documentation

## Development Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/auth.git
cd auth

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build ./...
```

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Assume good intentions

## Questions?

Feel free to open an issue for any questions or clarifications!

---

Thank you for contributing to GuardRail! 🛡️
