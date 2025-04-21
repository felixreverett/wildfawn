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
	"time"
)

type ProgramConfig struct {
	RespectRobots      bool `json:"RespectRobots"`      // unimplemented
	MaxCrawlDepth      int  `json:"MaxCrawlDepth"`      // unimplemented
	MaxCrawlsPerSecond int  `json:"MaxCrawlsPerSecond"` // unimplemented
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
	Root           string `json:"Root"`
	CrawlStart     string `json:"CrawlStart"`
	CrawlFrequency string `json:"CrawlFrequency"`
	SheetName      string `json:"SheetName"`
	SheetID        string `json:"SheetID"`
	KeepOldCrawls  bool   `json:"KeepOldCrawls"` // Writes over LatestCrawl and makes a dated copy
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

func (s CrawlConfig) isSiteDue() (bool, error) {
	startDate, err := time.Parse("2006-01-02", s.CrawlStart)
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
