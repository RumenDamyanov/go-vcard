package vcard

import (
	"os"
	"strings"
	"testing"
	"time"
)

// Test additional methods that aren't covered
func TestAdditionalMethods(t *testing.T) {
	card := New()

	// Test middle name, prefix, suffix
	card.AddName("John", "Doe").AddMiddleName("William").AddPrefix("Dr.").AddSuffix("Jr.")

	if card.name.Middle != "William" {
		t.Errorf("Expected middle name William, got %s", card.name.Middle)
	}
	if card.name.Prefix != "Dr." {
		t.Errorf("Expected prefix Dr., got %s", card.name.Prefix)
	}
	if card.name.Suffix != "Jr." {
		t.Errorf("Expected suffix Jr., got %s", card.name.Suffix)
	}

	// Test SetName
	newName := Name{
		First:  "Jane",
		Last:   "Smith",
		Middle: "Marie",
		Prefix: "Ms.",
		Suffix: "PhD",
	}
	card.SetName(newName)
	if card.name.First != "Jane" {
		t.Errorf("Expected first name Jane, got %s", card.name.First)
	}

	// Test email with preference
	card.AddEmailWithPreference("jane@example.com", EmailWork, true)
	if len(card.emails) == 0 || !card.emails[0].Preferred {
		t.Error("Email preference not set correctly")
	}

	// Test AddEmails
	emails := []Email{
		{Address: "email1@example.com", Type: EmailWork},
		{Address: "email2@example.com", Type: EmailHome},
	}
	card.AddEmails(emails)
	if len(card.emails) < 3 {
		t.Error("AddEmails should have added 2 more emails")
	}

	// Test phone with preference
	card.AddPhoneWithPreference("+1234567890", PhoneWork, true)
	if len(card.phones) == 0 || !card.phones[0].Preferred {
		t.Error("Phone preference not set correctly")
	}

	// Test AddPhones
	phones := []Phone{
		{Number: "+1111111111", Type: PhoneHome},
		{Number: "+2222222222", Type: PhoneMobile},
	}
	card.AddPhones(phones)
	if len(card.phones) < 3 {
		t.Error("AddPhones should have added 2 more phones")
	}

	// Test extended address
	card.AddAddressExtended("123 Main St", "Suite 100", "Springfield", "IL", "62701", "USA", AddressWork)
	if len(card.addresses) == 0 {
		t.Error("Extended address not added")
	}

	// Test address with preference
	card.AddAddressWithPreference("456 Oak Ave", "Hometown", "CA", "90210", "USA", AddressHome, true)
	if len(card.addresses) < 2 || !card.addresses[1].Preferred {
		t.Error("Address preference not set correctly")
	}

	// Test AddAddresses
	addresses := []Address{
		{Street: "789 Pine St", City: "Denver", State: "CO", PostalCode: "80202", Country: "USA", Type: AddressWork},
	}
	card.AddAddresses(addresses)
	if len(card.addresses) < 3 {
		t.Error("AddAddresses should have added 1 more address")
	}

	// Test department and role
	card.AddDepartment("Engineering").AddRole("Software Engineer")
	if card.organization.Department != "Engineering" {
		t.Errorf("Expected department Engineering, got %s", card.organization.Department)
	}
	if card.organization.Role != "Software Engineer" {
		t.Errorf("Expected role Software Engineer, got %s", card.organization.Role)
	}

	// Test SetOrganization
	org := Organization{
		Name:       "ACME Corp",
		Department: "IT",
		Title:      "Senior Developer",
		Role:       "Lead",
	}
	card.SetOrganization(org)
	if card.organization.Name != "ACME Corp" {
		t.Errorf("Expected organization ACME Corp, got %s", card.organization.Name)
	}

	// Test URL methods
	card.AddURL("https://example.com", URLWork)
	if len(card.urls) == 0 {
		t.Error("URL not added")
	}

	card.AddURLWithPreference("https://home.example.com", URLHome, true)
	if len(card.urls) < 2 || !card.urls[1].Preferred {
		t.Error("URL preference not set correctly")
	}

	urls := []URL{
		{Address: "https://social.example.com", Type: URLSocial},
	}
	card.AddURLs(urls)
	if len(card.urls) < 3 {
		t.Error("AddURLs should have added 1 more URL")
	}

	// Test photo
	card.AddPhoto("https://example.com/photo.jpg")
	if card.photo != "https://example.com/photo.jpg" {
		t.Errorf("Expected photo URL, got %s", card.photo)
	}

	// Test note
	card.AddNote("This is a test note")
	if card.note != "This is a test note" {
		t.Errorf("Expected note, got %s", card.note)
	}

	// Test birthday from string
	err := card.AddBirthdayFromString("1990-01-15")
	if err != nil {
		t.Errorf("Error adding birthday from string: %v", err)
	}

	// Test anniversary from string
	err = card.AddAnniversaryFromString("2020-06-20")
	if err != nil {
		t.Errorf("Error adding anniversary from string: %v", err)
	}

	// Test custom properties
	card.AddCustomProperties(map[string]string{
		"X-SOCIAL-PROFILE": "linkedin.com/in/johndoe",
		"X-COMPANY-ID":     "12345",
	})
	if len(card.GetCustomProperties()) != 2 {
		t.Error("Custom properties not added correctly")
	}
}

func TestGetterMethods(t *testing.T) {
	card := New()
	card.AddName("John", "Doe")
	card.AddEmail("john@example.com", EmailWork)
	card.AddPhone("+1234567890", PhoneWork)
	card.AddAddress("123 Main St", "Springfield", "IL", "62701", "USA", AddressWork)
	card.AddOrganization("ACME Corp")
	card.AddURL("https://example.com", URLWork)
	card.AddPhoto("https://example.com/photo.jpg")
	card.AddNote("Test note")
	card.AddBirthday(time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC))
	card.AddAnniversary(time.Date(2020, 6, 20, 0, 0, 0, 0, time.UTC))
	card.AddCustomProperty("X-TEST", "value")

	// Test all getter methods
	name := card.GetName()
	if name.FormattedName() == "" {
		t.Error("GetName should return formatted name")
	}

	emails := card.GetEmails()
	if len(emails) == 0 {
		t.Error("GetEmails should return emails")
	}

	phones := card.GetPhones()
	if len(phones) == 0 {
		t.Error("GetPhones should return phones")
	}

	addresses := card.GetAddresses()
	if len(addresses) == 0 {
		t.Error("GetAddresses should return addresses")
	}

	urls := card.GetURLs()
	if len(urls) == 0 {
		t.Error("GetURLs should return URLs")
	}

	photo := card.GetPhoto()
	if photo == "" {
		t.Error("GetPhoto should return photo")
	}

	note := card.GetNote()
	if note == "" {
		t.Error("GetNote should return note")
	}

	birthday := card.GetBirthday()
	if birthday.IsZero() {
		t.Error("GetBirthday should return birthday")
	}

	anniversary := card.GetAnniversary()
	if anniversary.IsZero() {
		t.Error("GetAnniversary should return anniversary")
	}

	customProps := card.GetCustomProperties()
	if len(customProps) == 0 {
		t.Error("GetCustomProperties should return custom properties")
	}
}

func TestFileMethods(t *testing.T) {
	card := New()
	card.AddName("John", "Doe")
	card.AddEmail("john@example.com")

	// Test Bytes method
	data, err := card.Bytes()
	if err != nil {
		t.Errorf("Bytes() error: %v", err)
	}
	if len(data) == 0 {
		t.Error("Bytes() should return data")
	}

	// Test SaveToFile method
	filename := "test_contact.vcf"
	defer os.Remove(filename) // Clean up

	err = card.SaveToFile(filename)
	if err != nil {
		t.Errorf("SaveToFile() error: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("SaveToFile() should create file")
	}
}

func TestVersionMethods(t *testing.T) {
	// Test NewWithVersion
	card := NewWithVersion(Version40)
	if card.GetVersion() != Version40 {
		t.Errorf("Expected version %s, got %s", Version40, card.GetVersion())
	}

	// Test SetVersion
	card.SetVersion(Version30)
	if card.GetVersion() != Version30 {
		t.Errorf("Expected version %s, got %s", Version30, card.GetVersion())
	}
}

func TestValidationMethods(t *testing.T) {
	card := New()

	// Test IsValid method
	if card.IsValid() {
		t.Error("Empty card should not be valid")
	}

	card.AddName("John", "Doe")
	if !card.IsValid() {
		t.Error("Card with name should be valid")
	}
}

func TestCloneMethod(t *testing.T) {
	original := New()
	original.AddName("John", "Doe")
	original.AddEmail("john@example.com")
	original.AddCustomProperty("X-TEST", "value")

	// Test Clone
	cloned := original.Clone()
	if cloned.GetFormattedName() != original.GetFormattedName() {
		t.Error("Cloned card should have same formatted name")
	}

	// Modify original to ensure independence
	original.AddName("Jane", "Smith")
	if cloned.GetFormattedName() == original.GetFormattedName() {
		t.Error("Cloned card should be independent of original")
	}
}

func TestContactMethod(t *testing.T) {
	card := New()

	contact := Contact{
		Name: Name{
			First: "John",
			Last:  "Doe",
		},
		Emails: []Email{
			{Address: "john@example.com", Type: EmailWork},
		},
		Phones: []Phone{
			{Number: "+1234567890", Type: PhoneWork},
		},
		Organization: Organization{
			Name: "ACME Corp",
		},
	}

	card.AddContact(contact)

	if card.GetFormattedName() != "John Doe" {
		t.Error("Contact should set name correctly")
	}
	if len(card.GetEmails()) == 0 {
		t.Error("Contact should set emails correctly")
	}
	if len(card.GetPhones()) == 0 {
		t.Error("Contact should set phones correctly")
	}
	orgName := card.GetOrganization()
	if orgName.Name != "ACME Corp" {
		t.Error("Contact should set organization correctly")
	}
}

func TestIndividualGetters(t *testing.T) {
	card := New()
	card.AddEmail("john@example.com", EmailWork)
	card.AddEmail("john@home.com", EmailHome)
	card.AddPhone("+1234567890", PhoneWork)
	card.AddPhone("+9876543210", PhoneHome)
	card.AddAddress("123 Main St", "Springfield", "IL", "62701", "USA", AddressWork)
	card.AddURL("https://example.com", URLWork)

	// Test GetEmail - just get first email
	workEmail := card.GetEmail()
	if workEmail == "" {
		t.Error("GetEmail should return first email")
	}

	// Test GetPhone - just get first phone
	workPhone := card.GetPhone()
	if workPhone == "" {
		t.Error("GetPhone should return first phone")
	}

	// Test GetAddress - just get first address
	address := card.GetAddress()
	if address == nil {
		t.Error("GetAddress should return first address")
	}

	// Test GetURL - just get first URL
	url := card.GetURL()
	if url == "" {
		t.Error("GetURL should return first URL")
	}
}

func TestPhotoFromFile(t *testing.T) {
	card := New()

	// Test with non-existent file (should return error)
	err := card.AddPhotoFromFile("non-existent-file.jpg")
	if err == nil {
		t.Error("AddPhotoFromFile should return error for non-existent file")
	}

	// Create a temporary test file
	tempFile := "test-photo.jpg"
	testData := []byte("fake image data")
	err = os.WriteFile(tempFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tempFile)

	// Test with existing file
	err = card.AddPhotoFromFile(tempFile)
	if err != nil {
		t.Errorf("AddPhotoFromFile should not return error for existing file: %v", err)
	}
}

func TestUnescapeValue(t *testing.T) {
	// Test the unescapeValue function through edge cases
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello\\,World", "Hello,World"},
		{"Line1\\nLine2", "Line1\nLine2"},
		{"Tab\\tSeparated", "Tab\tSeparated"},
		{"Back\\\\slash", "Back\\slash"},
		{"Normal text", "Normal text"},
	}

	for _, tc := range testCases {
		result := unescapeValue(tc.input)
		if result != tc.expected {
			t.Errorf("unescapeValue(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestFoldLine(t *testing.T) {
	// Test line folding for very long lines
	longLine := strings.Repeat("A", 100) // Create a line longer than 75 characters
	folded := foldLine(longLine)

	// Should contain line breaks for folding
	if !strings.Contains(folded, "\r\n ") {
		t.Error("Long line should be folded with CRLF + space")
	}

	// Test short line (should not be folded)
	shortLine := "Short"
	unfolded := foldLine(shortLine)
	if unfolded != shortLine {
		t.Error("Short line should not be folded")
	}
}

// Test additional methods that aren't covered
