package vcard

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestSaveToFile(t *testing.T) {
	card := New()
	card.AddName("Test", "User")

	filename := "test_save.vcf"
	defer os.Remove(filename) // Clean up after test

	err := card.SaveToFile(filename)
	if err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("File was not created")
	}

	// Read content
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "FN:Test User") {
		t.Error("File content is incorrect")
	}
}

func TestBytes(t *testing.T) {
	card := New()
	card.AddName("Test", "User")

	bytes, err := card.Bytes()
	if err != nil {
		t.Fatalf("Bytes() failed: %v", err)
	}

	if len(bytes) == 0 {
		t.Error("Bytes() returned empty slice")
	}

	content := string(bytes)
	if !strings.Contains(content, "FN:Test User") {
		t.Error("Bytes content is incorrect")
	}
}

func TestAddEmailWithPreference(t *testing.T) {
	card := New()
	card.AddName("Test", "User")
	card.AddEmailWithPreference("preferred@example.com", EmailWork, true)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "EMAIL;TYPE=WORK;PREF=1:preferred@example.com") {
		t.Error("Preferred email not properly formatted")
	}
}

func TestAddPhoneWithPreference(t *testing.T) {
	card := New()
	card.AddName("Test", "User")
	card.AddPhoneWithPreference("+1234567890", PhoneMobile, true)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "TEL;TYPE=MOBILE;PREF=1:+1234567890") {
		t.Error("Preferred phone not properly formatted")
	}
}

func TestAddAddressExtended(t *testing.T) {
	card := New()
	card.AddName("Test", "User")
	card.AddAddressExtended("123 Main St", "Suite 100", "City", "State", "12345", "Country", AddressWork)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "ADR;TYPE=WORK:;Suite 100;123 Main St;City;State;12345;Country") {
		t.Error("Extended address not properly formatted")
	}
}

func TestAddAddressWithPreference(t *testing.T) {
	card := New()
	card.AddName("Test", "User")
	card.AddAddressWithPreference("123 Main St", "City", "State", "12345", "Country", AddressHome, true)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "ADR;TYPE=HOME;PREF=1:;;123 Main St;City;State;12345;Country") {
		t.Error("Preferred address not properly formatted")
	}
}

func TestBatchOperations(t *testing.T) {
	card := New()
	card.AddName("Test", "User")

	// Test AddEmails
	emails := []Email{
		{Address: "email1@example.com", Type: EmailWork},
		{Address: "email2@example.com", Type: EmailHome},
	}
	card.AddEmails(emails)

	// Test AddPhones
	phones := []Phone{
		{Number: "+1111111111", Type: PhoneWork},
		{Number: "+2222222222", Type: PhoneMobile},
	}
	card.AddPhones(phones)

	// Test AddAddresses
	addresses := []Address{
		{Street: "Work St", City: "Work City", Type: AddressWork},
		{Street: "Home St", City: "Home City", Type: AddressHome},
	}
	card.AddAddresses(addresses)

	// Test AddURLs
	urls := []URL{
		{Address: "https://work.example.com", Type: URLWork},
		{Address: "https://personal.example.com", Type: URLHome},
	}
	card.AddURLs(urls)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	// Check that all items were added
	if !strings.Contains(content, "email1@example.com") {
		t.Error("First email not found")
	}
	if !strings.Contains(content, "email2@example.com") {
		t.Error("Second email not found")
	}
	if !strings.Contains(content, "+1111111111") {
		t.Error("First phone not found")
	}
	if !strings.Contains(content, "+2222222222") {
		t.Error("Second phone not found")
	}
}

func TestOrganizationMethods(t *testing.T) {
	card := New()
	card.AddName("Test", "User")
	card.AddOrganization("Test Corp")
	card.AddDepartment("Engineering")
	card.AddTitle("Developer")
	card.AddRole("Senior")

	org := card.GetOrganization()
	if org.Name != "Test Corp" {
		t.Errorf("Expected 'Test Corp', got '%s'", org.Name)
	}
	if org.Department != "Engineering" {
		t.Errorf("Expected 'Engineering', got '%s'", org.Department)
	}
	if org.Title != "Developer" {
		t.Errorf("Expected 'Developer', got '%s'", org.Title)
	}
	if org.Role != "Senior" {
		t.Errorf("Expected 'Senior', got '%s'", org.Role)
	}

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "ORG:Test Corp;Engineering") {
		t.Error("Organization not properly formatted")
	}
	if !strings.Contains(content, "TITLE:Developer") {
		t.Error("Title not found")
	}
	if !strings.Contains(content, "ROLE:Senior") {
		t.Error("Role not found")
	}
}

func TestURLWithPreference(t *testing.T) {
	card := New()
	card.AddName("Test", "User")
	card.AddURLWithPreference("https://example.com", URLWork, true)

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "URL;TYPE=WORK;PREF=1:https://example.com") {
		t.Error("Preferred URL not properly formatted")
	}
}

func TestPhoto(t *testing.T) {
	card := New()
	card.AddName("Test", "User")

	// Test URL photo
	card.AddPhoto("https://example.com/photo.jpg")

	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "PHOTO;VALUE=uri:https://example.com/photo.jpg") {
		t.Error("Photo URL not properly formatted")
	}

	// Test base64 photo
	card.Reset()
	card.AddName("Test", "User")
	card.AddPhoto("data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD")

	content, err = card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard: %v", err)
	}

	if !strings.Contains(content, "PHOTO;ENCODING=b:data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD") {
		t.Error("Photo base64 not properly formatted")
	}
}

func TestBirthdayFromString(t *testing.T) {
	card := New()
	card.AddName("Test", "User")

	err := card.AddBirthdayFromString("1990-05-15")
	if err != nil {
		t.Fatalf("Failed to add birthday from string: %v", err)
	}

	birthday := card.GetBirthday()
	if birthday == nil {
		t.Fatal("Birthday not set")
	}

	expected := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	if !birthday.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, *birthday)
	}

	// Test invalid date
	err = card.AddBirthdayFromString("invalid-date")
	if err == nil {
		t.Error("Expected error for invalid date")
	}
}

func TestAnniversaryFromString(t *testing.T) {
	card := NewWithVersion(Version40)
	card.AddName("Test", "User")

	err := card.AddAnniversaryFromString("2020-06-01")
	if err != nil {
		t.Fatalf("Failed to add anniversary from string: %v", err)
	}

	anniversary := card.GetAnniversary()
	if anniversary == nil {
		t.Fatal("Anniversary not set")
	}

	expected := time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)
	if !anniversary.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, *anniversary)
	}

	// Test invalid date
	err = card.AddAnniversaryFromString("invalid-date")
	if err == nil {
		t.Error("Expected error for invalid date")
	}
}

func TestAddCustomProperties(t *testing.T) {
	card := New()
	card.AddName("Test", "User")

	props := map[string]string{
		"X-DEPARTMENT":  "Engineering",
		"X-EMPLOYEE-ID": "EMP001",
	}
	card.AddCustomProperties(props)

	if card.GetCustomProperty("X-DEPARTMENT") != "Engineering" {
		t.Error("Custom property not set correctly")
	}

	if card.GetCustomProperty("X-EMPLOYEE-ID") != "EMP001" {
		t.Error("Custom property not set correctly")
	}

	allProps := card.GetCustomProperties()
	if len(allProps) != 2 {
		t.Errorf("Expected 2 custom properties, got %d", len(allProps))
	}
}

func TestContactStructure(t *testing.T) {
	birthday := "1990-05-15"
	anniversary := "2020-06-01"

	contact := Contact{
		Name: Name{
			First:  "John",
			Last:   "Doe",
			Middle: "William",
		},
		Emails: []Email{
			{Address: "john@work.com", Type: EmailWork},
			{Address: "john@home.com", Type: EmailHome},
		},
		Phones: []Phone{
			{Number: "+1234567890", Type: PhoneWork},
			{Number: "+1987654321", Type: PhoneMobile},
		},
		Addresses: []Address{
			{Street: "123 Work St", City: "Work City", Type: AddressWork},
		},
		Organization: Organization{
			Name:       "Test Corp",
			Department: "Engineering",
			Title:      "Developer",
			Role:       "Senior",
		},
		URLs: []URL{
			{Address: "https://johndoe.dev", Type: URLHome},
		},
		Photo:       "https://example.com/photo.jpg",
		Note:        "Test contact",
		Birthday:    &birthday,
		Anniversary: &anniversary,
		CustomProps: map[string]string{
			"X-TEST": "value",
		},
	}

	card := NewWithVersion(Version40)
	card.AddContact(contact)

	// Verify all data was added
	if card.GetFormattedName() != "John William Doe" {
		t.Error("Name not set correctly from contact")
	}

	if len(card.GetEmails()) != 2 {
		t.Error("Emails not set correctly from contact")
	}

	if len(card.GetPhones()) != 2 {
		t.Error("Phones not set correctly from contact")
	}

	if card.GetPhoto() != "https://example.com/photo.jpg" {
		t.Error("Photo not set correctly from contact")
	}

	if card.GetNote() != "Test contact" {
		t.Error("Note not set correctly from contact")
	}

	if card.GetCustomProperty("X-TEST") != "value" {
		t.Error("Custom property not set correctly from contact")
	}
}

func TestGetters(t *testing.T) {
	card := New()
	card.AddName("Test", "User")
	card.AddEmail("test@example.com")
	card.AddPhone("+1234567890")
	card.AddAddress("123 Main St", "City", "State", "12345", "Country")
	card.AddURL("https://example.com")

	// Test getters
	if card.GetEmail() != "test@example.com" {
		t.Error("GetEmail() returned wrong value")
	}

	if card.GetPhone() != "+1234567890" {
		t.Error("GetPhone() returned wrong value")
	}

	addr := card.GetAddress()
	if addr == nil || addr.Street != "123 Main St" {
		t.Error("GetAddress() returned wrong value")
	}

	if card.GetURL() != "https://example.com" {
		t.Error("GetURL() returned wrong value")
	}

	// Test empty card
	emptyCard := New()
	emptyCard.AddName("Empty", "User")

	if emptyCard.GetEmail() != "" {
		t.Error("GetEmail() should return empty string for empty card")
	}

	if emptyCard.GetPhone() != "" {
		t.Error("GetPhone() should return empty string for empty card")
	}

	if emptyCard.GetAddress() != nil {
		t.Error("GetAddress() should return nil for empty card")
	}

	if emptyCard.GetURL() != "" {
		t.Error("GetURL() should return empty string for empty card")
	}
}
