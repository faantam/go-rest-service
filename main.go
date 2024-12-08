package main

import (
	"log"
	"net/http"

	"github.com/faantam/go-rest-service/controllers"
)

func main() {
	// Route setup
	http.HandleFunc("/receipts/process", controllers.HandleProcessReceipt)
	http.HandleFunc("/receipts/", controllers.HandleGetPoints)

	// Start the server
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
