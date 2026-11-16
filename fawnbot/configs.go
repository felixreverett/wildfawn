// configs.go
// Defines the program's various configs/data objects

package fawnbot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

/*
Defines how the program executes a crawl instance. Can persist between crawls.
*/
type ProgramConfig struct {
	RespectRobots      bool   `json:"RespectRobots"`
	MaxCrawlDepth      int    `json:"MaxCrawlDepth"`      // unimplemented
	MaxCrawlsPerSecond int    `json:"MaxCrawlsPerSecond"` // unimplemented
	ReadSheetID        string `json:"ReadSheetID"`
	ReadSheetName      string `json:"ReadSheetName"`
}

/*
Defines the crawl data that is used to start a crawl. Unique to each crawl.
  - Todo: currently contains config information relevant to crawl scheduling, though
    this isn't inherent to the sitecrawlconfig and could be refactored
*/
type SiteCrawlConfig struct {
	Root              string `json:"Root"`
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
Loads a ProgramConfig object from a filepath
*/
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

/*
Loads a single crawl config from a filepath
*/
func LoadSiteCrawlConfig(filepath string) (SiteCrawlConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return SiteCrawlConfig{}, fmt.Errorf("failed to load file: %v", err)
	}

	var crawlConfig SiteCrawlConfig
	if err := json.Unmarshal(data, &crawlConfig); err != nil {
		return SiteCrawlConfig{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return crawlConfig, nil
}

/*
Load multiple crawl configs
Todo: Currently hardcoded to search for CrawlConfig.json
*/
func LoadSiteCrawlConfigs() ([]SiteCrawlConfig, error) {
	var crawlConfigs []SiteCrawlConfig

	matches, err := filepath.Glob("configs/*CrawlConfig.json")
	if err != nil {
		return nil, fmt.Errorf("error reading configs in directory %w", err)
	}

	for _, filepath := range matches {
		cfg, err := LoadSiteCrawlConfig(filepath)
		if err != nil {
			fmt.Printf("[!] Error loading %s: %v\n", filepath, err)
			continue
		}

		crawlConfigs = append(crawlConfigs, cfg)
	}

	return crawlConfigs, nil
}
