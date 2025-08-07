package main

import (
	"encoding/json"
	"fmt"
	"go-url-shortener/internal/dependencies"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storages"
	"net/http"
)

func main() {
	db, err := storages.NewDB()
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	// Initialize Redis client
	redisClient, err := storages.NewRedis()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}
	defer storages.CloseRedis()
	fmt.Println("Starting HTTP server on port 8080...")

	app := &dependencies.ShortenerApp{
		DB:    db,
		Redis: redisClient,
	}

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("web")))

	// Handle frontend API endpoint
	http.HandleFunc("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			URL string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		handlers.CreateShortURLHandler(w, r, app)
	})

	// Handle shortener endpoint
	http.HandleFunc("/v1/shortener", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateShortURLHandler(w, r, app)
	})

	// Handle redirect endpoint
	http.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		handlers.RedirectHandler(w, r, app)
	})

	// Handle get all links endpoint
	http.HandleFunc("/v1/links", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllLinksHandler(w, r, app)
	})

	// Start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}
