package main

// Place for code to export crawl data to:
// - Excel
// - General-purpose JSON export

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

//func CreateNewSheet(sheetID string)

//func WriteNewSheet(sheetID string, data map[string]*URLObject) error {}

func WriteToSheet(sheetID string, sheetName string, data map[string]*URLObject) error {

	// 1. Load credentials
	credentials, err := os.ReadFile("service_account.json")
	if err != nil {
		return fmt.Errorf("error: failed to read credentials file: %v", err)
	}

	// 2. Configure a JSON Web Token to authenticate requests to the Google Sheets API
	config, err := google.JWTConfigFromJSON(credentials, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("error: failed to parse credentials: %v", err)
	}

	client := config.Client(context.Background())

	// Creating Google Sheets service
	service, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("error: failed to create Sheets service: %v", err)
	}

	// Convert URLObject to interface for Sheets
	var values [][]interface{}
	values = append(values, []interface{}{"URL", "Inlinks", "Outlinks", "Page Status", "Crawl Depth"}) //headers

	for url, obj := range data {
		row := []interface{}{url, obj.Inlinks, obj.Outlinks, obj.PageStatus, obj.CrawlDepth}
		values = append(values, row)
	}

	// Define write range of sheet (e.g., "Sheet1!A1")
	writeRange := fmt.Sprintf("%s!A1", sheetName)

	// Prepare value range
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Write data to the sheet
	_, err = service.Spreadsheets.Values.Update(sheetID, writeRange, valueRange).ValueInputOption("RAW").Do()

	if err != nil {
		return fmt.Errorf("error: failed to write data to sheet; %v", err)
	}

	fmt.Println("Data successfully written to Sheet")

	return nil
}
