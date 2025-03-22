package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/addspin/acomm/pkg/models"
)

const (
	dbPath = "./db/comm.db"
	addr   = ":8082"
)

type CommentRequest struct {
	Text string `json:"text"`
}

// handleGetComments retrieves all comments for a news ID
func handleGetComments(w http.ResponseWriter, r *http.Request) {
	// Extract news ID from query parameter
	newsIDStr := r.URL.Query().Get("id")
	if newsIDStr == "" {
		http.Error(w, "Missing news ID", http.StatusBadRequest)
		return
	}

	newsID, err := strconv.ParseInt(newsIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid news ID", http.StatusBadRequest)
		return
	}

	comments, err := models.GetCommentsByNewsID(newsID)
	if err != nil {
		log.Printf("Error getting comments: %v", err)
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// handleAddComment adds a new comment for a news ID
func handleAddComment(w http.ResponseWriter, r *http.Request) {
	// Extract news ID from query parameter
	newsIDStr := r.URL.Query().Get("id")
	if newsIDStr == "" {
		http.Error(w, "Missing news ID", http.StatusBadRequest)
		return
	}

	newsID, err := strconv.ParseInt(newsIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid news ID", http.StatusBadRequest)
		return
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, "Comment text cannot be empty", http.StatusBadRequest)
		return
	}

	commentID, err := models.AddComment(newsID, req.Text)
	if err != nil {
		log.Printf("Error adding comment: %v", err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID int64 `json:"id"`
	}{
		ID: commentID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// setupRoutes configures the HTTP routes
func setupRoutes() {
	// Endpoint for getting comments for a specific news ID
	http.HandleFunc("/api/comm_news", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handleGetComments(w, r)
	})

	// Endpoint for adding comments for a specific news ID
	http.HandleFunc("/api/comm_add_news", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handleAddComment(w, r)
	})
}

func main() {
	// Initialize database
	if err := models.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Setup routes
	setupRoutes()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-c
	log.Println("Shutting down server...")
}
