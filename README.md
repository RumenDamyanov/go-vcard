# go-vcard

[![CI](https://github.com/rumendamyanov/go-vcard/actions/workflows/ci.yml/badge.svg)](https://github.com/rumendamyanov/go-vcard/actions/workflows/ci.yml)
![CodeQL](https://github.com/rumendamyanov/go-vcard/actions/workflows/github-code-scanning/codeql/badge.svg)
![Dependabot](https://github.com/rumendamyanov/go-vcard/actions/workflows/dependabot/dependabot-updates/badge.svg)
[![codecov](https://codecov.io/gh/rumendamyanov/go-vcard/branch/master/graph/badge.svg)](https://codecov.io/gh/rumendamyanov/go-vcard)
[![Go Report Card](https://goreportcard.com/badge/github.com/rumendamyanov/go-vcard?)](https://goreportcard.com/report/github.com/rumendamyanov/go-vcard)
[![Go Reference](https://pkg.go.dev/badge/github.com/rumendamyanov/go-vcard.svg)](https://pkg.go.dev/github.com/rumendamyanov/go-vcard)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/rumendamyanov/go-vcard/blob/master/LICENSE.md)

A framework-agnostic Go module for generating vCard files (.vcf) compatible with major contact managers (iOS, Android, Gmail, iCloud, etc.). Inspired by [php-vcard](https://github.com/RumenDamyanov/php-vcard), this package works seamlessly with any Go web framework including Gin, Echo, Fiber, Chi, and standard net/http.

## Features

• **Framework-agnostic**: Use with Gin, Echo, Fiber, Chi, or standard net/http
• **vCard 3.0/4.0 support**: Generate standards-compliant vCard files
• **Rich content**: Supports names, emails, phones, addresses, organizations, photos, and more
• **Modern Go**: Type-safe, extensible, and robust (Go 1.22+)
• **High test coverage**: 90+% coverage with comprehensive test suite and CI/CD integration
• **Easy integration**: Simple API, drop-in for handlers/middleware
• **Production ready**: Used in production environments
• **Multiple formats**: File output and HTTP response generation
• **Framework examples**: Ready-to-use examples for popular Go web frameworks

## Quick Links

• 📖 [Installation](#installation)
• 🚀 [Usage Examples](#usage)
• 🔧 [Framework Adapters](#framework-adapters)
• 📚 [Documentation Wiki](https://github.com/RumenDamyanov/go-vcard/wiki)
• 🧪 [Testing & Development](#testing--development)
• 🤝 [Contributing](https://github.com/RumenDamyanov/go-vcard/blob/master/CONTRIBUTING.md)
• 🔒 [Security Policy](https://github.com/RumenDamyanov/go-vcard/blob/master/SECURITY.md)
• 💝 [Support & Funding](https://github.com/RumenDamyanov/go-vcard/blob/master/FUNDING.md)
• 📄 [License](#license)

## Installation

```bash
go get github.com/rumendamyanov/go-vcard
```

## Usage

### Basic Example (net/http)

```go
package main

import (
    "net/http"
    "github.com/rumendamyanov/go-vcard"
)

func vcardHandler(w http.ResponseWriter, r *http.Request) {
    // Create new vCard
    card := vcard.New()

    // Add contact information
    card.AddName("John", "Doe")
    card.AddEmail("john.doe@example.com")
    card.AddPhone("+1234567890")
    card.AddAddress("123 Main St", "Anytown", "CA", "12345", "USA")
    card.AddOrganization("Acme Corp")

    // Generate vCard content
    content, err := card.String()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send as downloadable file
    w.Header().Set("Content-Type", "text/vcard")
    w.Header().Set("Content-Disposition", "attachment; filename=\"john_doe.vcf\"")
    w.Write([]byte(content))
}

func main() {
    http.HandleFunc("/vcard", vcardHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Advanced Features

```go
card := vcard.New()

// Set vCard version (3.0 or 4.0)
card.SetVersion(vcard.Version40)

// Add comprehensive contact details
card.AddName("Jane", "Smith").
    AddMiddleName("Elizabeth").
    AddPrefix("Dr.").
    AddSuffix("PhD")

card.AddEmail("jane@company.com", vcard.EmailWork).
    AddEmail("jane.personal@gmail.com", vcard.EmailHome)

card.AddPhone("+1-555-123-4567", vcard.PhoneWork).
    AddPhone("+1-555-987-6543", vcard.PhoneHome).
    AddPhone("+1-555-555-5555", vcard.PhoneMobile)

// Add detailed address
card.AddAddress(
    "123 Business Ave",     // Street
    "Suite 100",            // Extended
    "Business City",        // City
    "CA",                   // Region
    "90210",               // Postal Code
    "United States",        // Country
    vcard.AddressWork,      // Type
)

// Add organization details
card.AddOrganization("Tech Corp Inc.").
    AddTitle("Senior Software Engineer").
    AddRole("Backend Developer")

// Add social/web presence
card.AddURL("https://janesmith.dev", vcard.URLWork).
    AddURL("https://linkedin.com/in/janesmith", vcard.URLSocial)

// Add photo (base64 encoded or URL)
card.AddPhoto("https://example.com/photo.jpg")

// Add custom properties
card.AddCustomProperty("X-CUSTOM-FIELD", "Custom Value")

// Save to file
err := card.SaveToFile("jane_smith.vcf")
```

## Framework Adapters

Ready-to-use examples for popular Go web frameworks are available in the [`examples/`](./examples/) directory:

| Framework | Port | Example | Description |
|-----------|------|---------|-------------|
| [Gin](./examples/gin-adapter/) | 8080 | [`gin-adapter`](./examples/gin-adapter/main.go) | Gin web framework integration |
| [Echo](./examples/echo-adapter/) | 8081 | [`echo-adapter`](./examples/echo-adapter/main.go) | Echo web framework integration |
| [Fiber](./examples/fiber-adapter/) | 8082 | [`fiber-adapter`](./examples/fiber-adapter/main.go) | Fiber web framework integration |
| [Chi](./examples/chi-adapter/) | 8083 | [`chi-adapter`](./examples/chi-adapter/main.go) | Chi web framework integration |

### Two Integration Approaches

**1. Import Adapter Packages** — For clean middleware integration:
```bash
go get github.com/rumendamyanov/go-vcard/adapters/gin    # or echo, fiber, chi
```

**2. Copy Example Applications** — For quick start with full applications:
```bash
# Clone and run complete example servers
git clone https://github.com/rumendamyanov/go-vcard.git
cd go-vcard/examples/gin-adapter && go run main.go
```

### Quick Start with Framework Examples

```bash
# Clone the repository
git clone https://github.com/rumendamyanov/go-vcard.git
cd go-vcard

# Run Gin example (port 8080)
cd examples/gin-adapter && go run main.go

# Run Echo example (port 8081)
cd examples/echo-adapter && go run main.go

# Run Fiber example (port 8082)
cd examples/fiber-adapter && go run main.go

# Run Chi example (port 8083)
cd examples/chi-adapter && go run main.go
```

### Test the Examples

```bash
# Download vCard file
curl "http://localhost:8080/vcard/John/Doe?email=john@example.com" -o contact.vcf

# Get JSON response
curl "http://localhost:8080/contact-json?firstName=Jane&lastName=Smith&email=jane@example.com"
```

### Integration Pattern Example

Each framework adapter follows the same pattern:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/rumendamyanov/go-vcard"
    ginadapter "github.com/rumendamyanov/go-vcard/adapters/gin"
)

func main() {
    r := gin.Default()

    r.GET("/contact/:name", ginadapter.VCard(func(c *gin.Context) *vcard.VCard {
        name := c.Param("name")

        card := vcard.New()
        card.AddName(name, "Doe")
        card.AddEmail(name + "@example.com")

        return card
    }))

    r.Run(":8080")
}
```

### Fiber Example

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/rumendamyanov/go-vcard"
    fiberadapter "github.com/rumendamyanov/go-vcard/adapters/fiber"
)

func main() {
    app := fiber.New()

    app.Get("/contact/:name", fiberadapter.VCard(func(c *fiber.Ctx) *vcard.VCard {
        name := c.Params("name")

        card := vcard.New()
        card.AddName(name, "Smith")
        card.AddPhone("+1234567890")

        return card
    }))

    app.Listen(":8080")
}
```

### Echo Example

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/rumendamyanov/go-vcard"
    echoadapter "github.com/rumendamyanov/go-vcard/adapters/echo"
)

func main() {
    e := echo.New()

    e.GET("/contact/:name", echoadapter.VCard(func(c echo.Context) *vcard.VCard {
        name := c.Param("name")

        card := vcard.New()
        card.AddName(name, "Johnson")
        card.AddEmail(name + "@company.com")

        return card
    }))

    e.Start(":8080")
}
```

### Chi Example

```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/rumendamyanov/go-vcard"
    chiadapter "github.com/rumendamyanov/go-vcard/adapters/chi"
)

func main() {
    r := chi.NewRouter()

    r.Get("/contact/{name}", chiadapter.VCard(func(w http.ResponseWriter, r *http.Request) *vcard.VCard {
        name := chi.URLParam(r, "name")

        card := vcard.New()
        card.AddName(name, "Wilson")
        card.AddPhone("+1234567890")

        return card
    }))

    http.ListenAndServe(":8080", r)
}
```

## Multiple Methods for Adding Information

### Add() vs AddItem()

You can add vCard information using either individual methods or structured data:

**Individual methods** — Simple, type-safe, chainable:

```go
// Recommended for most use cases
card.AddName("John", "Doe").
    AddEmail("john@example.com").
    AddPhone("+1234567890").
    AddAddress("123 Main St", "City", "State", "12345", "Country")
```

**Structured data** — Advanced, batch operations:

```go
// Add contact info with detailed structure
contact := vcard.Contact{
    Name: vcard.Name{
        First:  "John",
        Last:   "Doe",
        Middle: "William",
        Prefix: "Mr.",
        Suffix: "Jr.",
    },
    Emails: []vcard.Email{
        {Address: "john@work.com", Type: vcard.EmailWork},
        {Address: "john@home.com", Type: vcard.EmailHome},
    },
    Phones: []vcard.Phone{
        {Number: "+1234567890", Type: vcard.PhoneWork},
        {Number: "+1987654321", Type: vcard.PhoneMobile},
    },
}

card.AddContact(contact)
```

## Documentation

For comprehensive documentation and examples:

• 📚 [Quick Start Guide](https://github.com/RumenDamyanov/go-vcard/wiki/Quick-Start.md) - Get up and running quickly
• 🔧 [Basic Usage](https://github.com/RumenDamyanov/go-vcard/wiki/Basic-Usage.md) - Core functionality and examples
• 🚀 [Advanced Usage](https://github.com/RumenDamyanov/go-vcard/wiki/Advanced-Usage.md) - Advanced features and customization
• 🔌 [Framework Integration](https://github.com/RumenDamyanov/go-vcard/wiki/Framework-Integration.md) - Integration with popular frameworks
• 🎯 [Best Practices](https://github.com/RumenDamyanov/go-vcard/wiki/Best-Practices.md) - Performance tips and recommendations
• 🤝 [Contributing Guidelines](https://github.com/RumenDamyanov/go-vcard/blob/master/CONTRIBUTING.md) - How to contribute to this project
• 🔒 [Security Policy](https://github.com/RumenDamyanov/go-vcard/blob/master/SECURITY.md) - Security guidelines and vulnerability reporting
• 💝 [Funding & Support](https://github.com/RumenDamyanov/go-vcard/blob/master/FUNDING.md) - Support and sponsorship information

## Testing & Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Code Quality

```bash
# Run static analysis
go vet ./...

# Format code
go fmt ./...

# Run linter (if installed)
golangci-lint run
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](https://github.com/RumenDamyanov/go-vcard/blob/master/CONTRIBUTING.md) for details on:

• Development setup
• Coding standards
• Testing requirements
• Pull request process

## Security

If you discover a security vulnerability, please review our [Security Policy](https://github.com/RumenDamyanov/go-vcard/blob/master/SECURITY.md) for responsible disclosure guidelines.

## Support

If you find this package helpful, consider:

• ⭐ Starring the repository
• 💝 [Supporting development](https://github.com/RumenDamyanov/go-vcard/blob/master/FUNDING.md)
• 🐛 [Reporting issues](https://github.com/rumendamyanov/go-vcard/issues)
• 🤝 [Contributing improvements](https://github.com/RumenDamyanov/go-vcard/blob/master/CONTRIBUTING.md)

## License

[MIT License](https://github.com/RumenDamyanov/go-vcard/blob/master/LICENSE.md)

---
