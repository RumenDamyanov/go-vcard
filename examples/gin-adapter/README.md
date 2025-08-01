# Gin Adapter Example

This example demonstrates how to use go-vcard with the Gin web framework.

## Running the Example

```bash
go mod tidy
go run main.go
```

## API Endpoints

### 1. Simple vCard from Path Parameters
```
GET /vcard/:firstName/:lastName?email=john@example.com&phone=555-1234
```

Example: `http://localhost:8080/vcard/John/Doe?email=john@example.com&phone=555-1234`

### 2. Complex vCard from Query Parameters
```
GET /contact?firstName=Jane&lastName=Smith&email=jane@example.com&organization=ACME
```

Example: `http://localhost:8080/contact?firstName=Jane&lastName=Smith&email=jane@example.com&organization=ACME&title=Developer`

### 3. Predefined Contact
```
GET /me
```

Example: `http://localhost:8080/me`

### 4. JSON Response
```
GET /contact-json?firstName=Test&lastName=User&email=test@example.com
```

Example: `http://localhost:8080/contact-json?firstName=Test&lastName=User&email=test@example.com`

## Framework Integration Pattern

The example shows how to:

1. **Create a VCard Handler Function**: Define a function that takes a Gin context and returns a vCard
2. **Use Middleware Pattern**: Wrap handlers with vCard generation logic
3. **Handle Validation**: Validate vCards before sending response
4. **Set Proper Headers**: Content-Type and Content-Disposition for .vcf files
5. **Generate Safe Filenames**: Create filenames from contact names
6. **Support Multiple Response Formats**: Both .vcf download and JSON responses

## Key Features Demonstrated

- ✅ URL parameter extraction
- ✅ Query parameter handling
- ✅ Form data processing
- ✅ Error handling and validation
- ✅ Proper HTTP headers for vCard files
- ✅ JSON API responses
- ✅ Safe filename generation
- ✅ Type-safe email/phone/URL types

This pattern can be adapted for other Go web frameworks like Echo, Fiber, or Chi.
