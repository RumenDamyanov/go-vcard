# Contributing to go-vcard

Thank you for your interest in contributing to go-vcard! We welcome contributions from the community and are pleased to have you join us.

## Code of Conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Issues

Before creating an issue, please:

1. **Search existing issues** to avoid duplicates
2. **Check the documentation** to ensure it's not a usage question
3. **Use the issue templates** when available

When reporting bugs, include:
- Go version (`go version`)
- Operating system and architecture
- Minimal code example that reproduces the issue
- Expected vs actual behavior

### Suggesting Enhancements

Enhancement suggestions are welcome! Please:

1. **Check existing feature requests** first
2. **Provide clear use cases** for the enhancement
3. **Consider backwards compatibility**
4. **Include examples** of how it would work

### Pull Requests

We actively welcome pull requests. Here's how to get started:

#### Setup Development Environment

1. **Fork the repository**
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-vcard.git
   cd go-vcard
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/rumendamyanov/go-vcard.git
   ```

#### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
2. **Make your changes**
3. **Follow coding standards** (see below)
4. **Add tests** for new functionality
5. **Update documentation** if needed

#### Testing

Before submitting, ensure:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run static analysis
go vet ./...

# Format code
go fmt ./...
```

#### Submitting

1. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
2. **Create a pull request** from your fork to the main repository
3. **Fill out the PR template** completely
4. **Respond to feedback** promptly

## Coding Standards

### Go Code Style

- Follow standard Go conventions (`go fmt`, `go vet`)
- Use meaningful variable and function names
- Write clear, concise comments for public APIs
- Keep functions focused and reasonably sized

### Testing

- Write tests for new functionality
- Maintain or improve test coverage
- Use table-driven tests when appropriate
- Test edge cases and error conditions

### Documentation

- Update README.md for new features
- Add godoc comments for public APIs
- Include examples in documentation
- Update wiki documentation when needed

## Project Structure

```
go-vcard/
├── .github/          # GitHub workflows and templates
├── adapters/         # Framework adapters (gin, fiber, etc.)
├── examples/         # Usage examples
├── wiki/            # Local documentation
├── *.go            # Core library files
├── *_test.go       # Test files
└── README.md       # Main documentation
```

## Commit Message Guidelines

Use clear, descriptive commit messages:

```
feat: add support for custom vCard properties
fix: handle empty phone numbers correctly
docs: update README with new examples
test: add tests for email validation
refactor: simplify vcard generation logic
```

## Review Process

1. **Automated checks** must pass (CI, tests, linting)
2. **Code review** by maintainers
3. **Discussion** and feedback incorporation
4. **Approval** and merge

## Questions?

- **General questions**: Open a [Discussion](https://github.com/rumendamyanov/go-vcard/discussions)
- **Bug reports**: Create an [Issue](https://github.com/rumendamyanov/go-vcard/issues)
- **Security concerns**: See our [Security Policy](SECURITY.md)

## Recognition

Contributors are recognized in:
- README.md contributors section
- Release notes
- Project acknowledgments

Thank you for contributing to go-vcard! 🎉
