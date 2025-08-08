// Package fiber provides Fiber framework adapter for go-vcard
package fiber

import (
	"github.com/gofiber/fiber/v2"
	"go.rumenx.com/vcard"
)

// VCardHandler is a function that returns a VCard
type VCardHandler func(c *fiber.Ctx) *vcard.VCard

// Options configures the vCard response
type Options struct {
	// Filename generates the filename for the vCard download
	Filename func(c *fiber.Ctx) string

	// ContentDisposition sets how the file should be handled (attachment/inline)
	ContentDisposition string
}

// DefaultOptions provides sensible defaults
var DefaultOptions = Options{
	Filename: func(c *fiber.Ctx) string {
		return "contact.vcf"
	},
	ContentDisposition: "attachment",
}

// VCard middleware for Fiber that generates vCard responses
func VCard(handler VCardHandler, opts ...Options) fiber.Handler {
	options := DefaultOptions
	if len(opts) > 0 {
		options = opts[0]
		// Apply defaults for missing fields
		if options.Filename == nil {
			options.Filename = DefaultOptions.Filename
		}
		if options.ContentDisposition == "" {
			options.ContentDisposition = DefaultOptions.ContentDisposition
		}
	}

	return func(c *fiber.Ctx) error {
		// Generate vCard
		card := handler(c)
		if card == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate vCard",
			})
		}

		// Generate vCard content
		content, err := card.String()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate vCard content",
			})
		}

		// Set headers
		filename := options.Filename(c)
		c.Set("Content-Type", "text/vcard")
		c.Set("Content-Disposition", options.ContentDisposition+"; filename="+filename)

		return c.SendString(content)
	}
}

// VCardJSON middleware for Fiber that returns vCard data as JSON
func VCardJSON(handler VCardHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate vCard
		card := handler(c)
		if card == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate vCard",
			})
		}

		// Convert to JSON-friendly structure
		response := fiber.Map{
			"name":         card.GetName(),
			"emails":       card.GetEmails(),
			"phones":       card.GetPhones(),
			"addresses":    card.GetAddresses(),
			"organization": card.GetOrganization(),
			"urls":         card.GetURLs(),
			"photo":        card.GetPhoto(),
			"note":         card.GetNote(),
		}

		return c.JSON(response)
	}
}

// CreateFromParams creates a vCard from Fiber context parameters and query values
func CreateFromParams(c *fiber.Ctx) *vcard.VCard {
	card := vcard.New()

	// Name from path parameters or query parameters
	if firstName := c.Query("firstName"); firstName == "" {
		firstName = c.Params("firstName")
		if firstName != "" {
			lastName := c.Query("lastName")
			if lastName == "" {
				lastName = c.Params("lastName")
			}
			card.AddName(firstName, lastName)
		}
	} else {
		lastName := c.Query("lastName")
		card.AddName(firstName, lastName)
	}

	// Email
	if email := c.Query("email"); email != "" {
		emailType := c.Query("emailType")
		switch emailType {
		case "home":
			card.AddEmail(email, vcard.EmailHome)
		case "mobile":
			card.AddEmail(email, vcard.EmailMobile)
		default:
			card.AddEmail(email, vcard.EmailWork)
		}
	}

	// Phone
	if phone := c.Query("phone"); phone != "" {
		phoneType := c.Query("phoneType")
		switch phoneType {
		case "home":
			card.AddPhone(phone, vcard.PhoneHome)
		case "mobile", "cell":
			card.AddPhone(phone, vcard.PhoneMobile)
		case "fax":
			card.AddPhone(phone, vcard.PhoneFax)
		default:
			card.AddPhone(phone, vcard.PhoneWork)
		}
	}

	// Organization
	if org := c.Query("organization"); org != "" {
		card.AddOrganization(org)
	}

	// Title
	if title := c.Query("title"); title != "" {
		card.AddTitle(title)
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
