# Framework Integration Examples

This directory contains examples showing how to integrate go-vcard with popular Go web frameworks.

## Available Examples

| Framework | Port | Directory | Description |
|-----------|------|-----------|-------------|
| [Gin](./gin-adapter/) | 8080 | `gin-adapter/` | Gin web framework integration |
| [Echo](./echo-adapter/) | 8081 | `echo-adapter/` | Echo web framework integration |
| [Fiber](./fiber-adapter/) | 8082 | `fiber-adapter/` | Fiber web framework integration |
| [Chi](./chi-adapter/) | 8083 | `chi-adapter/` | Chi web framework integration |

## Quick Start

Each example can be run independently:

```bash
# Gin example
cd gin-adapter && go run main.go

# Echo example
cd echo-adapter && go run main.go

# Fiber example
cd fiber-adapter && go run main.go

# Chi example
cd chi-adapter && go run main.go
```

## Common API Endpoints

All examples implement the same API patterns:

### 1. Path Parameter vCard
```
GET /{framework}/vcard/:firstName/:lastName?email=...&phone=...&org=...
```

### 2. Query Parameter vCard
```
GET /{framework}/contact?firstName=...&lastName=...&email=...&organization=...&title=...
```

### 3. Predefined Contact
```
GET /{framework}/me
```

### 4. JSON API Response
```
GET /{framework}/contact-json?firstName=...&lastName=...&email=...
```

## Integration Patterns

Each example demonstrates:

### 🔧 **Middleware Pattern**
```go
func VCardMiddleware(handler VCardHandler) FrameworkHandlerFunc {
    return func(c FrameworkContext) FrameworkError {
        card := handler(c)
        // Validation, headers, response
    }
}
```

### 📝 **Parameter Extraction**
```go
func CreateVCardFromParams(c FrameworkContext) *vcard.VCard {
    card := vcard.New()

    // Extract from path params, query params, form data
    if firstName := extractParam(c, "firstName"); firstName != "" {
        card.AddName(firstName, extractParam(c, "lastName"))
    }

    return card
}
```

### ✅ **Validation & Error Handling**
```go
if err := card.Validate(); err != nil {
    return frameworkErrorResponse(400, err.Error())
}
```

### 📄 **Proper HTTP Headers**
```go
setHeader("Content-Type", "text/vcard; charset=utf-8")
setHeader("Content-Disposition", "attachment; filename=\"contact.vcf\"")
```

### 📁 **Safe Filename Generation**
```go
filename := "contact.vcf"
if name := card.GetFormattedName(); name != "" {
    filename = strings.ReplaceAll(strings.ToLower(name), " ", "-") + ".vcf"
}
```

### 🔄 **Multiple Response Formats**
- `.vcf` file download for vCard clients
- JSON response for API integration

## Test All Examples

You can test all examples simultaneously by running them in different terminals:

```bash
# Terminal 1 - Gin (port 8080)
cd gin-adapter && go run main.go

# Terminal 2 - Echo (port 8081)
cd echo-adapter && go run main.go

# Terminal 3 - Fiber (port 8082)
cd fiber-adapter && go run main.go
```

Then test with curl:

```bash
# Test Gin
curl "http://localhost:8080/vcard/John/Doe?email=john@gin.com" -o john-gin.vcf

# Test Echo
curl "http://localhost:8081/vcard/Jane/Echo?email=jane@echo.com" -o jane-echo.vcf

# Test Fiber
curl "http://localhost:8082/vcard/Bob/Fiber?email=bob@fiber.com" -o bob-fiber.vcf

# Test JSON responses
curl "http://localhost:8080/contact-json?firstName=Test&lastName=User&email=test@example.com"
```

## Framework-Specific Features

### Gin Advantages
- Extensive middleware ecosystem
- JSON binding and validation
- Built-in testing support
- Mature and stable

### Echo Advantages
- High performance
- Extensive middleware
- Built-in WebSocket support
- Automatic TLS

### Fiber Advantages
- Express.js-like API
- Zero memory allocation router
- Built-in rate limiting
- WebSocket support

## Adaptation Guide

To adapt these patterns for other frameworks:

1. **Replace Framework Context**: Change `gin.Context`, `echo.Context`, or `fiber.Ctx`
2. **Update Parameter Extraction**: Adapt `c.Param()`, `c.Query()`, `c.FormValue()` methods
3. **Modify Response Methods**: Change `c.JSON()`, `c.String()`, error handling
4. **Adjust Header Setting**: Update header setting methods
5. **Update Middleware Signature**: Match framework's middleware interface

## Production Considerations

- **Rate Limiting**: Add rate limiting middleware
- **Authentication**: Implement auth middleware if needed
- **Validation**: Add input validation and sanitization
- **Logging**: Include structured logging
- **Metrics**: Add prometheus metrics or similar
- **Security**: CORS, CSRF protection, secure headers
- **Caching**: Cache frequently generated vCards
- **File Storage**: Consider storing generated files temporarily

## Performance Tips

- **Connection Pooling**: Use connection pooling for database access
- **Caching**: Cache vCard templates and generated content
- **Compression**: Enable gzip compression for large vCards
- **Background Processing**: Use workers for complex vCard generation
- **Memory Management**: Be mindful of memory usage with large contact lists
