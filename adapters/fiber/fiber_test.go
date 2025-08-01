package fiber

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	vcard "github.com/rumendamyanov/go-vcard"
)

func TestVCardMiddleware(t *testing.T) {
	app := fiber.New()

	// Handler that creates a simple vCard
	handler := func(c *fiber.Ctx) *vcard.VCard {
		card := vcard.New()
		card.AddName("John", "Doe")
		card.AddEmail("john@example.com")
		return card
	}

	// Add middleware
	app.Get("/test", VCard(handler))

	// Test request
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/vcard" {
		t.Errorf("Expected Content-Type text/vcard, got %s", contentType)
	}

	// Check Content-Disposition
	contentDisposition := resp.Header.Get("Content-Disposition")
	if !strings.Contains(contentDisposition, "attachment") {
		t.Errorf("Expected Content-Disposition to contain 'attachment', got %s", contentDisposition)
	}
}

func TestVCardWithCustomOptions(t *testing.T) {
	app := fiber.New()

	// Handler that creates a simple vCard
	handler := func(c *fiber.Ctx) *vcard.VCard {
		card := vcard.New()
		card.AddName("Jane", "Smith")
		return card
	}

	// Custom options
	options := Options{
		Filename: func(c *fiber.Ctx) string {
			return "jane-smith.vcf"
		},
		ContentDisposition: "inline",
	}

	// Add middleware with custom options
	app.Get("/test", VCard(handler, options))

	// Test request
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// Check custom filename
	contentDisposition := resp.Header.Get("Content-Disposition")
	if !strings.Contains(contentDisposition, "jane-smith.vcf") {
		t.Errorf("Expected filename 'jane-smith.vcf' in Content-Disposition, got %s", contentDisposition)
	}

	if !strings.Contains(contentDisposition, "inline") {
		t.Errorf("Expected 'inline' in Content-Disposition, got %s", contentDisposition)
	}
}

func TestVCardJSONMiddleware(t *testing.T) {
	app := fiber.New()

	// Handler that creates a vCard with data
	handler := func(c *fiber.Ctx) *vcard.VCard {
		card := vcard.New()
		card.AddName("Jane", "Smith")
		card.AddEmail("jane@example.com")
		return card
	}

	// Add JSON middleware
	app.Get("/vcard", VCardJSON(handler))

	// Test request
	req := httptest.NewRequest("GET", "/vcard", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain application/json, got %s", contentType)
	}

	// Parse and verify the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if response["name"] == nil {
		t.Error("Response should contain name field")
	}

	if response["emails"] == nil {
		t.Error("Response should contain emails field")
	}
}

func TestCreateFromParams(t *testing.T) {
	app := fiber.New()

	// Test handler that uses CreateFromParams
	app.Get("/create", func(c *fiber.Ctx) error {
		card := CreateFromParams(c)
		if card.GetFormattedName() == "" {
			return c.SendString("No name provided")
		}
		return c.SendString("OK: " + card.GetFormattedName())
	})

	// Test with query parameters
	req := httptest.NewRequest("GET", "/create?firstName=John&lastName=Doe&email=john@example.com&phone=123-456-7890&organization=ACME%20Corp", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	bodyStr := string(body)
	if !strings.Contains(bodyStr, "John Doe") {
		t.Errorf("Expected response to contain 'John Doe', got %s", bodyStr)
	}
}

func TestCreateFromParamsEmailTypes(t *testing.T) {
	app := fiber.New()

	// Test handler for email types
	app.Get("/email", func(c *fiber.Ctx) error {
		card := CreateFromParams(c)
		emails := card.GetEmails()
		if len(emails) == 0 {
			return c.SendString("No email")
		}
		return c.SendString(string(emails[0].Type))
	})

	// Test home email type
	req := httptest.NewRequest("GET", "/email?email=test@example.com&emailType=home", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "HOME" {
		t.Errorf("Expected email type 'HOME', got %s", string(body))
	}
}

func TestVCardErrorHandling(t *testing.T) {
	app := fiber.New()

	// Handler that returns nil (should cause error)
	handler := func(c *fiber.Ctx) *vcard.VCard {
		return nil
	}

	app.Get("/error", VCard(handler))

	req := httptest.NewRequest("GET", "/error", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}
}
