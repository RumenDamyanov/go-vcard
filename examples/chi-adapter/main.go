package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	vcard "go.rumenx.com/vcard"
)

// VCardResponse represents the JSON response structure
type VCardResponse struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Phone        string `json:"phone,omitempty"`
	Organization string `json:"organization,omitempty"`
	Title        string `json:"title,omitempty"`
}

// VCardMiddleware is a Chi middleware that adds vCard creation functionality
func VCardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CreateVCardFromParams creates a vCard from URL parameters
func CreateVCardFromParams(r *http.Request) (*vcard.VCard, error) {
	query := r.URL.Query()

	firstName := query.Get("firstName")
	lastName := query.Get("lastName")
	email := query.Get("email")
	phone := query.Get("phone")
	organization := query.Get("organization")
	title := query.Get("title")

	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("firstName and lastName are required")
	}

	vc := vcard.New()
	vc.AddName(firstName, lastName)

	if email != "" {
		vc.AddEmail(email, vcard.EmailWork)
	}
	if phone != "" {
		vc.AddPhone(phone, vcard.PhoneWork)
	}
	if organization != "" {
		vc.AddOrganization(organization)
	}
	if title != "" {
		vc.AddTitle(title)
	}

	return vc, nil
}

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(VCardMiddleware)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Go-VCard Chi Example</title>
</head>
<body>
    <h1>Go-VCard Chi Framework Example</h1>
    <h2>Available Endpoints:</h2>
    <ul>
        <li><a href="/vcard/John/Doe?email=john@example.com">/vcard/{firstName}/{lastName}</a> - Download vCard file</li>
        <li><a href="/contact-json?firstName=Jane&lastName=Smith&email=jane@example.com">/contact-json</a> - Get JSON response</li>
        <li><a href="/health">/health</a> - Health check</li>
    </ul>
    <h3>Example URLs:</h3>
    <ul>
        <li><a href="/vcard/John/Doe?email=john@example.com&phone=+1234567890">/vcard/John/Doe?email=john@example.com&phone=+1234567890</a></li>
        <li><a href="/contact-json?firstName=Jane&lastName=Smith&email=jane@example.com&organization=Acme Corp&title=Developer">/contact-json?firstName=Jane&lastName=Smith&email=jane@example.com&organization=Acme Corp&title=Developer</a></li>
    </ul>
</body>
</html>`
		w.Write([]byte(html))
	})

	// vCard download endpoint with path parameters
	r.Get("/vcard/{firstName}/{lastName}", func(w http.ResponseWriter, r *http.Request) {
		firstName := chi.URLParam(r, "firstName")
		lastName := chi.URLParam(r, "lastName")

		// Create a new request with the path parameters as query parameters
		query := r.URL.Query()
		query.Set("firstName", firstName)
		query.Set("lastName", lastName)
		r.URL.RawQuery = query.Encode()

		vc, err := CreateVCardFromParams(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		vCardData, err := vc.String()
		if err != nil {
			http.Error(w, "Failed to generate vCard", http.StatusInternalServerError)
			return
		}
		filename := fmt.Sprintf("%s_%s.vcf", firstName, lastName)

		w.Header().Set("Content-Type", "text/vcard")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Write([]byte(vCardData))
	})

	// JSON response endpoint
	r.Get("/contact-json", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		response := VCardResponse{
			FirstName:    query.Get("firstName"),
			LastName:     query.Get("lastName"),
			Email:        query.Get("email"),
			Phone:        query.Get("phone"),
			Organization: query.Get("organization"),
			Title:        query.Get("title"),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "healthy",
			"framework": "chi",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// POST endpoint for creating vCards with JSON body
	r.Post("/vcard", func(w http.ResponseWriter, r *http.Request) {
		var req VCardResponse
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.FirstName == "" || req.LastName == "" {
			http.Error(w, "firstName and lastName are required", http.StatusBadRequest)
			return
		}

		vc := vcard.New()
		vc.AddName(req.FirstName, req.LastName)

		if req.Email != "" {
			vc.AddEmail(req.Email, vcard.EmailWork)
		}
		if req.Phone != "" {
			vc.AddPhone(req.Phone, vcard.PhoneWork)
		}
		if req.Organization != "" {
			vc.AddOrganization(req.Organization)
		}
		if req.Title != "" {
			vc.AddTitle(req.Title)
		}

		vCardData, err := vc.String()
		if err != nil {
			http.Error(w, "Failed to generate vCard", http.StatusInternalServerError)
			return
		}
		filename := fmt.Sprintf("%s_%s.vcf", req.FirstName, req.LastName)

		w.Header().Set("Content-Type", "text/vcard")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Write([]byte(vCardData))
	})

	port := ":8083"
	fmt.Printf("Chi server starting on port %s\n", port)
	fmt.Printf("Visit http://localhost%s for available endpoints\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
