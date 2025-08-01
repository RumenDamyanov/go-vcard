package vcard

import (
	"strings"
)

// EmailType represents the type of email address
type EmailType string

const (
	// EmailInternet represents an internet email address (default)
	EmailInternet EmailType = "INTERNET"

	// EmailWork represents a work email address
	EmailWork EmailType = "WORK"

	// EmailHome represents a home email address
	EmailHome EmailType = "HOME"

	// EmailMobile represents a mobile email address
	EmailMobile EmailType = "MOBILE"
)

// PhoneType represents the type of phone number
type PhoneType string

const (
	// PhoneVoice represents a voice phone number (default)
	PhoneVoice PhoneType = "VOICE"

	// PhoneWork represents a work phone number
	PhoneWork PhoneType = "WORK"

	// PhoneHome represents a home phone number
	PhoneHome PhoneType = "HOME"

	// PhoneMobile represents a mobile phone number
	PhoneMobile PhoneType = "MOBILE"

	// PhoneFax represents a fax number
	PhoneFax PhoneType = "FAX"
)

// AddressType represents the type of address
type AddressType string

const (
	// AddressWork represents a work address
	AddressWork AddressType = "WORK"

	// AddressHome represents a home address
	AddressHome AddressType = "HOME"

	// AddressPostal represents a postal address
	AddressPostal AddressType = "POSTAL"
)

// URLType represents the type of URL
type URLType string

const (
	// URLWork represents a work-related URL
	URLWork URLType = "WORK"

	// URLHome represents a personal URL
	URLHome URLType = "HOME"

	// URLSocial represents a social media URL
	URLSocial URLType = "SOCIAL"
)

// Name represents the structured name information
type Name struct {
	// Last name (family name)
	Last string

	// First name (given name)
	First string

	// Middle name(s) (additional names)
	Middle string

	// Name prefix (Mr., Dr., etc.)
	Prefix string

	// Name suffix (Jr., PhD, etc.)
	Suffix string
}

// FormattedName returns the full formatted name
func (n Name) FormattedName() string {
	var parts []string

	if n.Prefix != "" {
		parts = append(parts, n.Prefix)
	}

	if n.First != "" {
		parts = append(parts, n.First)
	}

	if n.Middle != "" {
		parts = append(parts, n.Middle)
	}

	if n.Last != "" {
		parts = append(parts, n.Last)
	}

	if n.Suffix != "" {
		parts = append(parts, n.Suffix)
	}

	return strings.Join(parts, " ")
}

// StructuredName returns the vCard structured name format (N property)
func (n Name) StructuredName() string {
	return strings.Join([]string{
		escapeValue(n.Last),
		escapeValue(n.First),
		escapeValue(n.Middle),
		escapeValue(n.Prefix),
		escapeValue(n.Suffix),
	}, ";")
}

// Email represents an email address with optional type
type Email struct {
	// The email address
	Address string

	// The type of email (optional)
	Type EmailType

	// Whether this is the preferred email
	Preferred bool
}

// Phone represents a phone number with optional type
type Phone struct {
	// The phone number
	Number string

	// The type of phone (optional)
	Type PhoneType

	// Whether this is the preferred phone
	Preferred bool
}

// Address represents a postal address
type Address struct {
	// Street address
	Street string

	// Extended address (apartment, suite, etc.)
	Extended string

	// City/locality
	City string

	// State/region/province
	State string

	// Postal code
	PostalCode string

	// Country
	Country string

	// Address type (optional)
	Type AddressType

	// Whether this is the preferred address
	Preferred bool
}

// StructuredAddress returns the vCard structured address format (ADR property)
func (a Address) StructuredAddress() string {
	return strings.Join([]string{
		"",                        // Post office box (not commonly used)
		escapeValue(a.Extended),   // Extended address
		escapeValue(a.Street),     // Street address
		escapeValue(a.City),       // Locality
		escapeValue(a.State),      // Region
		escapeValue(a.PostalCode), // Postal code
		escapeValue(a.Country),    // Country
	}, ";")
}

// FormattedAddress returns a human-readable address string
func (a Address) FormattedAddress() string {
	var parts []string

	if a.Street != "" {
		parts = append(parts, a.Street)
	}

	if a.Extended != "" {
		parts = append(parts, a.Extended)
	}

	var cityState []string
	if a.City != "" {
		cityState = append(cityState, a.City)
	}
	if a.State != "" {
		cityState = append(cityState, a.State)
	}
	if len(cityState) > 0 {
		parts = append(parts, strings.Join(cityState, ", "))
	}

	if a.PostalCode != "" {
		parts = append(parts, a.PostalCode)
	}

	if a.Country != "" {
		parts = append(parts, a.Country)
	}

	return strings.Join(parts, "\n")
}

// Organization represents organization/work information
type Organization struct {
	// Organization name
	Name string

	// Department
	Department string

	// Job title
	Title string

	// Role/position
	Role string
}

// URL represents a website or URL with optional type
type URL struct {
	// The URL
	Address string

	// The type of URL (optional)
	Type URLType

	// Whether this is the preferred URL
	Preferred bool
}

// Contact represents a complete contact structure for batch operations
type Contact struct {
	Name         Name
	Emails       []Email
	Phones       []Phone
	Addresses    []Address
	Organization Organization
	URLs         []URL
	Photo        string
	Note         string
	Birthday     *string // Date string in YYYY-MM-DD format
	Anniversary  *string // Date string in YYYY-MM-DD format
	CustomProps  map[string]string
}
