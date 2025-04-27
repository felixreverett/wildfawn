package fawnbot

/*
| - - import.go - -
| Contains all relevant input data for main processing
|
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ProgramConfig struct {
	RespectRobots      bool   `json:"RespectRobots"`
	MaxCrawlDepth      int    `json:"MaxCrawlDepth"`      // unimplemented
	MaxCrawlsPerSecond int    `json:"MaxCrawlsPerSecond"` // unimplemented
	ReadSheetID        string `json:"ReadSheetID"`
	ReadSheetName      string `json:"ReadSheetName"`
}

func LoadProgramConfig(filename string) (ProgramConfig, error) {
	defaultConfig := ProgramConfig{RespectRobots: false, MaxCrawlDepth: 99, MaxCrawlsPerSecond: 10}
	data, err := os.ReadFile(filename)
	if err != nil {
		return defaultConfig, fmt.Errorf("failed to load config file: %v", err)
	}

	var config ProgramConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return defaultConfig, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return config, nil
}

// - - -

type CrawlConfig struct {
	Root              string `json:"Root"`
	CrawlStart        string `json:"CrawlStart"`
	CrawlFrequency    string `json:"CrawlFrequency"`
	SheetName         string `json:"SheetName"`
	AnalysisSheetName string `json:"AnalysisSheetName"`
	SheetID           string `json:"SheetID"`
	KeepOldCrawls     bool   `json:"KeepOldCrawls"` // Writes over LatestCrawl and makes a dated copy
}

func loadCrawlConfig(filepath string) (CrawlConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return CrawlConfig{}, fmt.Errorf("failed to load file: %v", err)
	}

	var crawlConfig CrawlConfig
	if err := json.Unmarshal(data, &crawlConfig); err != nil {
		return CrawlConfig{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return crawlConfig, nil
}

func LoadCrawlConfigs() ([]CrawlConfig, error) {
	var crawlConfigs []CrawlConfig

	matches, err := filepath.Glob("configs/*CrawlConfig.json")
	if err != nil {
		return nil, fmt.Errorf("error reading configs in directory %w", err)
	}

	for _, filepath := range matches {
		cfg, err := loadCrawlConfig(filepath)
		if err != nil {
			fmt.Printf("[!] Error loading %s: %v\n", filepath, err)
			continue
		}

		crawlConfigs = append(crawlConfigs, cfg)
	}

	return crawlConfigs, nil
}

// Note: this read is highly-dependent on the configuration of the read sheet. If values or column order changes, the program will not return any configs.
func FetchCrawlConfigsFromSheet(sheetID, sheetName string) ([]CrawlConfig, error) {
	service, err := startNewSheetsService()
	if err != nil {
		return []CrawlConfig{}, fmt.Errorf("failed to start new Sheets service: %v", err)
	}

	readRange := sheetName + "!A2:F" // skip header, cols A-F
	response, err := service.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	var crawlConfigs []CrawlConfig
	for _, row := range response.Values {
		if len(row) < 6 {
			continue //ignore incomplete rows
		}

		keepOldCrawls := false
		if val, ok := row[5].(string); ok && strings.ToLower(val) == "true" {
			keepOldCrawls = true
		}

		config := CrawlConfig{
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

func IsSiteDue(s CrawlConfig) (bool, error) {
	startDate, err := time.Parse("2006-01-02", s.CrawlStart)
	if err != nil {
		return false, fmt.Errorf("invalid date format: %v", err)
	}

	today := time.Now().Truncate(24 * time.Hour)
	startDate = startDate.Truncate(24 * time.Hour)
	daysPassed := int(today.Sub(startDate).Hours() / 24)

	var interval int

	switch strings.ToLower(s.CrawlFrequency) {
	case "daily":
		interval = 1
	case "weekly":
		interval = 7
	case "fortnightly":
		interval = 14
	case "monthly":
		interval = 28
	default:
		return false, fmt.Errorf("invalid crawl frequency: '%s'. Expected one of: daily, weekly, fortnightly, monthly", s.CrawlFrequency)
	}

	//fmt.Printf("Root: %s, Start: %s, Days Passed: %d\n", s.Root, startDate.Format("2006-01-02"), daysPassed)

	return daysPassed >= 0 && daysPassed%interval == 0, nil
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
