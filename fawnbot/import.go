// ================================================================
// import.go
// - Contains logic for importing data from external
//   sources, such as Google Sheets.
//================================================================

package fawnbot

import (
	"fmt"
	"strings"
	"time"
)

// Note: this read is highly-dependent on the configuration of the read sheet. If values or column order changes, the program will not return any configs.
func FetchCrawlConfigsFromSheet(sheetID, sheetName string) ([]SiteCrawlConfig, error) {
	service, err := startNewSheetsService()
	if err != nil {
		return []SiteCrawlConfig{}, fmt.Errorf("failed to start new Sheets service: %v", err)
	}

	readRange := sheetName + "!A2:G" // skip header, cols A-F
	response, err := service.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	var crawlConfigs []SiteCrawlConfig
	for _, row := range response.Values {
		if len(row) < 7 {
			continue //ignore incomplete rows
		}

		keepOldCrawls := false
		if val, ok := row[6].(string); ok && strings.ToLower(val) == "true" {
			keepOldCrawls = true
		}

		config := SiteCrawlConfig{
			Root:              row[0].(string),
			CrawlStart:        row[1].(string),
			CrawlFrequency:    row[2].(string),
			SheetName:         row[3].(string),
			AnalysisSheetName: row[4].(string),
			SheetID:           extractSheetIDFromURL(row[5].(string)),
			KeepOldCrawls:     keepOldCrawls,
		}

		crawlConfigs = append(crawlConfigs, config)
	}

	return crawlConfigs, nil
}

func IsSiteDue(s SiteCrawlConfig) (bool, error) {
	startDate, err := time.Parse("2006-01-02", s.CrawlStart)
	if err != nil {
		return false, fmt.Errorf("invalid date format: %v", err)
	}

	today := time.Now().Truncate(24 * time.Hour)
	startDate = startDate.Truncate(24 * time.Hour)

	switch strings.ToLower(s.CrawlFrequency) {
	case "daily":
		return true, nil
	case "weekly":
		daysPassed := int(today.Sub(startDate).Hours() / 24)
		return daysPassed >= 0 && daysPassed%7 == 0, nil
	case "fortnightly":
		daysPassed := int(today.Sub(startDate).Hours() / 24)
		return daysPassed >= 0 && daysPassed%14 == 0, nil
	case "monthly":
		return startDate.Day() == today.Day(), nil
	default:
		return false, fmt.Errorf("invalid crawl frequency: '%s'. Expected one of: daily, weekly, fortnightly, monthly", s.CrawlFrequency)
	}
}

func extractSheetIDFromURL(url string) string {
	if strings.Contains(url, "docs.google.com") {
		parts := strings.Split(url, "/")
		for i := 0; i < len(parts); i++ {
			if parts[i] == "d" && i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}
	return url
}
