package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price"`
}

type ProcessResponse struct {
	ID     string `json:"id"`
	Points int    `json:"points"`
}

func TestProcessReceipt(t *testing.T) {
	httpposturl := "http://localhost:8080/receipts/process"

	fmt.Println("TEST CASE: 1", httpposturl)
	var receipt1jsonData = []byte(`{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
		  {
			"shortDescription": "Mountain Dew 12PK",
			"price": "6.49"
		  },{
			"shortDescription": "Emils Cheese Pizza",
			"price": "12.25"
		  },{
			"shortDescription": "Knorr Creamy Chicken",
			"price": "1.26"
		  },{
			"shortDescription": "Doritos Nacho Cheese",
			"price": "3.35"
		  },{
			"shortDescription": "Klarbrunn 12-PK 12 FL OZ",
			"price": "12.00"
		  }
		],
		"total": "35.35"
	  }`)

	// Send the receipt and get the ID and points
	id, points, err := sendReceiptAndGetPoints(httpposturl, receipt1jsonData)
	if err != nil {
		t.Fatalf("Failed to send receipt and get points: %v", err)
	}

	// Verify the awarded points
	expectedPoints := 28
	if points != expectedPoints {
		t.Errorf("Expected points: %d, but got: %d", expectedPoints, points)
	}

	fmt.Println("ID:", id)
	fmt.Println("Points:", points)
}

func sendReceiptAndGetPoints(url string, receipt []byte) (string, int, error) {
	// Convert receipt object to JSON

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(receipt))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read response body: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("received non-OK response: %s", response.Status)
	}

	var processResp ProcessResponse
	err = json.Unmarshal(body, &processResp)
	if err != nil {
		return "", 0, fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	return processResp.ID, processResp.Points, nil
}
