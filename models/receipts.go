package models

import (
	"errors"
	"sync"
)

// Receipt represents a receipt with all its details
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Item represents a single item in a receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var (
	receipts = make(map[string]Receipt)
	mutex    sync.Mutex
)

func AddReceipt(id string, receipt Receipt) error {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := receipts[id]; exists {
		return errors.New("receipt with this ID already exists")
	}

	receipts[id] = receipt
	return nil
}

func GetReceipt(id string) (Receipt, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	receipt, exists := receipts[id]
	return receipt, exists
}
