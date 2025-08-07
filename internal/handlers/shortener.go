package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-url-shortener/internal/dependencies"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// ShortenerRequest represents the request body for creating a short URL
// swagger:parameters createShortUrl
// in: body
// required: true
// schema:
//
//	$ref: "#/definitions/ShortenerRequest"
type ShortenerRequest struct {
	URL string `json:"url"`
}

// ShortenerResponse represents the response body for creating a short URL
// swagger:response createShortUrlResponse
type Link struct {
	RedirectCode string    `json:"redirect_code"`
	DestinyURL   string    `json:"destiny_url"`
	ShortURL     string    `json:"short_url"`
	Clicks       int       `json:"clicks"`
	CreatedAt    time.Time `json:"created_at"`
}

type LinksResponse struct {
	Links []Link `json:"links"`
}

type ShortenerResponse struct {
	ShortURL string `json:"short_url"`
}

// GetAllLinksHandler returns all links from the database
func GetAllLinksHandler(w http.ResponseWriter, r *http.Request, app *dependencies.ShortenerApp) {
	rows, err := app.DB.Query("SELECT redirect_code, destiny_url, clicks, created_at FROM links ORDER BY created_at DESC")
	if err != nil {
		log.Printf("Error querying links: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		if err := rows.Scan(&link.RedirectCode, &link.DestinyURL, &link.Clicks, &link.CreatedAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Construct the short URL
		link.ShortURL = fmt.Sprintf("http://%s/r/%s", r.Host, link.RedirectCode)
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := LinksResponse{Links: links}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RedirectHandler handles URL redirection
func RedirectHandler(w http.ResponseWriter, r *http.Request, app *dependencies.ShortenerApp) {
	ctx := r.Context()
	redirectCode := strings.TrimPrefix(r.URL.Path, "/r/")

	if redirectCode == r.URL.Path {
		http.Error(w, "Invalid redirect code", http.StatusBadRequest)
		return
	}

	var destinyURL string

	// try to get the short URL from Redis
	destinyURL, err := app.Redis.Get(ctx, redirectCode).Result()

	if err != nil && err == redis.Nil {
		// Query the database for the redirect
		row := app.DB.QueryRow("SELECT destiny_url FROM links WHERE redirect_code = $1", redirectCode)

		if err := row.Scan(&destinyURL); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Redirect code not found", http.StatusNotFound)
				return
			}
			log.Printf("Error querying database: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set the short URL in Redis
		go func(ctx context.Context) {
			if err := app.Redis.Set(ctx, redirectCode, destinyURL, 24*time.Hour).Err(); err != nil {
				log.Printf("Error setting short URL in Redis: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}(context.WithoutCancel(ctx))
	}

	// Increment clicks in a goroutine (call and forget)
	go func() {
		_, err := app.DB.Exec("UPDATE links SET clicks = clicks + 1 WHERE redirect_code = $1", redirectCode)
		if err != nil {
			log.Printf("Error updating clicks: %v", err)
		}
	}()

	// Redirect to the destination URL
	http.Redirect(w, r, destinyURL, http.StatusTemporaryRedirect)
}

// CreateShortURLHandler handles POST requests to create a new short URL
func CreateShortURLHandler(w http.ResponseWriter, r *http.Request, app *dependencies.ShortenerApp) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate URL
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Generate redirect code (simple implementation using random string)
	redirectCode := generateRedirectCode()

	// Insert into database
	query := "INSERT INTO links (redirect_code, destiny_url) VALUES ($1, $2)"
	_, err := app.DB.Exec(query, redirectCode, req.URL)
	if err != nil {
		log.Printf("Error inserting into database: %v", err)
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	// Create response
	response := ShortenerResponse{
		ShortURL: fmt.Sprintf("http://%s/r/%s", r.Host, redirectCode),
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to create response", http.StatusInternalServerError)
	}
}

// generateRedirectCode generates a random string for redirect code
func generateRedirectCode() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	for i := 0; i < 6; i++ {
		result.WriteByte(chars[rand.Intn(len(chars))])
	}
	return result.String()
}
