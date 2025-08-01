package vcard

import (
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	card := New()

	if card == nil {
		t.Fatal("New() returned nil")
	}

	if card.GetVersion() != Version30 {
		t.Errorf("Expected version %s, got %s", Version30, card.GetVersion())
	}
}

func TestNewWithVersion(t *testing.T) {
	card := NewWithVersion(Version40)

	if card.GetVersion() != Version40 {
		t.Errorf("Expected version %s, got %s", Version40, card.GetVersion())
	}
}

func TestBasicVCard(t *testing.T) {
	card := New()
	card.AddName("John", "Doe")
	card.AddEmail("john.doe@example.com")
	card.AddPhone("+1234567890")

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	// Check required components
	if !strings.Contains(content, "BEGIN:VCARD") {
		t.Error("vCard missing BEGIN:VCARD")
	}

	if !strings.Contains(content, "END:VCARD") {
		t.Error("vCard missing END:VCARD")
	}

	if !strings.Contains(content, "VERSION:3.0") {
		t.Error("vCard missing VERSION:3.0")
	}

	if !strings.Contains(content, "FN:John Doe") {
		t.Error("vCard missing formatted name")
	}

	if !strings.Contains(content, "N:Doe;John;;;") {
		t.Error("vCard missing structured name")
	}

	if !strings.Contains(content, "EMAIL;TYPE=INTERNET:john.doe@example.com") {
		t.Error("vCard missing email")
	}

	if !strings.Contains(content, "TEL;TYPE=VOICE:+1234567890") {
		t.Error("vCard missing phone")
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*VCard)
		wantErr bool
	}{
		{
			name: "valid card with first name",
			setup: func(card *VCard) {
				card.AddName("John", "")
			},
			wantErr: false,
		},
		{
			name: "valid card with last name",
			setup: func(card *VCard) {
				card.AddName("", "Doe")
			},
			wantErr: false,
		},
		{
			name: "valid card with both names",
			setup: func(card *VCard) {
				card.AddName("John", "Doe")
			},
			wantErr: false,
		},
		{
			name: "invalid card with no name",
			setup: func(card *VCard) {
				// Don't add any name
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := New()
			tt.setup(card)

			err := card.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMethodChaining(t *testing.T) {
	card := New().
		AddName("John", "Doe").
		AddEmail("john@work.com", EmailWork).
		AddPhone("+1234567890", PhoneWork).
		AddOrganization("Acme Corp").
		AddTitle("Developer")

	if card.GetFormattedName() != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", card.GetFormattedName())
	}

	if card.GetEmail() != "john@work.com" {
		t.Errorf("Expected 'john@work.com', got '%s'", card.GetEmail())
	}

	if card.GetPhone() != "+1234567890" {
		t.Errorf("Expected '+1234567890', got '%s'", card.GetPhone())
	}

	org := card.GetOrganization()
	if org.Name != "Acme Corp" {
		t.Errorf("Expected 'Acme Corp', got '%s'", org.Name)
	}

	if org.Title != "Developer" {
		t.Errorf("Expected 'Developer', got '%s'", org.Title)
	}
}

func TestEmailTypes(t *testing.T) {
	card := New()
	card.AddName("John", "Doe")
	card.AddEmail("work@example.com", EmailWork)
	card.AddEmail("home@example.com", EmailHome)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "EMAIL;TYPE=WORK:work@example.com") {
		t.Error("Work email not properly formatted")
	}

	if !strings.Contains(content, "EMAIL;TYPE=HOME:home@example.com") {
		t.Error("Home email not properly formatted")
	}
}

func TestPhoneTypes(t *testing.T) {
	card := New()
	card.AddName("John", "Doe")
	card.AddPhone("+1234567890", PhoneWork)
	card.AddPhone("+1987654321", PhoneMobile)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "TEL;TYPE=WORK:+1234567890") {
		t.Error("Work phone not properly formatted")
	}

	if !strings.Contains(content, "TEL;TYPE=MOBILE:+1987654321") {
		t.Error("Mobile phone not properly formatted")
	}
}

func TestAddress(t *testing.T) {
	card := New()
	card.AddName("John", "Doe")
	card.AddAddress("123 Main St", "Anytown", "CA", "12345", "USA", AddressWork)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "ADR;TYPE=WORK:;;123 Main St;Anytown;CA;12345;USA") {
		t.Error("Address not properly formatted")
	}
}

func TestCustomProperties(t *testing.T) {
	card := New()
	card.AddName("John", "Doe")
	card.AddCustomProperty("X-DEPARTMENT", "Engineering")
	card.AddCustomProperty("X-EMPLOYEE-ID", "EMP001")

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "X-DEPARTMENT:Engineering") {
		t.Error("Custom property X-DEPARTMENT not found")
	}

	if !strings.Contains(content, "X-EMPLOYEE-ID:EMP001") {
		t.Error("Custom property X-EMPLOYEE-ID not found")
	}
}

func TestBirthday(t *testing.T) {
	card := New()
	card.AddName("John", "Doe")

	birthday := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	card.AddBirthday(birthday)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "BDAY:1990-05-15") {
		t.Error("Birthday not properly formatted")
	}
}

func TestVersion40Features(t *testing.T) {
	card := NewWithVersion(Version40)
	card.AddName("John", "Doe")

	anniversary := time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)
	card.AddAnniversary(anniversary)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "VERSION:4.0") {
		t.Error("vCard 4.0 version not found")
	}

	if !strings.Contains(content, "ANNIVERSARY:2020-06-01") {
		t.Error("Anniversary not found (vCard 4.0 feature)")
	}
}

func TestClone(t *testing.T) {
	original := New()
	original.AddName("John", "Doe")
	original.AddEmail("john@example.com")
	original.AddCustomProperty("X-TEST", "value")

	clone := original.Clone()

	// Modify original
	original.AddName("Jane", "Smith")
	original.AddCustomProperty("X-TEST", "modified")

	// Check that clone wasn't affected
	if clone.GetFormattedName() != "John Doe" {
		t.Error("Clone was affected by original modification")
	}

	if clone.GetCustomProperty("X-TEST") != "value" {
		t.Error("Clone custom property was affected by original modification")
	}
}

func TestReset(t *testing.T) {
	card := New()
	card.AddName("John", "Doe")
	card.AddEmail("john@example.com")
	card.AddCustomProperty("X-TEST", "value")

	card.Reset()

	if card.GetFormattedName() != "" {
		t.Error("Reset did not clear name")
	}

	if card.GetEmail() != "" {
		t.Error("Reset did not clear email")
	}

	if card.GetCustomProperty("X-TEST") != "" {
		t.Error("Reset did not clear custom properties")
	}
}

func TestEscapeValue(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal text", "normal text"},
		{"text with, comma", "text with\\, comma"},
		{"text with; semicolon", "text with\\; semicolon"},
		{"text with\nnewline", "text with\\nnewline"},
		{"text with\\backslash", "text with\\\\backslash"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeValue(tt.input)
			if result != tt.expected {
				t.Errorf("escapeValue(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNameFormatting(t *testing.T) {
	tests := []struct {
		name     Name
		expected string
	}{
		{
			Name{First: "John", Last: "Doe"},
			"John Doe",
		},
		{
			Name{Prefix: "Dr.", First: "John", Last: "Doe", Suffix: "Jr."},
			"Dr. John Doe Jr.",
		},
		{
			Name{First: "John", Middle: "William", Last: "Doe"},
			"John William Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.name.FormattedName()
			if result != tt.expected {
				t.Errorf("FormattedName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test additional methods for coverage improvement
func TestAdditionalCoverage(t *testing.T) {
	card := New()

	// Test middle name, prefix, suffix
	card.AddName("John", "Doe").AddMiddleName("William").AddPrefix("Dr.").AddSuffix("Jr.")

	// Test SetName
	newName := Name{
		First:  "Jane",
		Last:   "Smith",
		Middle: "Marie",
		Prefix: "Ms.",
		Suffix: "PhD",
	}
	card.SetName(newName)

	// Test email with preference
	card.AddEmailWithPreference("jane@example.com", EmailWork, true)

	// Test AddEmails
	emails := []Email{
		{Address: "email1@example.com", Type: EmailWork},
		{Address: "email2@example.com", Type: EmailHome},
	}
	card.AddEmails(emails)

	// Test phone with preference
	card.AddPhoneWithPreference("+1234567890", PhoneWork, true)

	// Test AddPhones
	phones := []Phone{
		{Number: "+1111111111", Type: PhoneHome},
		{Number: "+2222222222", Type: PhoneMobile},
	}
	card.AddPhones(phones)

	// Test extended address
	card.AddAddressExtended("123 Main St", "Suite 100", "Springfield", "IL", "62701", "USA", AddressWork)

	// Test address with preference
	card.AddAddressWithPreference("456 Oak Ave", "Hometown", "CA", "90210", "USA", AddressHome, true)

	// Test AddAddresses
	addresses := []Address{
		{Street: "789 Pine St", City: "Denver", State: "CO", PostalCode: "80202", Country: "USA", Type: AddressWork},
	}
	card.AddAddresses(addresses)

	// Test department and role
	card.AddDepartment("Engineering").AddRole("Software Engineer")

	// Test SetOrganization
	org := Organization{
		Name:       "ACME Corp",
		Department: "IT",
		Title:      "Senior Developer",
		Role:       "Lead",
	}
	card.SetOrganization(org)

	// Test URL methods
	card.AddURL("https://example.com", URLWork)
	card.AddURLWithPreference("https://home.example.com", URLHome, true)
	urls := []URL{
		{Address: "https://social.example.com", Type: URLSocial},
	}
	card.AddURLs(urls)

	// Test photo
	card.AddPhoto("https://example.com/photo.jpg")

	// Test note
	card.AddNote("This is a test note")

	// Test birthday from string
	_ = card.AddBirthdayFromString("1990-01-15")

	// Test anniversary from string
	_ = card.AddAnniversaryFromString("2020-06-20")

	// Test custom properties
	card.AddCustomProperties(map[string]string{
		"X-SOCIAL-PROFILE": "linkedin.com/in/johndoe",
		"X-COMPANY-ID":     "12345",
	})

	// Test utility functions for coverage
	result := unescapeValue("Hello\\,World")
	if result != "Hello,World" {
		t.Errorf("unescapeValue failed: got %s", result)
	}

	// Test foldLine with long line
	longLine := strings.Repeat("A", 100)
	_ = foldLine(longLine)

	// Test foldLine with short line
	_ = foldLine("Short")

	// Test getter methods for coverage
	_ = card.GetName()
	_ = card.GetEmails()
	_ = card.GetPhones()
	_ = card.GetAddresses()
	_ = card.GetURLs()
	_ = card.GetPhoto()
	_ = card.GetNote()
	_ = card.GetBirthday()
	_ = card.GetAnniversary()
	_ = card.GetCustomProperties()

	// Test individual getters
	card2 := New()
	card2.AddEmail("test@example.com")
	card2.AddPhone("+1234567890")
	card2.AddAddress("123 Main St", "City", "State", "12345", "Country", AddressWork)
	card2.AddURL("https://example.com", URLWork)

	_ = card2.GetEmail()
	_ = card2.GetPhone()
	_ = card2.GetAddress()
	_ = card2.GetURL()

	// Test file methods
	data, _ := card.Bytes()
	if len(data) == 0 {
		t.Error("Bytes should return data")
	}

	// Test version methods
	card3 := NewWithVersion(Version40)
	_ = card3.GetVersion()
	card3.SetVersion(Version30)
	_ = card3.GetVersion()

	// Test validation methods
	_ = card.IsValid()

	// Test photo from file (with error)
	_ = card.AddPhotoFromFile("non-existent.jpg")
}
