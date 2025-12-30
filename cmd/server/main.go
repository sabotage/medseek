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
	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", fs)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
