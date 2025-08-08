// Package main demonstrates how to use go-vcard with Fiber framework
package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.rumenx.com/vcard"
)

// VCardHandler is a function that returns a VCard
type VCardHandler func(c *fiber.Ctx) *vcard.VCard

// VCardMiddleware creates a Fiber middleware for generating vCard responses
func VCardMiddleware(handler VCardHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate vCard
		card := handler(c)
		if card == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate vCard",
			})
		}

		// Validate vCard
		if err := card.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid vCard: %v", err),
			})
		}

		// Generate filename from name or use default
		filename := "contact.vcf"
		if name := card.GetFormattedName(); name != "" {
			// Create safe filename from name
			filename = strings.ReplaceAll(strings.ToLower(name), " ", "-") + ".vcf"
		}

		// Set headers
		c.Set("Content-Type", "text/vcard; charset=utf-8")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

		// Generate vCard content
		content, err := card.String()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to generate vCard content: %v", err),
			})
		}

		return c.SendString(content)
	}
}

// CreateVCardFromParams creates a vCard from URL parameters and form data
func CreateVCardFromParams(c *fiber.Ctx) *vcard.VCard {
	card := vcard.New()

	// Name from path parameters or query parameters
	firstName := c.Query("firstName")
	if firstName == "" {
		firstName = c.Params("firstName")
	}
	if firstName != "" {
		lastName := c.Query("lastName")
		if lastName == "" {
			lastName = c.Params("lastName")
		}
		card.AddName(firstName, lastName)
	}

	// Email
	if email := c.Query("email"); email != "" {
		emailType := c.Query("emailType", "work")
		var eType vcard.EmailType
		switch strings.ToLower(emailType) {
		case "home":
			eType = vcard.EmailHome
		case "mobile":
			eType = vcard.EmailMobile
		default:
			eType = vcard.EmailWork
		}
		card.AddEmail(email, eType)
	}

	// Phone
	if phone := c.Query("phone"); phone != "" {
		phoneType := c.Query("phoneType", "work")
		var pType vcard.PhoneType
		switch strings.ToLower(phoneType) {
		case "home":
			pType = vcard.PhoneHome
		case "mobile":
			pType = vcard.PhoneMobile
		case "fax":
			pType = vcard.PhoneFax
		default:
			pType = vcard.PhoneWork
		}
		card.AddPhone(phone, pType)
	}

	// Organization
	if org := c.Query("organization"); org != "" {
		card.AddOrganization(org)

		if title := c.Query("title"); title != "" {
			card.AddTitle(title)
		}
	}

	// URL
	if url := c.Query("url"); url != "" {
		card.AddURL(url, vcard.URLWork)
	}

	// Note
	if note := c.Query("note"); note != "" {
		card.AddNote(note)
	}

	return card
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "go-vcard Fiber Example",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Example 1: Simple vCard generation from path parameters
	app.Get("/vcard/:firstName/:lastName", VCardMiddleware(func(c *fiber.Ctx) *vcard.VCard {
		firstName := c.Params("firstName")
		lastName := c.Params("lastName")

		card := vcard.New()
		card.AddName(firstName, lastName)

		// Add optional parameters
		if email := c.Query("email"); email != "" {
			card.AddEmail(email, vcard.EmailWork)
		}
		if phone := c.Query("phone"); phone != "" {
			card.AddPhone(phone, vcard.PhoneWork)
		}
		if org := c.Query("org"); org != "" {
			card.AddOrganization(org)
		}

		return card
	}))

	// Example 2: Complex vCard from query parameters
	app.Get("/contact", VCardMiddleware(CreateVCardFromParams))

	// Example 3: Predefined contact
	app.Get("/me", VCardMiddleware(func(c *fiber.Ctx) *vcard.VCard {
		card := vcard.New()
		card.AddName("Bob", "Fiber")
		card.AddEmail("bob@fiber-example.com", vcard.EmailWork)
		card.AddPhone("+1-555-456-7890", vcard.PhoneWork)
		card.AddOrganization("Fiber Corp")
		card.AddTitle("Fiber Developer")
		card.AddURL("https://bobfiber.com", vcard.URLWork)
		card.AddNote("High-performance Go web framework specialist")
		return card
	}))

	// Example 4: JSON response with vCard data
	app.Get("/contact-json", func(c *fiber.Ctx) error {
		card := CreateVCardFromParams(c)
		if card == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No contact data provided"})
		}

		if err := card.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		vcardContent, err := card.String()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate vCard"})
		}

		return c.JSON(fiber.Map{
			"vcard": vcardContent,
			"data": fiber.Map{
				"name":         card.GetFormattedName(),
				"emails":       card.GetEmails(),
				"phones":       card.GetPhones(),
				"addresses":    card.GetAddresses(),
				"organization": card.GetOrganization(),
				"urls":         card.GetURLs(),
			},
		})
	})

	fmt.Println("Starting Fiber server on :8082")
	fmt.Println("Try these endpoints:")
	fmt.Println("  GET /vcard/Bob/Fiber?email=bob@fiber.com&phone=555-4567")
	fmt.Println("  GET /contact?firstName=Charlie&lastName=Brown&email=charlie@example.com&organization=Fiber")
	fmt.Println("  GET /me")
	fmt.Println("  GET /contact-json?firstName=Test&lastName=Fiber&email=test@fiber.com")

	app.Listen(":8082")
}
