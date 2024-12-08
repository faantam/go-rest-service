package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/faantam/go-rest-service/models"
	"github.com/faantam/go-rest-service/services"
	"github.com/google/uuid"
)

// ProcessReceiptResponse represents the response containing the receipt ID
type ProcessReceiptResponse struct {
	ID string `json:"id"`
}

// GetPointsResponse represents the response for points
type GetPointsResponse struct {
	Points int `json:"points"`
}

func handleProcessReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Generate a UUID for the receipt
	id := strings.ToLower(uuid.New().String())

	// Use AddReceipt to add the receipt
	err := models.AddReceipt(id, receipt)
	if err != nil {
		http.Error(w, "UUID collision detected, receipt not processed", http.StatusConflict)
		return
	}

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from the URL
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	receipt, exists := models.GetReceipt(id)

	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	points := services.CalculatePoints(receipt)

	response := GetPointsResponse{Points: points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Route setup
	http.HandleFunc("/receipts/process", handleProcessReceipt)
	http.HandleFunc("/receipts/", handleGetPoints)

	// Start the server
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
