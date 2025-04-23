package main

import (
	"fmt"

	"github.com/felixreverett/wildfawn/fawnbot"
)

func main() {
	// a. Load program config
	programConfig, err := fawnbot.LoadProgramConfig("configs/programConfig.json")
	if err != nil {
		fmt.Println("[!] Error loading program config. Using default:", err)
	} else {
		fmt.Println("(i) Successfully loaded program config")
	}

	// b. Load crawl configs
	//crawlConfigs, err := fawnbot.LoadCrawlConfigs()
	crawlConfigs, err := fawnbot.FetchCrawlConfigsFromSheet(programConfig.ReadSheetID, programConfig.ReadSheetName)
	if err != nil {
		fmt.Println("[!] Error loading crawl configs:", err)
	}

	// c. Crawl and export all
	for _, crawlConfig := range crawlConfigs {
		URLObjectList, err := fawnbot.GoWild(crawlConfig, programConfig)
		if err != nil {
			fmt.Println("[!] Error crawling root URL, aborting:", err)
			continue
		}
		analysis := fawnbot.AnalyseCrawl(URLObjectList)

		fawnbot.WriteWild(URLObjectList, analysis, crawlConfig)
	}
}
