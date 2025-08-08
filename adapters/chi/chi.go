// Package chi provides Chi framework adapter for go-vcard
package chi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.rumenx.com/vcard"
)

// VCardHandler is a function that returns a VCard
type VCardHandler func(w http.ResponseWriter, r *http.Request) *vcard.VCard

// Options configures the vCard response
type Options struct {
	// Filename generates the filename for the vCard download
	Filename func(w http.ResponseWriter, r *http.Request) string

	// ContentDisposition sets how the file should be handled (attachment/inline)
	ContentDisposition string
}

// DefaultOptions provides sensible defaults
var DefaultOptions = Options{
	Filename: func(w http.ResponseWriter, r *http.Request) string {
		return "contact.vcf"
	},
	ContentDisposition: "attachment",
}

// VCard middleware for Chi that generates vCard responses
func VCard(handler VCardHandler, opts ...Options) http.HandlerFunc {
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

	return func(w http.ResponseWriter, r *http.Request) {
		// Generate vCard
		card := handler(w, r)
		if card == nil {
			http.Error(w, "Failed to generate vCard", http.StatusInternalServerError)
			return
		}

		// Generate vCard content
		content, err := card.String()
		if err != nil {
			http.Error(w, "Failed to generate vCard content", http.StatusInternalServerError)
			return
		}

		// Set headers
		filename := options.Filename(w, r)
		w.Header().Set("Content-Type", "text/vcard")
		w.Header().Set("Content-Disposition", options.ContentDisposition+"; filename="+filename)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(content))
	}
}

// VCardJSON middleware for Chi that returns vCard data as JSON
func VCardJSON(handler VCardHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate vCard
		card := handler(w, r)
		if card == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to generate vCard",
			})
			return
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// CreateFromParams creates a vCard from Chi context parameters and query values
func CreateFromParams(w http.ResponseWriter, r *http.Request) *vcard.VCard {
	card := vcard.New()

	// Name from path parameters or query parameters
	if firstName := r.URL.Query().Get("firstName"); firstName == "" {
		firstName = chi.URLParam(r, "firstName")
		if firstName != "" {
			lastName := r.URL.Query().Get("lastName")
			if lastName == "" {
				lastName = chi.URLParam(r, "lastName")
			}
			card.AddName(firstName, lastName)
		}
	} else {
		lastName := r.URL.Query().Get("lastName")
		card.AddName(firstName, lastName)
	}

	// Email
	if email := r.URL.Query().Get("email"); email != "" {
		emailType := r.URL.Query().Get("emailType")
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
	if phone := r.URL.Query().Get("phone"); phone != "" {
		phoneType := r.URL.Query().Get("phoneType")
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
	if org := r.URL.Query().Get("organization"); org != "" {
		card.AddOrganization(org)
	}

	// Title
	if title := r.URL.Query().Get("title"); title != "" {
		card.AddTitle(title)
	}

	// URL
	if url := r.URL.Query().Get("url"); url != "" {
		card.AddURL(url, vcard.URLWork)
	}

	// Note
	if note := r.URL.Query().Get("note"); note != "" {
		card.AddNote(note)
	}

	return card
}
