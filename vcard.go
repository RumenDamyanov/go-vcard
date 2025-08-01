// Package vcard provides functionality for generating vCard files (.vcf)
// compatible with major contact managers (iOS, Android, Gmail, iCloud, etc.).
//
// This package is framework-agnostic and works seamlessly with any Go web
// framework including Gin, Echo, Fiber, Chi, and standard net/http.
//
// Basic usage:
//
//	card := vcard.New()
//	card.AddName("John", "Doe")
//	card.AddEmail("john.doe@example.com")
//	card.AddPhone("+1234567890")
//
//	content, err := card.String()
//	if err != nil {
//		// handle error
//	}
//
//	// Save to file
//	err = card.SaveToFile("contact.vcf")
//
// The package supports both vCard 3.0 and 4.0 specifications and provides
// type-safe methods for adding various contact information including names,
// emails, phones, addresses, organizations, photos, and custom properties.
package vcard

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Version represents a vCard version specification
type Version string

const (
	// Version30 represents vCard version 3.0 (RFC 2426)
	Version30 Version = "3.0"

	// Version40 represents vCard version 4.0 (RFC 6350)
	Version40 Version = "4.0"
)

// String returns the string representation of the version
func (v Version) String() string {
	return string(v)
}

// VCard represents a vCard contact entry with all supported properties
type VCard struct {
	version      Version
	name         Name
	emails       []Email
	phones       []Phone
	addresses    []Address
	organization Organization
	urls         []URL
	photo        string
	note         string
	birthday     *time.Time
	anniversary  *time.Time
	customProps  map[string]string
}

// New creates a new vCard instance with default settings (version 3.0)
func New() *VCard {
	return &VCard{
		version:     Version30,
		emails:      make([]Email, 0),
		phones:      make([]Phone, 0),
		addresses:   make([]Address, 0),
		urls:        make([]URL, 0),
		customProps: make(map[string]string),
	}
}

// NewWithVersion creates a new vCard instance with the specified version
func NewWithVersion(version Version) *VCard {
	card := New()
	card.version = version
	return card
}

// SetVersion sets the vCard version
func (v *VCard) SetVersion(version Version) *VCard {
	v.version = version
	return v
}

// GetVersion returns the current vCard version
func (v *VCard) GetVersion() Version {
	return v.version
}

// String generates the vCard content as a string
func (v *VCard) String() (string, error) {
	if err := v.Validate(); err != nil {
		return "", fmt.Errorf("vcard validation failed: %w", err)
	}

	var builder strings.Builder

	// Begin vCard
	builder.WriteString("BEGIN:VCARD\n")
	builder.WriteString(fmt.Sprintf("VERSION:%s\n", v.version))

	// Add name information
	if err := v.writeNameProperties(&builder); err != nil {
		return "", err
	}

	// Add contact information
	v.writeEmailProperties(&builder)
	v.writePhoneProperties(&builder)
	v.writeAddressProperties(&builder)
	v.writeOrganizationProperties(&builder)
	v.writeURLProperties(&builder)

	// Add optional properties
	if v.photo != "" {
		v.writePhotoProperty(&builder)
	}

	if v.note != "" {
		builder.WriteString(fmt.Sprintf("NOTE:%s\n", escapeValue(v.note)))
	}

	if v.birthday != nil {
		v.writeBirthdayProperty(&builder)
	}

	if v.anniversary != nil {
		v.writeAnniversaryProperty(&builder)
	}

	// Add custom properties
	v.writeCustomProperties(&builder)

	// End vCard
	builder.WriteString("END:VCARD\n")

	return builder.String(), nil
}

// Bytes generates the vCard content as a byte slice
func (v *VCard) Bytes() ([]byte, error) {
	content, err := v.String()
	if err != nil {
		return nil, err
	}
	return []byte(content), nil
}

// SaveToFile saves the vCard content to a file
func (v *VCard) SaveToFile(filename string) error {
	content, err := v.String()
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(content), 0644)
}

// Validate checks if the vCard has required fields and valid data
func (v *VCard) Validate() error {
	// Check if name is provided (required field)
	if v.name.First == "" && v.name.Last == "" {
		return fmt.Errorf("vcard must have at least first name or last name")
	}

	// Validate emails
	for _, email := range v.emails {
		if email.Address == "" {
			return fmt.Errorf("email address cannot be empty")
		}
	}

	// Validate phones
	for _, phone := range v.phones {
		if phone.Number == "" {
			return fmt.Errorf("phone number cannot be empty")
		}
	}

	return nil
}

// IsValid returns true if the vCard has valid required fields
func (v *VCard) IsValid() bool {
	return v.Validate() == nil
}

// Reset clears all vCard data, allowing reuse of the instance
func (v *VCard) Reset() *VCard {
	v.version = Version30
	v.name = Name{}
	v.emails = v.emails[:0]
	v.phones = v.phones[:0]
	v.addresses = v.addresses[:0]
	v.organization = Organization{}
	v.urls = v.urls[:0]
	v.photo = ""
	v.note = ""
	v.birthday = nil
	v.anniversary = nil

	// Clear custom properties map
	for k := range v.customProps {
		delete(v.customProps, k)
	}

	return v
}

// Clone creates a deep copy of the vCard
func (v *VCard) Clone() *VCard {
	clone := &VCard{
		version:      v.version,
		name:         v.name,
		emails:       make([]Email, len(v.emails)),
		phones:       make([]Phone, len(v.phones)),
		addresses:    make([]Address, len(v.addresses)),
		organization: v.organization,
		urls:         make([]URL, len(v.urls)),
		photo:        v.photo,
		note:         v.note,
		customProps:  make(map[string]string),
	}

	// Copy slices
	copy(clone.emails, v.emails)
	copy(clone.phones, v.phones)
	copy(clone.addresses, v.addresses)
	copy(clone.urls, v.urls)

	// Copy time pointers
	if v.birthday != nil {
		birthday := *v.birthday
		clone.birthday = &birthday
	}

	if v.anniversary != nil {
		anniversary := *v.anniversary
		clone.anniversary = &anniversary
	}

	// Copy custom properties
	for k, v := range v.customProps {
		clone.customProps[k] = v
	}

	return clone
}

// GetFormattedName returns the formatted full name
func (v *VCard) GetFormattedName() string {
	return v.name.FormattedName()
}

// GetName returns the name structure
func (v *VCard) GetName() Name {
	return v.name
}

// GetEmails returns all email addresses
func (v *VCard) GetEmails() []Email {
	return v.emails
}

// GetPhones returns all phone numbers
func (v *VCard) GetPhones() []Phone {
	return v.phones
}

// GetAddresses returns all addresses
func (v *VCard) GetAddresses() []Address {
	return v.addresses
}

// GetOrganization returns the organization information
func (v *VCard) GetOrganization() Organization {
	return v.organization
}

// GetURLs returns all URLs
func (v *VCard) GetURLs() []URL {
	return v.urls
}

// GetPhoto returns the photo data/URL
func (v *VCard) GetPhoto() string {
	return v.photo
}

// GetNote returns the note text
func (v *VCard) GetNote() string {
	return v.note
}

// GetBirthday returns the birthday if set
func (v *VCard) GetBirthday() *time.Time {
	return v.birthday
}

// GetAnniversary returns the anniversary if set
func (v *VCard) GetAnniversary() *time.Time {
	return v.anniversary
}

// GetCustomProperties returns all custom properties
func (v *VCard) GetCustomProperties() map[string]string {
	props := make(map[string]string)
	for k, v := range v.customProps {
		props[k] = v
	}
	return props
}

// GetCustomProperty returns a specific custom property value
func (v *VCard) GetCustomProperty(name string) string {
	return v.customProps[name]
}
