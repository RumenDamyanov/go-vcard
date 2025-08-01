// Package main demonstrates how to use go-vcard with Gin framework
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rumendamyanov/go-vcard"
)

// VCardHandler is a function that returns a VCard
type VCardHandler func(c *gin.Context) *vcard.VCard

// VCardMiddleware creates a Gin middleware for generating vCard responses
func VCardMiddleware(handler VCardHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate vCard
		card := handler(c)
		if card == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate vCard",
			})
			return
		}

		// Validate vCard
		if err := card.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid vCard: %v", err),
			})
			return
		}

		// Generate filename from name or use default
		filename := "contact.vcf"
		if name := card.GetFormattedName(); name != "" {
			// Create safe filename from name
			filename = strings.ReplaceAll(strings.ToLower(name), " ", "-") + ".vcf"
		}

		// Set headers
		c.Header("Content-Type", "text/vcard; charset=utf-8")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

		// Generate vCard content
		content, err := card.String()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to generate vCard content: %v", err),
			})
			return
		}

		c.String(http.StatusOK, content)
	}
}

// CreateVCardFromParams creates a vCard from URL parameters and form data
func CreateVCardFromParams(c *gin.Context) *vcard.VCard {
	card := vcard.New()

	// Name from path parameters or form data
	if firstName := c.DefaultQuery("firstName", c.Param("firstName")); firstName != "" {
		lastName := c.DefaultQuery("lastName", c.Param("lastName"))
		card.AddName(firstName, lastName)
	}

	// Email
	if email := c.DefaultQuery("email", ""); email != "" {
		emailType := c.DefaultQuery("emailType", "work")
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
	if phone := c.DefaultQuery("phone", ""); phone != "" {
		phoneType := c.DefaultQuery("phoneType", "work")
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
	if org := c.DefaultQuery("organization", ""); org != "" {
		card.AddOrganization(org)

		if title := c.DefaultQuery("title", ""); title != "" {
			card.AddTitle(title)
		}
	}

	// URL
	if url := c.DefaultQuery("url", ""); url != "" {
		card.AddURL(url, vcard.URLWork)
	}

	// Note
	if note := c.DefaultQuery("note", ""); note != "" {
		card.AddNote(note)
	}

	return card
}

func main() {
	r := gin.Default()

	// Example 1: Simple vCard generation from URL parameters
	r.GET("/vcard/:firstName/:lastName", VCardMiddleware(func(c *gin.Context) *vcard.VCard {
		firstName := c.Param("firstName")
		lastName := c.Param("lastName")

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
	r.GET("/contact", VCardMiddleware(CreateVCardFromParams))

	// Example 3: Predefined contact
	r.GET("/me", VCardMiddleware(func(c *gin.Context) *vcard.VCard {
		card := vcard.New()
		card.AddName("John", "Developer")
		card.AddEmail("john@example.com", vcard.EmailWork)
		card.AddPhone("+1-555-123-4567", vcard.PhoneWork)
		card.AddOrganization("Example Corp")
		card.AddTitle("Senior Developer")
		card.AddURL("https://johndeveloper.com", vcard.URLWork)
		card.AddNote("Software engineer with 10+ years experience")
		return card
	}))

	// Example 4: JSON response with vCard data
	r.GET("/contact-json", func(c *gin.Context) {
		card := CreateVCardFromParams(c)
		if card == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No contact data provided"})
			return
		}

		if err := card.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		vcardContent, err := card.String()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate vCard"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"vcard": vcardContent,
			"data": gin.H{
				"name":         card.GetFormattedName(),
				"emails":       card.GetEmails(),
				"phones":       card.GetPhones(),
				"addresses":    card.GetAddresses(),
				"organization": card.GetOrganization(),
				"urls":         card.GetURLs(),
			},
		})
	})

	fmt.Println("Starting Gin server on :8080")
	fmt.Println("Try these endpoints:")
	fmt.Println("  GET /vcard/John/Doe?email=john@example.com&phone=555-1234")
	fmt.Println("  GET /contact?firstName=Jane&lastName=Smith&email=jane@example.com&organization=ACME")
	fmt.Println("  GET /me")
	fmt.Println("  GET /contact-json?firstName=Test&lastName=User&email=test@example.com")

	r.Run(":8080")
}
