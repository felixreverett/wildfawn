package main

/*
| - - import.go - -
| Contains all relevant input data for main processing
|
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ProgramConfig struct {
	RespectRobots      bool `json:"RespectRobots"`      // unimplemented
	MaxCrawlDepth      int  `json:"MaxCrawlDepth"`      // unimplemented
	MaxCrawlsPerSecond int  `json:"MaxCrawlsPerSecond"` // unimplemented
}

func LoadConfig(filename string) (ProgramConfig, error) {
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

type Secrets struct {
	SheetID   string `json:"SheetID"`
	SheetName string `json:"SheetName"`
}

func LoadSecrets(filename string) (*Secrets, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets file: %v\nDoes the file exist on your local machine?", err)
	}

	var secrets Secrets
	if err := json.Unmarshal(data, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &secrets, nil
}

// - - -

type CrawlConfig struct {
	Root           string
	FirstAdded     string
	CrawlFrequency string
	SheetName      string
	SheetID        string
	KeepOldCrawls  bool // Writes over LatestCrawl and makes a dated copy
}

func ImportSiteCrawlInfo() /*[]SiteCrawlConfig*/ {
	// Import Site Crawl Configurations from dedicated Google Sheet (or appropriate source)
}

func (s CrawlConfig) IsSiteDue() (bool, error) {
	startDate, err := time.Parse("2006-01-02", s.FirstAdded)
	if err != nil {
		return false, fmt.Errorf("invalid date format: %v", err)
	}

	var duration time.Duration

	switch s.CrawlFrequency {
	case "daily":
		duration = 24 * time.Hour
	case "weekly":
		duration = 7 * 24 * time.Hour
	case "fortnightly":
		duration = 14 * 24 * time.Hour
	case "monthly":
		duration = 28 * 24 * time.Hour
	default:
		return false, fmt.Errorf("invalid crawl frequency: '%s'. Expected one of: daily, weekly, fortnightly, monthly", s.CrawlFrequency)
	}

	targetDate := startDate.Add(duration)

	return time.Now().After(targetDate), nil
}
