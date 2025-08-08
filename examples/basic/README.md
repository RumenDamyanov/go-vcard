# Basic vCard Example
A simple example of creating and saving a vCard.

## Code

```go
package main

import (
	"fmt"
	"go.rumenx.com/vcard"
)

func main() {
	// Create a new vCard
	card := vcard.New()

	// Add contact information
	card.AddName("John", "Doe")
	card.AddEmail("john.doe@example.com")
	card.AddPhone("+1-555-123-4567")
	card.AddAddress("123 Main Street", "Anytown", "CA", "12345", "USA")
	card.AddOrganization("Acme Corporation")
	card.AddTitle("Software Engineer")

	// Generate vCard content
	content, err := card.String()
	if err != nil {
		fmt.Printf("Error generating vCard: %v\n", err)
		return
	}

	// Print the vCard
	fmt.Println("Generated vCard:")
	fmt.Println(content)

	// Save to file
	err = card.SaveToFile("john_doe.vcf")
	if err != nil {
		fmt.Printf("Error saving vCard: %v\n", err)
		return
	}

	fmt.Println("vCard saved to john_doe.vcf")
}
```

## Running

```bash
cd examples/basic
go run main.go
```

This will create a file called `john_doe.vcf` that you can open with any contact manager application.
