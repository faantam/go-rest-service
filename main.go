package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/faantam/go-rest-service/models"
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

	points := calculatePoints(receipt)

	// Respond with points
	response := GetPointsResponse{Points: points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Points calculation logic
func calculatePoints(receipt models.Receipt) int {
	points := 0

	// Rule 1: 1 point for every alphanumeric character in the retailer name
	points += len(regexp.MustCompile(`[a-zA-Z0-9]`).FindAllString(receipt.Retailer, -1))

	// Rule 2: 50 points if the total is a round dollar amount
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if total, err := strconv.ParseFloat(receipt.Total, 64); err == nil && int(total*100)%25 == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: If the trimmed item description length is a multiple of 3, multiply the item price by 0.2
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int(price * 0.2)
			}
		}
	}

	// Rule 6: 6 points if the purchase date is odd
	if date, err := strconv.Atoi(strings.Split(receipt.PurchaseDate, "-")[2]); err == nil && date%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the purchase time is between 2:00 PM and 4:00 PM
	if time, err := strconv.Atoi(strings.Replace(receipt.PurchaseTime, ":", "", 1)); err == nil && time >= 1400 && time < 1600 {
		points += 10
	}

	return points
}

func main() {
	// Route setup
	http.HandleFunc("/receipts/process", handleProcessReceipt)
	http.HandleFunc("/receipts/", handleGetPoints)

	// Start the server
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
