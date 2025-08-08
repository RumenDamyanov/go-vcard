// Package main demonstrates how to use go-vcard with Echo framework
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.rumenx.com/vcard"
)

// VCardHandler is a function that returns a VCard
type VCardHandler func(c echo.Context) *vcard.VCard

// VCardMiddleware creates an Echo middleware for generating vCard responses
func VCardMiddleware(handler VCardHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Generate vCard
		card := handler(c)
		if card == nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to generate vCard",
			})
		}

		// Validate vCard
		if err := card.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
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
		c.Response().Header().Set("Content-Type", "text/vcard; charset=utf-8")
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

		// Generate vCard content
		content, err := card.String()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("Failed to generate vCard content: %v", err),
			})
		}

		return c.String(http.StatusOK, content)
	}
}

// CreateVCardFromParams creates a vCard from URL parameters and form data
func CreateVCardFromParams(c echo.Context) *vcard.VCard {
	card := vcard.New()

	// Name from path parameters or query parameters
	if firstName := c.QueryParam("firstName"); firstName == "" {
		firstName = c.Param("firstName")
	} else {
		lastName := c.QueryParam("lastName")
		if lastName == "" {
			lastName = c.Param("lastName")
		}
		card.AddName(firstName, lastName)
	}

	// Email
	if email := c.QueryParam("email"); email != "" {
		emailType := c.QueryParam("emailType")
		if emailType == "" {
			emailType = "work"
		}
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
	if phone := c.QueryParam("phone"); phone != "" {
		phoneType := c.QueryParam("phoneType")
		if phoneType == "" {
			phoneType = "work"
		}
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
	if org := c.QueryParam("organization"); org != "" {
		card.AddOrganization(org)

		if title := c.QueryParam("title"); title != "" {
			card.AddTitle(title)
		}
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

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Example 1: Simple vCard generation from path parameters
	e.GET("/vcard/:firstName/:lastName", VCardMiddleware(func(c echo.Context) *vcard.VCard {
		firstName := c.Param("firstName")
		lastName := c.Param("lastName")

		card := vcard.New()
		card.AddName(firstName, lastName)

		// Add optional parameters
		if email := c.QueryParam("email"); email != "" {
			card.AddEmail(email, vcard.EmailWork)
		}
		if phone := c.QueryParam("phone"); phone != "" {
			card.AddPhone(phone, vcard.PhoneWork)
		}
		if org := c.QueryParam("org"); org != "" {
			card.AddOrganization(org)
		}

		return card
	}))

	// Example 2: Complex vCard from query parameters
	e.GET("/contact", VCardMiddleware(CreateVCardFromParams))

	// Example 3: Predefined contact
	e.GET("/me", VCardMiddleware(func(c echo.Context) *vcard.VCard {
		card := vcard.New()
		card.AddName("Jane", "Echo")
		card.AddEmail("jane@echo-example.com", vcard.EmailWork)
		card.AddPhone("+1-555-987-6543", vcard.PhoneWork)
		card.AddOrganization("Echo Corp")
		card.AddTitle("Echo Developer")
		card.AddURL("https://janeecho.com", vcard.URLWork)
		card.AddNote("Echo framework specialist")
		return card
	}))

	// Example 4: JSON response with vCard data
	e.GET("/contact-json", func(c echo.Context) error {
		card := CreateVCardFromParams(c)
		if card == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "No contact data provided"})
		}

		if err := card.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		vcardContent, err := card.String()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate vCard"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"vcard": vcardContent,
			"data": map[string]interface{}{
				"name":         card.GetFormattedName(),
				"emails":       card.GetEmails(),
				"phones":       card.GetPhones(),
				"addresses":    card.GetAddresses(),
				"organization": card.GetOrganization(),
				"urls":         card.GetURLs(),
			},
		})
	})

	fmt.Println("Starting Echo server on :8081")
	fmt.Println("Try these endpoints:")
	fmt.Println("  GET /vcard/Jane/Echo?email=jane@echo.com&phone=555-9876")
	fmt.Println("  GET /contact?firstName=Alice&lastName=Johnson&email=alice@example.com&organization=Echo")
	fmt.Println("  GET /me")
	fmt.Println("  GET /contact-json?firstName=Test&lastName=Echo&email=test@echo.com")

	e.Logger.Fatal(e.Start(":8081"))
}
