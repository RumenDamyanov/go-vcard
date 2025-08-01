package echo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rumendamyanov/go-vcard"
)

func TestVCard(t *testing.T) {
	// Create a test handler that returns a vCard
	handler := func(c echo.Context) *vcard.VCard {
		card := vcard.New()
		card.AddName("John", "Doe")
		card.AddEmail("john@example.com", vcard.EmailWork)
		return card
	}

	// Create Echo instance and test request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test the VCard middleware
	vcardHandler := VCard(handler)
	err := vcardHandler(c)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Header().Get("Content-Type"), "text/vcard") {
		t.Errorf("Expected Content-Type to contain 'text/vcard', got %s", rec.Header().Get("Content-Type"))
	}

	if !strings.Contains(rec.Header().Get("Content-Disposition"), "attachment") {
		t.Errorf("Expected Content-Disposition to contain 'attachment', got %s", rec.Header().Get("Content-Disposition"))
	}

	body := rec.Body.String()
	if !strings.Contains(body, "BEGIN:VCARD") {
		t.Error("Expected vCard content to contain 'BEGIN:VCARD'")
	}

	if !strings.Contains(body, "FN:John Doe") {
		t.Error("Expected vCard content to contain formatted name")
	}

	if !strings.Contains(body, "john@example.com") {
		t.Error("Expected vCard content to contain email")
	}
}

func TestVCardJSON(t *testing.T) {
	// Create a test handler that returns a vCard
	handler := func(c echo.Context) *vcard.VCard {
		card := vcard.New()
		card.AddName("Jane", "Smith")
		card.AddEmail("jane@example.com", vcard.EmailHome)
		card.AddPhone("+1234567890", vcard.PhoneMobile)
		return card
	}

	// Create Echo instance and test request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test the VCardJSON middleware
	jsonHandler := VCardJSON(handler)
	err := jsonHandler(c)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
		t.Errorf("Expected Content-Type to contain 'application/json', got %s", rec.Header().Get("Content-Type"))
	}

	body := rec.Body.String()
	if !strings.Contains(body, "jane@example.com") {
		t.Error("Expected JSON response to contain email")
	}
}

func TestCreateFromParams(t *testing.T) {
	// Create Echo instance and test request with query parameters
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?firstName=Test&lastName=User&email=test@example.com&phone=123-456-7890&organization=Test+Corp&title=Developer", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test CreateFromParams
	card := CreateFromParams(c)

	// Assertions
	if card == nil {
		t.Fatal("Expected card to be created, got nil")
	}

	// Verify the card contains expected data
	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard content: %v", err)
	}

	if !strings.Contains(content, "Test User") {
		t.Error("Expected vCard to contain name 'Test User'")
	}

	if !strings.Contains(content, "test@example.com") {
		t.Error("Expected vCard to contain email")
	}

	if !strings.Contains(content, "123-456-7890") {
		t.Error("Expected vCard to contain phone number")
	}

	if !strings.Contains(content, "Test Corp") {
		t.Error("Expected vCard to contain organization")
	}

	if !strings.Contains(content, "Developer") {
		t.Error("Expected vCard to contain title")
	}
}

func TestCreateFromParamsWithPathParams(t *testing.T) {
	// Create Echo instance and test request with path parameters
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/contact/John/Doe?email=john@example.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set path parameters
	c.SetParamNames("firstName", "lastName")
	c.SetParamValues("John", "Doe")

	// Test CreateFromParams
	card := CreateFromParams(c)

	// Assertions
	if card == nil {
		t.Fatal("Expected card to be created, got nil")
	}

	// Verify the card contains expected data
	content, err := card.String()
	if err != nil {
		t.Fatalf("Failed to generate vCard content: %v", err)
	}

	if !strings.Contains(content, "John Doe") {
		t.Error("Expected vCard to contain name 'John Doe'")
	}

	if !strings.Contains(content, "john@example.com") {
		t.Error("Expected vCard to contain email")
	}
}

func TestVCardWithCustomOptions(t *testing.T) {
	// Create a test handler that returns a vCard
	handler := func(c echo.Context) *vcard.VCard {
		card := vcard.New()
		card.AddName("Custom", "User")
		return card
	}

	// Custom options
	options := Options{
		Filename: func(c echo.Context) string {
			return "custom-contact.vcf"
		},
		ContentDisposition: "inline",
	}

	// Create Echo instance and test request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test the VCard middleware with custom options
	vcardHandler := VCard(handler, options)
	err := vcardHandler(c)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	disposition := rec.Header().Get("Content-Disposition")
	if !strings.Contains(disposition, "inline") {
		t.Errorf("Expected Content-Disposition to contain 'inline', got %s", disposition)
	}

	if !strings.Contains(disposition, "custom-contact.vcf") {
		t.Errorf("Expected Content-Disposition to contain custom filename, got %s", disposition)
	}
}

func TestVCardNilHandler(t *testing.T) {
	// Create a test handler that returns nil
	handler := func(c echo.Context) *vcard.VCard {
		return nil
	}

	// Create Echo instance and test request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test the VCard middleware
	vcardHandler := VCard(handler)
	err := vcardHandler(c)

	// Assertions
	if err == nil {
		t.Error("Expected error when handler returns nil, got nil")
	}

	echoErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Errorf("Expected echo.HTTPError, got %T", err)
	}

	if echoErr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", echoErr.Code)
	}
}
