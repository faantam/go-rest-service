package services

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/faantam/go-rest-service/models"
)

func CalculatePoints(receipt models.Receipt) int {
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
