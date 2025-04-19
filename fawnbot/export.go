package main

/*
| - - export.go - -
| Contains functionality for post-crawl data exports to:
| - Excel
| - General-purpose JSON export (WIP)
*/

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func startNewSheetsService() (*sheets.Service, error) {
	credentials, err := os.ReadFile("service_account.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %v", err)
	}

	config, err := google.JWTConfigFromJSON(credentials, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %v", err)
	}

	client := config.Client(context.Background())

	service, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Sheets service: %v", err)
	}

	return service, nil
}

func sheetExists(service *sheets.Service, sheetID string, sheetName string) (bool, error) {
	spreadsheet, err := service.Spreadsheets.Get(sheetID).Do()
	if err != nil {
		return false, fmt.Errorf("failed to fetch spreadsheet: %v", err)
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			return true, nil
		}
	}

	return false, nil
}

func CreateNewSheet(service *sheets.Service, sheetID string, sheetName string) (int64, error) {
	exists, err := sheetExists(service, sheetID, sheetName)
	if err != nil {
		return 0, fmt.Errorf("error checking sheet existence: %v", err)
	}
	if exists {
		fmt.Printf(">   Sheet '%s' already exists, skipping creation\n", sheetName)
		return 0, nil
	}

	addSheetRequest := &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetName,
			},
		},
	}

	batchUpdate := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{addSheetRequest},
	}

	resp, err := service.Spreadsheets.BatchUpdate(sheetID, batchUpdate).Do()
	if err != nil {
		return 0, fmt.Errorf("failed to create new sheet: %v", err)
	}

	sheetIDNum := resp.Replies[0].AddSheet.Properties.SheetId
	return sheetIDNum, nil
}

func WriteToSheet(service *sheets.Service, sheetID string, sheetName string, URLObjectList URLObjectList) error {
	var err error
	data := URLObjectList.URLObjects
	start := time.Now()
	fmt.Printf(">   Writing export to sheets\n")

	// Convert URLObject to interface for Sheets
	var values [][]interface{}
	values = append(values, []interface{}{
		"URL", "Inlinks", "Outlinks", "Page Status", "Crawl Depth", "No Index", "Indexability", "Canonical", "Self-Canonicalises", "Is Canonical Indexable", "Is Orphan", "Blocked by Robots"}) //headers

	for url, obj := range data {
		row := []interface{}{
			url, obj.Inlinks, obj.Outlinks, obj.PageStatus, obj.CrawlDepth, obj.NoIndex, obj.Indexability, obj.Canonical, obj.IsSelfCanonicalising, obj.IsCanonicalIndexable, obj.IsOrphan, obj.IsBlockedByRobots}
		values = append(values, row)
	}

	// Define write range of sheet (e.g., "Sheet1!A1")
	writeRange := fmt.Sprintf("%s!A1", sheetName)

	// Prepare value range
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Clear any existing data on sheet
	clearRange := fmt.Sprintf("%s!A1:Z1000", sheetName)
	_, err = service.Spreadsheets.Values.Clear(sheetID, clearRange, &sheets.ClearValuesRequest{}).Do()
	if err != nil {
		return fmt.Errorf("failed to clear sheet before writing: %v", err)
	}

	// Write data to the sheet
	_, err = service.Spreadsheets.Values.Update(sheetID, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("failed to write data to sheet: %v", err)
	}

	fmt.Printf(">   Data successfully written to Sheet in %s\n", time.Since(start))

	return nil
}

func WriteWild(URLObjectList URLObjectList) {
	var err error

	// To be replaced with crawl config information
	keepOldCrawls := true

	// Load secrets (to be replaced with crawl config information)
	secrets, err := LoadSecrets("secrets.json")
	if err != nil {
		fmt.Println("[!] Error loading secrets:", err)
		return
	}

	// Establish new service
	service, err := startNewSheetsService()
	if err != nil {
		fmt.Println("[!] Could not create Sheets service:", err)
	}

	// Write analysis
	// todo

	// Write crawl
	_, err = CreateNewSheet(service, secrets.SheetID, secrets.SheetName)
	if err != nil {
		fmt.Println("[!] Error creating new sheet:", err)
	}

	if err := WriteToSheet(service, secrets.SheetID, secrets.SheetName, URLObjectList); err != nil {
		fmt.Println("[!] Error writing to sheet:", err)
	}

	// Export copy of crawl
	if keepOldCrawls {
		// Create timestamped sheetname
		timestamp := time.Now().Format("2006-01-02")
		newSheetName := fmt.Sprintf("Crawl %s", timestamp)

		_, err = CreateNewSheet(service, secrets.SheetID, newSheetName)
		if err != nil {
			fmt.Println("[!] Error creating new sheet:", err)
		}

		if err := WriteToSheet(service, secrets.SheetID, newSheetName, URLObjectList); err != nil {
			fmt.Println("[!] Error writing to sheet:", err)
		}
	}
}
