package fawnbot

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
	credentials, err := os.ReadFile("configs/service_account.json")
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

func createNewSheet(service *sheets.Service, sheetID string, sheetName string) (int64, error) {
	exists, err := sheetExists(service, sheetID, sheetName)
	if err != nil {
		return 0, fmt.Errorf("error checking sheet existence: %v", err)
	}
	if exists {
		fmt.Printf("(i) Sheet '%s' already exists, skipping creation.\n", sheetName)
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

func writeCrawlToSheet(service *sheets.Service, sheetID string, sheetName string, URLObjectList URLObjectList) error {
	var err error
	data := URLObjectList.URLObjects
	start := time.Now()
	fmt.Println("(i) Writing Crawl...")

	// Convert URLObject to interface for Sheets
	var values [][]interface{}
	values = append(values, []interface{}{
		"URL", "Inlinks", "Outlinks", "Page Status", "Crawl Depth",
		"No Index", "Indexability", "Canonical", "Self-Canonicalises", "Is Canonical Indexable",
		"Is Orphan", "Blocked by Robots",
		"Meta Title", "Meta Title Length", "Meta Description", "Meta Description Length", "H1", "H1 Length"}) //headers

	for url, obj := range data {
		row := []interface{}{
			url, obj.Inlinks, obj.Outlinks, obj.PageStatus, obj.CrawlDepth,
			obj.NoIndex, obj.Indexability, obj.Canonical, obj.IsSelfCanonicalising, obj.IsCanonicalIndexable,
			obj.IsOrphan, obj.IsBlockedByRobots,
			obj.MetaTitle, obj.MetaTitleLength, obj.MetaDescription, obj.MetaDescriptionLength, obj.H1, obj.H1Length}
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

	fmt.Printf("(i) Data successfully written to %s in %s\n", sheetName, time.Since(start))

	return nil
}

func writeAnalysis(service *sheets.Service, analysis CrawlAnalysis, crawlConfig CrawlConfig) error {
	start := time.Now()
	fmt.Println("(i) Writing Analysis...")

	// check if sheet exists
	exists, err := sheetExists(service, crawlConfig.SheetID, crawlConfig.AnalysisSheetName)
	if err != nil {
		return fmt.Errorf("failed to verify if analysis sheet exists: %v", err)
	}
	if !exists {
		_, err := createNewSheet(service, crawlConfig.SheetID, crawlConfig.AnalysisSheetName)
		if err != nil {
			return fmt.Errorf("failed to create analysis sheet: %v", err)
		}
		fmt.Printf("(i) Created new analysis sheet: %s\n", crawlConfig.AnalysisSheetName)
	}

	// Find first free row
	readRange := fmt.Sprintf("%s!A:A", crawlConfig.AnalysisSheetName)
	resp, err := service.Spreadsheets.Values.Get(crawlConfig.SheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("failed to read sheet to find first free row: %v", err)
	}

	firstFreeRow := len(resp.Values) + 1

	var values [][]interface{}

	if firstFreeRow == 1 {
		values = append(values, []interface{}{
			"Crawl Date", "Internal URLs", "200s", "300s", "400s", "500s",
			"Empty Meta Titles", "Empty Meta Descriptions", "Missing Canonicals", "No Indexes", "URLs Not In Sitemaps", "Non-Indexable URLs In Sitemaps", "Orphan URLs"})
	}

	today := time.Now().Format("2006-01-02")

	values = append(values, []interface{}{
		today, analysis.TotalInternalURLs, analysis.Total200s, analysis.Total300s, analysis.Total400s, analysis.Total500s,
		analysis.TotalEmptyMetaTitles, analysis.TotalEmptyMetaDescriptions, analysis.TotalMissingCanonicals, analysis.TotalNoIndexes, analysis.TotalNotInSitemap, analysis.TotalNonIndexableInSitemap, analysis.TotalOrphans})

	writeRange := fmt.Sprintf("%s!A%d", crawlConfig.AnalysisSheetName, firstFreeRow)

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Write data to the sheet
	_, err = service.Spreadsheets.Values.Update(crawlConfig.SheetID, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("failed to write data to sheet: %v", err)
	}

	fmt.Printf("(i) Data successfully written to Analysis sheet in %s.\n", time.Since(start))

	return nil
}

func WriteWild(URLObjectList URLObjectList, analysis CrawlAnalysis, crawlConfig CrawlConfig) {
	var err error

	// Establish new service
	service, err := startNewSheetsService()
	if err != nil {
		fmt.Println("[!] Could not create Sheets service:", err)
	}

	// Write analysis
	if err = writeAnalysis(service, analysis, crawlConfig); err != nil {
		fmt.Println("[!] Error writing crawl analysis:", err)
	}

	// Write crawl
	_, err = createNewSheet(service, crawlConfig.SheetID, crawlConfig.SheetName)
	if err != nil {
		fmt.Println("[!] Error creating new sheet:", err)
	}

	if err := writeCrawlToSheet(service, crawlConfig.SheetID, crawlConfig.SheetName, URLObjectList); err != nil {
		fmt.Println("[!] Error writing to sheet:", err)
	}

	// Export copy of crawl
	if crawlConfig.KeepOldCrawls {
		// Create timestamped sheetname
		timestamp := time.Now().Format("2006-01-02")
		newSheetName := fmt.Sprintf("Crawl %s", timestamp)

		_, err = createNewSheet(service, crawlConfig.SheetID, newSheetName)
		if err != nil {
			fmt.Println("[!] Error creating new sheet:", err)
		}

		if err := writeCrawlToSheet(service, crawlConfig.SheetID, newSheetName, URLObjectList); err != nil {
			fmt.Println("[!] Error writing to sheet:", err)
		}
	}
}
