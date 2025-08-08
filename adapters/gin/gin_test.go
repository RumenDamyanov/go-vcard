package gin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.rumenx.com/vcard"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestVCard(t *testing.T) {
	// Create a test handler that returns a vCard
	handler := func(c *gin.Context) *vcard.VCard {
		card := vcard.New()
		card.AddName("John", "Doe")
		card.AddEmail("john@example.com", vcard.EmailWork)
		return card
	}

	// Create Gin engine and test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/", nil)
	c.Request = req

	// Test the VCard middleware
	vcardHandler := VCard(handler)
	vcardHandler(c)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Header().Get("Content-Type"), "text/vcard") {
		t.Errorf("Expected Content-Type to contain 'text/vcard', got %s", w.Header().Get("Content-Type"))
	}

	if !strings.Contains(w.Header().Get("Content-Disposition"), "attachment") {
		t.Errorf("Expected Content-Disposition to contain 'attachment', got %s", w.Header().Get("Content-Disposition"))
	}

	body := w.Body.String()
	if !strings.Contains(body, "BEGIN:VCARD") {
		t.Error("Expected vCard content to contain 'BEGIN:VCARD'")
	}

	if !strings.Contains(body, "john@example.com") {
		t.Error("Expected vCard content to contain email")
	}
}

func TestVCardNilHandler(t *testing.T) {
	// Create a test handler that returns nil
	handler := func(c *gin.Context) *vcard.VCard {
		return nil
	}

	// Create Gin engine and test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/", nil)
	c.Request = req

	// Test the VCard middleware
	vcardHandler := VCard(handler)
	vcardHandler(c)

	// Assertions
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}
