package main

import (
	"fmt"

	"github.com/rumendamyanov/go-vcard"
)

func main() {
	card := vcard.New()
	card.AddName("John", "Doe")
	card.AddEmail("john.doe@example.com")
	card.AddPhone("+1234567890")
	card.AddAddress("123 Main St", "Anytown", "CA", "12345", "USA", vcard.AddressWork)

	content, err := card.String()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("=== vCard Output ===")
	fmt.Println(content)
	fmt.Println("=== End Output ===")
}
