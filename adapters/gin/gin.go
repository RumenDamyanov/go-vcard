// Package gin provides Gin framework adapter for go-vcard
package gin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.rumenx.com/vcard"
)

// VCardHandler is a function that returns a VCard
type VCardHandler func(c *gin.Context) *vcard.VCard

// Options configures the vCard response
type Options struct {
	// Filename generates the filename for the vCard download
	Filename func(c *gin.Context) string

	// ContentDisposition sets how the file should be handled (attachment/inline)
	ContentDisposition string
}

// DefaultOptions provides sensible defaults
var DefaultOptions = Options{
	Filename: func(c *gin.Context) string {
		return "contact.vcf"
	},
	ContentDisposition: "attachment",
}

// VCard middleware for Gin that generates vCard responses
func VCard(handler VCardHandler, opts ...Options) gin.HandlerFunc {
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

		// Generate filename
		filename := options.Filename(c)
		if !strings.HasSuffix(strings.ToLower(filename), ".vcf") {
			filename += ".vcf"
		}

		// Set headers
		c.Header("Content-Type", "text/vcard; charset=utf-8")
		c.Header("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"",
			options.ContentDisposition, filename))

		// Send vCard content
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

// VCardJSON middleware that returns vCard data as JSON
func VCardJSON(handler VCardHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		card := handler(c)
		if card == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate vCard",
			})
			return
		}

		if err := card.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid vCard: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"vcard": func() string {
				content, err := card.String()
				if err != nil {
					return ""
				}
				return content
			}(),
			"data": map[string]interface{}{
				"name":         card.GetName(),
				"emails":       card.GetEmails(),
				"phones":       card.GetPhones(),
				"addresses":    card.GetAddresses(),
				"organization": card.GetOrganization(),
				"urls":         card.GetURLs(),
				"photo":        card.GetPhoto(),
				"birthday":     card.GetBirthday(),
				"anniversary":  card.GetAnniversary(),
				"note":         card.GetNote(),
			},
		})
	}
}

// FromParams creates a vCard from Gin context parameters and form data
func FromParams(c *gin.Context) *vcard.VCard {
	card := vcard.New()

	// Name from path parameters or form data
	if firstName := c.DefaultPostForm("firstName", c.Param("firstName")); firstName != "" {
		lastName := c.DefaultPostForm("lastName", c.Param("lastName"))
		card.AddName(firstName, lastName)
	}

	// Email
	if email := c.DefaultPostForm("email", c.Param("email")); email != "" {
		emailType := c.DefaultPostForm("emailType", "work")
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
	if phone := c.DefaultPostForm("phone", c.Param("phone")); phone != "" {
		phoneType := c.DefaultPostForm("phoneType", "work")
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
	if org := c.DefaultPostForm("organization", c.Param("organization")); org != "" {
		card.AddOrganization(org)

		// Add additional organization details if provided
		if dept := c.DefaultPostForm("department", ""); dept != "" {
			card.AddDepartment(dept)
		}
		if title := c.DefaultPostForm("title", ""); title != "" {
			card.AddTitle(title)
		}
		if role := c.DefaultPostForm("role", ""); role != "" {
			card.AddRole(role)
		}
	}

	// URL
	if url := c.DefaultPostForm("url", c.Param("url")); url != "" {
		urlType := c.DefaultPostForm("urlType", "work")
		var uType vcard.URLType
		switch strings.ToLower(urlType) {
		case "home":
			uType = vcard.URLHome
		case "social":
			uType = vcard.URLSocial
		default:
			uType = vcard.URLWork
		}
		card.AddURL(url, uType)
	}

	// Note
	if note := c.DefaultPostForm("note", c.Param("note")); note != "" {
		card.AddNote(note)
	}

	return card
}
