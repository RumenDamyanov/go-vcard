package chi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	vcard "go.rumenx.com/vcard"
)

func TestVCardMiddleware(t *testing.T) {
	r := chi.NewRouter()

	// Handler that creates a simple vCard
	handler := func(w http.ResponseWriter, r *http.Request) *vcard.VCard {
		card := vcard.New()
		card.AddName("John", "Doe")
		card.AddEmail("john@example.com")
		return card
	}

	// Add middleware
	r.Get("/test", VCard(handler))

	// Test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/vcard" {
		t.Errorf("Expected Content-Type text/vcard, got %s", contentType)
	}

	// Check Content-Disposition
	contentDisposition := rr.Header().Get("Content-Disposition")
	if !strings.Contains(contentDisposition, "attachment") {
		t.Errorf("Expected Content-Disposition to contain 'attachment', got %s", contentDisposition)
	}
}

func TestVCardWithCustomOptions(t *testing.T) {
	r := chi.NewRouter()

	// Handler that creates a simple vCard
	handler := func(w http.ResponseWriter, r *http.Request) *vcard.VCard {
		card := vcard.New()
		card.AddName("Jane", "Smith")
		return card
	}

	// Custom options
	options := Options{
		Filename: func(w http.ResponseWriter, r *http.Request) string {
			return "jane-smith.vcf"
		},
		ContentDisposition: "inline",
	}

	// Add middleware with custom options
	r.Get("/test", VCard(handler, options))

	// Test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check custom filename
	contentDisposition := rr.Header().Get("Content-Disposition")
	if !strings.Contains(contentDisposition, "jane-smith.vcf") {
		t.Errorf("Expected filename 'jane-smith.vcf' in Content-Disposition, got %s", contentDisposition)
	}

	if !strings.Contains(contentDisposition, "inline") {
		t.Errorf("Expected 'inline' in Content-Disposition, got %s", contentDisposition)
	}
}

func TestVCardJSONMiddleware(t *testing.T) {
	r := chi.NewRouter()

	// Handler that creates a vCard with data
	handler := func(w http.ResponseWriter, r *http.Request) *vcard.VCard {
		card := vcard.New()
		card.AddName("Jane", "Smith")
		card.AddEmail("jane@example.com")
		return card
	}

	// Add JSON middleware
	r.Get("/vcard", VCardJSON(handler))

	// Test request
	req := httptest.NewRequest("GET", "/vcard", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check Content-Type
	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain application/json, got %s", contentType)
	}

	// Parse and verify the JSON response
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
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
	r := chi.NewRouter()

	// Test handler that uses CreateFromParams
	r.Get("/create", func(w http.ResponseWriter, r *http.Request) {
		card := CreateFromParams(w, r)
		if card.GetFormattedName() == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No name provided"))
			return
		}
		w.Write([]byte("OK: " + card.GetFormattedName()))
	})

	// Test with query parameters
	req := httptest.NewRequest("GET", "/create?firstName=John&lastName=Doe&email=john@example.com&phone=123-456-7890&organization=ACME%20Corp", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "John Doe") {
		t.Errorf("Expected response to contain 'John Doe', got %s", body)
	}
}

func TestCreateFromParamsWithURLParams(t *testing.T) {
	r := chi.NewRouter()

	// Test handler with URL parameters
	r.Get("/user/{firstName}/{lastName}", func(w http.ResponseWriter, r *http.Request) {
		card := CreateFromParams(w, r)
		w.Write([]byte(card.GetFormattedName()))
	})

	// Test with URL parameters
	req := httptest.NewRequest("GET", "/user/John/Doe", nil)

	// Add context with chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("firstName", "John")
	rctx.URLParams.Add("lastName", "Doe")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "John Doe") {
		t.Errorf("Expected response to contain 'John Doe', got %s", body)
	}
}

func TestCreateFromParamsEmailTypes(t *testing.T) {
	r := chi.NewRouter()

	// Test handler for email types
	r.Get("/email", func(w http.ResponseWriter, r *http.Request) {
		card := CreateFromParams(w, r)
		emails := card.GetEmails()
		if len(emails) == 0 {
			w.Write([]byte("No email"))
			return
		}
		w.Write([]byte(string(emails[0].Type)))
	})

	// Test home email type
	req := httptest.NewRequest("GET", "/email?email=test@example.com&emailType=home", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	body := rr.Body.String()
	if body != "HOME" {
		t.Errorf("Expected email type 'HOME', got %s", body)
	}
}

func TestVCardErrorHandling(t *testing.T) {
	r := chi.NewRouter()

	// Handler that returns nil (should cause error)
	handler := func(w http.ResponseWriter, r *http.Request) *vcard.VCard {
		return nil
	}

	r.Get("/error", VCard(handler))

	req := httptest.NewRequest("GET", "/error", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rr.Code)
	}
}

func TestVCardJSONErrorHandling(t *testing.T) {
	r := chi.NewRouter()

	// Handler that returns nil (should cause error)
	handler := func(w http.ResponseWriter, r *http.Request) *vcard.VCard {
		return nil
	}

	r.Get("/error", VCardJSON(handler))

	req := httptest.NewRequest("GET", "/error", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rr.Code)
	}

	// Check that response is JSON
	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected JSON response for error, got Content-Type: %s", contentType)
	}
}
