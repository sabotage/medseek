package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"medseek/internal/handlers"
	"medseek/internal/service"
	"medseek/internal/websocket"

	"github.com/joho/godotenv"
)

// corsMiddleware adds CORS headers to support iOS and other browsers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers for iOS compatibility
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// WebSocket-specific headers for better support
		w.Header().Set("Connection", "Upgrade")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// iOS Safari specific headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load environment variables
	godotenv.Load()

	deepseekAPIKey := os.Getenv("DEEPSEEK_API_KEY")
	if deepseekAPIKey == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize services
	chatService := service.NewChatService(deepseekAPIKey)
	wsHub := websocket.NewHub(chatService)

	// Start WebSocket hub
	go wsHub.Run()

	// Initialize handlers
	handler := handlers.NewHandler(chatService, wsHub)

	// Setup routes
	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/api/session/create", handler.CreateSession)
	http.HandleFunc("/api/session/messages", handler.GetSessionMessages)
	http.HandleFunc("/api/session/close", handler.CloseSession)
	http.HandleFunc("/ws", handler.WebSocket)

	// Serve static files from frontend
	// Try multiple possible locations
	staticDirs := []string{
		"./frontend/dist",
		"./dist",
		"/home/ecs-user/medseek-deploy-20251230_224927/frontend/dist",
		"/home/ecs-user/medseek-deploy-20251230_224927/dist",
	}

	var fs http.Handler
	for _, dir := range staticDirs {
		if _, err := os.Stat(dir); err == nil {
			log.Printf("Using static file directory: %s", dir)
			fs = http.FileServer(http.Dir(dir))
			break
		}
	}

	if fs == nil {
		log.Printf("Warning: No static file directory found. Tried: %v", staticDirs)
		// Fallback to current directory
		fs = http.FileServer(http.Dir("."))
	}

	http.Handle("/", fs)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), corsMiddleware(http.DefaultServeMux)))
}
