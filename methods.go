package vcard

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"
)

// AddName sets the contact's name
func (v *VCard) AddName(first, last string) *VCard {
	v.name.First = first
	v.name.Last = last
	return v
}

// AddMiddleName sets the middle name
func (v *VCard) AddMiddleName(middle string) *VCard {
	v.name.Middle = middle
	return v
}

// AddPrefix sets the name prefix (Mr., Dr., etc.)
func (v *VCard) AddPrefix(prefix string) *VCard {
	v.name.Prefix = prefix
	return v
}

// AddSuffix sets the name suffix (Jr., PhD, etc.)
func (v *VCard) AddSuffix(suffix string) *VCard {
	v.name.Suffix = suffix
	return v
}

// SetName sets the complete name structure
func (v *VCard) SetName(name Name) *VCard {
	v.name = name
	return v
}

// AddEmail adds an email address with optional type
func (v *VCard) AddEmail(address string, emailType ...EmailType) *VCard {
	email := Email{
		Address: address,
	}

	if len(emailType) > 0 {
		email.Type = emailType[0]
	} else {
		email.Type = EmailInternet
	}

	v.emails = append(v.emails, email)
	return v
}

// AddEmailWithPreference adds an email address with type and preference
func (v *VCard) AddEmailWithPreference(address string, emailType EmailType, preferred bool) *VCard {
	email := Email{
		Address:   address,
		Type:      emailType,
		Preferred: preferred,
	}

	v.emails = append(v.emails, email)
	return v
}

// AddEmails adds multiple email addresses at once
func (v *VCard) AddEmails(emails []Email) *VCard {
	v.emails = append(v.emails, emails...)
	return v
}

// AddPhone adds a phone number with optional type
func (v *VCard) AddPhone(number string, phoneType ...PhoneType) *VCard {
	phone := Phone{
		Number: number,
	}

	if len(phoneType) > 0 {
		phone.Type = phoneType[0]
	} else {
		phone.Type = PhoneVoice
	}

	v.phones = append(v.phones, phone)
	return v
}

// AddPhoneWithPreference adds a phone number with type and preference
func (v *VCard) AddPhoneWithPreference(number string, phoneType PhoneType, preferred bool) *VCard {
	phone := Phone{
		Number:    number,
		Type:      phoneType,
		Preferred: preferred,
	}

	v.phones = append(v.phones, phone)
	return v
}

// AddPhones adds multiple phone numbers at once
func (v *VCard) AddPhones(phones []Phone) *VCard {
	v.phones = append(v.phones, phones...)
	return v
}

// AddAddress adds an address with optional type
func (v *VCard) AddAddress(street, city, state, postalCode, country string, addressType ...AddressType) *VCard {
	address := Address{
		Street:     street,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
	}

	if len(addressType) > 0 {
		address.Type = addressType[0]
	}

	v.addresses = append(v.addresses, address)
	return v
}

// AddAddressExtended adds an address with extended information
func (v *VCard) AddAddressExtended(street, extended, city, state, postalCode, country string, addressType ...AddressType) *VCard {
	address := Address{
		Street:     street,
		Extended:   extended,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
	}

	if len(addressType) > 0 {
		address.Type = addressType[0]
	}

	v.addresses = append(v.addresses, address)
	return v
}

// AddAddressWithPreference adds an address with type and preference
func (v *VCard) AddAddressWithPreference(street, city, state, postalCode, country string, addressType AddressType, preferred bool) *VCard {
	address := Address{
		Street:     street,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
		Type:       addressType,
		Preferred:  preferred,
	}

	v.addresses = append(v.addresses, address)
	return v
}

// AddAddresses adds multiple addresses at once
func (v *VCard) AddAddresses(addresses []Address) *VCard {
	v.addresses = append(v.addresses, addresses...)
	return v
}

// AddOrganization sets the organization name
func (v *VCard) AddOrganization(name string) *VCard {
	v.organization.Name = name
	return v
}

// AddDepartment sets the department
func (v *VCard) AddDepartment(department string) *VCard {
	v.organization.Department = department
	return v
}

// AddTitle sets the job title
func (v *VCard) AddTitle(title string) *VCard {
	v.organization.Title = title
	return v
}

// AddRole sets the role/position
func (v *VCard) AddRole(role string) *VCard {
	v.organization.Role = role
	return v
}

// SetOrganization sets the complete organization structure
func (v *VCard) SetOrganization(org Organization) *VCard {
	v.organization = org
	return v
}

// AddURL adds a URL with optional type
func (v *VCard) AddURL(address string, urlType ...URLType) *VCard {
	url := URL{
		Address: address,
	}

	if len(urlType) > 0 {
		url.Type = urlType[0]
	}

	v.urls = append(v.urls, url)
	return v
}

// AddURLWithPreference adds a URL with type and preference
func (v *VCard) AddURLWithPreference(address string, urlType URLType, preferred bool) *VCard {
	url := URL{
		Address:   address,
		Type:      urlType,
		Preferred: preferred,
	}

	v.urls = append(v.urls, url)
	return v
}

// AddURLs adds multiple URLs at once
func (v *VCard) AddURLs(urls []URL) *VCard {
	v.urls = append(v.urls, urls...)
	return v
}

// AddPhoto sets the photo (URL or base64 data)
func (v *VCard) AddPhoto(photo string) *VCard {
	v.photo = photo
	return v
}

// AddPhotoFromFile loads a photo from file and encodes as base64
func (v *VCard) AddPhotoFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Encode as base64 data URI
	encoded := base64.StdEncoding.EncodeToString(data)
	v.photo = "data:image/jpeg;base64," + encoded
	return nil
}

// AddNote sets a note
func (v *VCard) AddNote(note string) *VCard {
	v.note = note
	return v
}

// AddBirthday sets the birthday
func (v *VCard) AddBirthday(birthday time.Time) *VCard {
	v.birthday = &birthday
	return v
}

// AddBirthdayFromString sets the birthday from a date string (YYYY-MM-DD)
func (v *VCard) AddBirthdayFromString(dateStr string) error {
	birthday, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	v.birthday = &birthday
	return nil
}

// AddAnniversary sets the anniversary (vCard 4.0 only)
func (v *VCard) AddAnniversary(anniversary time.Time) *VCard {
	v.anniversary = &anniversary
	return v
}

// AddAnniversaryFromString sets the anniversary from a date string (YYYY-MM-DD)
func (v *VCard) AddAnniversaryFromString(dateStr string) error {
	anniversary, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	v.anniversary = &anniversary
	return nil
}

// AddCustomProperty adds a custom X- property
func (v *VCard) AddCustomProperty(name, value string) *VCard {
	if v.customProps == nil {
		v.customProps = make(map[string]string)
	}
	v.customProps[name] = value
	return v
}

// AddCustomProperties adds multiple custom properties at once
func (v *VCard) AddCustomProperties(props map[string]string) *VCard {
	if v.customProps == nil {
		v.customProps = make(map[string]string)
	}

	for k, val := range props {
		v.customProps[k] = val
	}

	return v
}

// AddContact adds contact information from a Contact structure
func (v *VCard) AddContact(contact Contact) *VCard {
	// Set name
	v.SetName(contact.Name)

	// Add emails
	if len(contact.Emails) > 0 {
		v.AddEmails(contact.Emails)
	}

	// Add phones
	if len(contact.Phones) > 0 {
		v.AddPhones(contact.Phones)
	}

	// Add addresses
	if len(contact.Addresses) > 0 {
		v.AddAddresses(contact.Addresses)
	}

	// Set organization
	if contact.Organization.Name != "" {
		v.SetOrganization(contact.Organization)
	}

	// Add URLs
	if len(contact.URLs) > 0 {
		v.AddURLs(contact.URLs)
	}

	// Set photo
	if contact.Photo != "" {
		v.AddPhoto(contact.Photo)
	}

	// Set note
	if contact.Note != "" {
		v.AddNote(contact.Note)
	}

	// Set birthday
	if contact.Birthday != nil {
		if err := v.AddBirthdayFromString(*contact.Birthday); err == nil {
			// Birthday set successfully
		}
	}

	// Set anniversary
	if contact.Anniversary != nil {
		if err := v.AddAnniversaryFromString(*contact.Anniversary); err == nil {
			// Anniversary set successfully
		}
	}

	// Add custom properties
	if len(contact.CustomProps) > 0 {
		v.AddCustomProperties(contact.CustomProps)
	}

	return v
}

// GetEmail returns the first email address (if any)
func (v *VCard) GetEmail() string {
	if len(v.emails) > 0 {
		return v.emails[0].Address
	}
	return ""
}

// GetPhone returns the first phone number (if any)
func (v *VCard) GetPhone() string {
	if len(v.phones) > 0 {
		return v.phones[0].Number
	}
	return ""
}

// GetAddress returns the first address (if any)
func (v *VCard) GetAddress() *Address {
	if len(v.addresses) > 0 {
		return &v.addresses[0]
	}
	return nil
}

// GetURL returns the first URL (if any)
func (v *VCard) GetURL() string {
	if len(v.urls) > 0 {
		return v.urls[0].Address
	}
	return ""
}
