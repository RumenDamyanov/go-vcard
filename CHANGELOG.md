# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial release of go-vcard
- Support for vCard 3.0 and 4.0 specifications
- Core vCard functionality:
  - Name management (structured and formatted names)
  - Email addresses with types (work, home, mobile)
  - Phone numbers with types (work, home, mobile, fax)
  - Postal addresses with types (work, home, postal)
  - Organization information (name, department, title, role)
  - URLs with types (work, home, social)
  - Photo support (URL and base64 data)
  - Birthday and anniversary support
  - Notes and custom properties
- Type-safe API with method chaining support
- Input validation and error handling
- File operations (save to file, read from bytes)
- Comprehensive test suite with 50%+ coverage
- Framework-agnostic design
- Project infrastructure:
  - MIT License
  - Contributing guidelines
  - Security policy
  - Code of conduct
  - Funding information
  - GitHub Actions CI/CD
  - Dependabot configuration
  - Comprehensive documentation
  - Wiki with guides and examples
  - Makefile for development tasks

### Features

- **Type Safety**: Strongly typed API with enums for email, phone, address, and URL types
- **Validation**: Built-in validation for required fields and data integrity
- **Method Chaining**: Fluent API for easy vCard construction
- **Multiple Formats**: Support for both structured data and individual property methods
- **Extensibility**: Custom properties support with X- prefix validation
- **Standards Compliance**: Adheres to RFC 2426 (vCard 3.0) and RFC 6350 (vCard 4.0)
- **Memory Efficiency**: Reusable vCard instances with Reset() and Clone() methods
- **Error Handling**: Comprehensive error reporting and graceful degradation

### Documentation

- Quick Start guide for immediate productivity
- Basic Usage guide covering core functionality
- Advanced Usage guide for complex scenarios
- Framework Integration guide (prepared for future adapters)
- Best Practices guide for production deployment
- Complete API documentation
- Multiple code examples
- Wiki-based documentation system

### Testing

- Unit tests for all core functionality
- Validation testing for edge cases
- Method chaining tests
- File I/O testing
- Memory management tests
- Error condition testing
- Coverage reporting with HTML output

### Development Infrastructure

- Makefile with common development tasks
- GitHub Actions workflow for CI/CD
- Dependabot for dependency management
- Code quality tools integration
- Linting and formatting configuration
- Examples directory for demonstrations

## [0.1.0] - TBD

### Initial Release

- Core vCard generation functionality
- Basic documentation and examples
- Alpha version ready for feedback
