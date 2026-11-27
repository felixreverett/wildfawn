// configs.go
// Defines the program's various configs/data objects

package fawnbot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type CrawlConfig struct {
	Root      string     `json:"root"`
	Core      *Core      `json:"core"`
	Limits    *Limits    `json:"limits"`
	PostCrawl *PostCrawl `json:"postCrawl"`
	API       *API       `json:"api,omitempty"`
}

type Core struct {
	RespectRobots bool `json:"RespectRobots"`
	MaxCrawlDepth int  `json:"MaxCrawlDepth"`
}

type Limits struct {
	MaxRequestsPerSecond int  `json:"MaxRequestsPerSecond"` // unimplemented
	MaxThreads           int  `json:"MaxThreads"`           // unimplemented
	SlowDownIf429        bool `json:"SlowDownIf429"`        // unimplemented
	RetryIf429           bool `json:"RetryIf429"`           // unimplemented
}

type PostCrawl struct {
	CollectPostCrawlMetrics bool `json:"CollectPostCrawlMetrics"` // unimplemented
}

type API struct {
	ReadSheetID       string `json:"ReadSheetID"`       // deactivated
	ReadSheetName     string `json:"ReadSheetName"`     // deactivated
	AutoUploadResults bool   `json:"AutoUploadResults"` // unimplemented
	CrawlStart        string `json:"CrawlStart"`
	CrawlFrequency    string `json:"CrawlFrequency"`
	SheetName         string `json:"SheetName"`
	AnalysisSheetName string `json:"AnalysisSheetName"`
	SheetID           string `json:"SheetID"`
	KeepOldCrawls     bool   `json:"KeepOldCrawls"` // Writes over LatestCrawl and makes a dated copy
}

/*
Defines any config required for a particular implementation of fawnbot (i.e., Cloud, Tundra)
*/
type ImplementationConfig struct {
	OnProgress func(url string, obj URLObject) // Optional callback for streaming URLObject data back to the caller during an ongoing crawl
}

// ================================================================
// Associated functions
// ================================================================

/*
Loads a single crawl config from a filepath
*/
func LoadCrawlConfig(filepath string) (CrawlConfig, error) {
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

/*
Load multiple crawl configs
Todo: Currently hardcoded to search for CrawlConfig.json
*/
func LoadCrawlConfigs() ([]CrawlConfig, error) {
	var crawlConfigs []CrawlConfig

	matches, err := filepath.Glob("configs/*CrawlConfig.json")
	if err != nil {
		return nil, fmt.Errorf("error reading configs in directory %w", err)
	}

	for _, filepath := range matches {
		cfg, err := LoadCrawlConfig(filepath)
		if err != nil {
			fmt.Printf("[!] Error loading %s: %v\n", filepath, err)
			continue
		}

		crawlConfigs = append(crawlConfigs, cfg)
	}

	return crawlConfigs, nil
}
