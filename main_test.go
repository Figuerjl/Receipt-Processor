package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestProcessReceipt(t *testing.T) {
	receiptJSON := `
	{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{
				"shortDescription": "Mountain Dew 12PK",
				"price": 6.49
			},
			{
				"shortDescription": "Emils Cheese Pizza",
				"price": 12.25
			},
			{
				"shortDescription": "Knorr Creamy Chicken",
				"price": 1.26
			},
			{
				"shortDescription": "Doritos Nacho Cheese",
				"price": 3.35
			},
			{
				"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
				"price": 12.00
			}
		],
		"total": 35.35
	}
	`

	// Send the receipt JSON and retrieve the awarded points
	points, err := sendReceiptAndGetPoints(receiptJSON)
	if err != nil {
		t.Fatal("Failed to send receipt and get points:", err)
	}

	// Verify the awarded points
	expectedPoints := 28
	if points != expectedPoints {
		t.Errorf("Expected points: %d, but got: %d", expectedPoints, points)
	}
}

func sendReceiptAndGetPoints(receiptJSON string) (int, error) {
	// Send a request to the Process Receipts endpoint
	processReceiptURL := "http://localhost:8080/receipts/process"
	resp, err := http.Post(processReceiptURL, "application/json", bytes.NewBuffer([]byte(receiptJSON)))
	if err != nil {
		return 0, fmt.Errorf("failed to send request to Process Receipts endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response JSON
	var processResp struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&processResp)
	if err != nil {
		return 0, fmt.Errorf("failed to parse response from Process Receipts endpoint: %v", err)
	}

	// Send a request to the Get Points endpoint
	getPointsURL := fmt.Sprintf("http://localhost:8080/receipts/%s/points", processResp.ID)
	resp, err = http.Get(getPointsURL)
	if err != nil {
		return 0, fmt.Errorf("failed to send request to Get Points endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response JSON
	var pointsResp struct {
		Points int `json:"points"`
	}
	err = json.NewDecoder(resp.Body).Decode(&pointsResp)
	if err != nil {
		return 0, fmt.Errorf("failed to parse response from Get Points endpoint: %v", err)
	}

	return pointsResp.Points, nil
}
