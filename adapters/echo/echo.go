// Package echo provides Echo framework adapter for go-vcard
package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.rumenx.com/vcard"
)

// VCardHandler is a function that returns a VCard
type VCardHandler func(c echo.Context) *vcard.VCard

// Options configures the vCard response
type Options struct {
	// Filename generates the filename for the vCard download
	Filename func(c echo.Context) string

	// ContentDisposition sets how the file should be handled (attachment/inline)
	ContentDisposition string
}

// DefaultOptions provides sensible defaults
var DefaultOptions = Options{
	Filename: func(c echo.Context) string {
		return "contact.vcf"
	},
	ContentDisposition: "attachment",
}

// VCard middleware for Echo that generates vCard responses
func VCard(handler VCardHandler, opts ...Options) echo.HandlerFunc {
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

	return func(c echo.Context) error {
		// Generate vCard
		card := handler(c)
		if card == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate vCard")
		}

		// Generate vCard content
		content, err := card.String()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate vCard content")
		}

		// Set headers
		filename := options.Filename(c)
		c.Response().Header().Set("Content-Type", "text/vcard")
		c.Response().Header().Set("Content-Disposition", options.ContentDisposition+"; filename="+filename)

		return c.String(http.StatusOK, content)
	}
}

// VCardJSON middleware for Echo that returns vCard data as JSON
func VCardJSON(handler VCardHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Generate vCard
		card := handler(c)
		if card == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate vCard")
		}

		// Convert to JSON-friendly structure
		response := map[string]interface{}{
			"name":         card.GetName(),
			"emails":       card.GetEmails(),
			"phones":       card.GetPhones(),
			"addresses":    card.GetAddresses(),
			"organization": card.GetOrganization(),
			"urls":         card.GetURLs(),
			"photo":        card.GetPhoto(),
			"note":         card.GetNote(),
		}

		return c.JSON(http.StatusOK, response)
	}
}

// CreateFromParams creates a vCard from Echo context parameters and query values
func CreateFromParams(c echo.Context) *vcard.VCard {
	card := vcard.New()

	// Name from path parameters or query parameters
	if firstName := c.QueryParam("firstName"); firstName == "" {
		firstName = c.Param("firstName")
		if firstName != "" {
			lastName := c.QueryParam("lastName")
			if lastName == "" {
				lastName = c.Param("lastName")
			}
			card.AddName(firstName, lastName)
		}
	} else {
		lastName := c.QueryParam("lastName")
		card.AddName(firstName, lastName)
	}

	// Email
	if email := c.QueryParam("email"); email != "" {
		emailType := c.QueryParam("emailType")
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
	if phone := c.QueryParam("phone"); phone != "" {
		phoneType := c.QueryParam("phoneType")
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
	if org := c.QueryParam("organization"); org != "" {
		card.AddOrganization(org)
	}

	// Title
	if title := c.QueryParam("title"); title != "" {
		card.AddTitle(title)
	}

	// URL
	if url := c.QueryParam("url"); url != "" {
		card.AddURL(url, vcard.URLWork)
	}

	// Note
	if note := c.QueryParam("note"); note != "" {
		card.AddNote(note)
	}

	return card
}
