# Chi Framework Adapter for go-vcard

This example demonstrates how to integrate the go-vcard library with the Chi web framework.

## Features

- **Middleware Integration**: Custom Chi middleware for vCard generation
- **Multiple Endpoints**: Support for both GET and POST requests
- **Path Parameters**: Extract contact info from URL paths
- **Query Parameters**: Flexible parameter handling
- **JSON Support**: Both JSON input/output and vCard file generation
- **Type Safety**: Proper Chi router patterns and parameter extraction
- **CORS Support**: Cross-origin request handling

## Running the Example

```bash
# Install dependencies
go mod tidy

# Run the server
go run main.go
```

The server will start on port 8083. Visit http://localhost:8083 to see available endpoints.

## API Endpoints

### GET /vcard/{firstName}/{lastName}
Download a vCard file with path and query parameters.

**Example:**
```bash
curl "http://localhost:8083/vcard/John/Doe?email=john@example.com&phone=+1234567890" -o contact.vcf
```

### GET /contact-json
Get contact information as JSON response.

**Example:**
```bash
curl "http://localhost:8083/contact-json?firstName=Jane&lastName=Smith&email=jane@example.com"
```

### POST /vcard
Create a vCard from JSON input.

**Example:**
```bash
curl -X POST http://localhost:8083/vcard \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Alice",
    "lastName": "Johnson",
    "email": "alice@example.com",
    "phone": "+1987654321",
    "organization": "Tech Corp",
    "title": "Software Engineer"
  }' -o alice.vcf
```

### GET /health
Health check endpoint returning server status.

## Integration Pattern

```go
import (
    "github.com/go-chi/chi/v5"
    vcard "github.com/rumendamyanov/go-vcard"
)

// Create Chi router
r := chi.NewRouter()

// Add middleware
r.Use(VCardMiddleware)

// Define routes
r.Get("/vcard/{firstName}/{lastName}", func(w http.ResponseWriter, r *http.Request) {
    firstName := chi.URLParam(r, "firstName")
    lastName := chi.URLParam(r, "lastName")

    // Create vCard
    vc := vcard.New()
    vc.AddName(firstName, lastName)

    // Generate content
    content, err := vc.String()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send response
    w.Header().Set("Content-Type", "text/vcard")
    w.Write([]byte(content))
})
```

## Key Chi Features Used

1. **URL Parameters**: `chi.URLParam()` for extracting path parameters
2. **Middleware**: Custom middleware for CORS and request processing
3. **Router Groups**: Organized route definition with Chi patterns
4. **Request Context**: Proper context handling for request processing
5. **Type Safety**: Chi's type-safe parameter extraction

## Dependencies

- Chi v5: `github.com/go-chi/chi/v5`
- go-vcard: `github.com/rumendamyanov/go-vcard`

## Testing

Test the endpoints using curl or any HTTP client:

```bash
# Test vCard download
curl "http://localhost:8083/vcard/Test/User?email=test@example.com" -o test.vcf

# Test JSON response
curl "http://localhost:8083/contact-json?firstName=Demo&lastName=User&email=demo@example.com"

# Test health endpoint
curl "http://localhost:8083/health"
```
