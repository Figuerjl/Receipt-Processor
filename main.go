package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/receipts/process", submitReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getReceiptPointsHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

var receipts []Receipt

func submitReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt

	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the receipt
	if err := validateReceipt(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the received receipt and calculate points
	points := calculatePointsForReceipt(&receipt)

	// Generate a unique ID for the receipt
	receipt.ID = generateID()

	// Append the receipt to the receipts array
	receipts = append(receipts, receipt)

	// Return the ID and points of the processed receipt
	response := map[string]interface{}{
		"id":     receipt.ID,
		"points": points,
	}

	// Set the response headers and encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func calculatePointsForReceipt(receipt *Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
	points += len(regexp.MustCompile("[a-zA-Z0-9]").FindAllString(receipt.Retailer, -1))

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	totalFloat, _ := strconv.ParseFloat(receipt.Total, 64)
	if totalFloat == math.Trunc(totalFloat) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt.
	itemCount := len(receipt.Items)
	points += (itemCount / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	// The result is the number of points earned.
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			priceFloat, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(priceFloat * 0.2))
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	after2PM := time.Date(0, 0, 0, 14, 0, 0, 0, time.UTC)
	before4PM := time.Date(0, 0, 0, 16, 0, 0, 0, time.UTC)
	if purchaseTime.After(after2PM) && purchaseTime.Before(before4PM) {
		points += 10
	}

	return points
}

func isRoundDollarAmount(amount string) bool {
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return false
	}
	return amountFloat == math.Trunc(amountFloat)
}

func isMultipleOfQuarter(amount string) bool {
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return false
	}
	return math.Mod(amountFloat, 0.25) == 0
}

func getReceiptPointsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// Lookup the receipt points from the database or storage
	receipt, err := findReceiptByID(id)

	if err != nil {
		http.Error(w, "Invalid receipt payload", http.StatusBadRequest)
		return
	}
	// Perform the calculation based on the receipt data
	points := calculatePointsForReceipt(receipt)

	if err != nil {
		http.Error(w, "Invalid point payload", http.StatusBadRequest)
		return
	}
	// Return the points for the receipt
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": points})
}

func validateReceipt(receipt *Receipt) error {
	// Validate retailer
	if !regexp.MustCompile(`^\S+$`).MatchString(receipt.Retailer) {
		return fmt.Errorf("Retailer name is invalid")
	}

	// Validate purchase date
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(receipt.PurchaseDate) {
		return fmt.Errorf("Purchase date is invalid")
	}

	// Validate purchase time
	if !regexp.MustCompile(`^\d{2}:\d{2}$`).MatchString(receipt.PurchaseTime) {
		return fmt.Errorf("Purchase time is invalid")
	}

	// Validate total amount
	if !regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(receipt.Total) {
		return fmt.Errorf("Total amount is invalid")
	}

	// Validate items
	if len(receipt.Items) < 1 {
		return fmt.Errorf("At least one item is required")
	}

	for _, item := range receipt.Items {
		// Validate item description
		if !regexp.MustCompile(`^\w[\w\s-]*$`).MatchString(item.ShortDescription) {
			return fmt.Errorf("Item description is invalid")
		}

		// Validate item price
		if !regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(item.Price) {
			return fmt.Errorf("Item price is invalid")
		}
	}

	return nil
}

func findReceiptByID(id string) (*Receipt, error) {
	for i := range receipts {
		if receipts[i].ID == id {
			return &receipts[i], nil
		}
	}
	return nil, fmt.Errorf("receipt not found")
}

func generateID() string {
	// Generate a unique ID for the receipt
	// Replace this with your actual logic to generate IDs
	// Here, we are generating a random ID for simplicity

	return "adb6b560-0eef-42bc-9d16-df48f30e89b2"
}
