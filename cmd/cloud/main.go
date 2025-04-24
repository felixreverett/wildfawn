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

		ok, err := fawnbot.IsSiteDue(crawlConfig)
		if err != nil {
			fmt.Println("[!] Error determining if site is due:", err)
			continue
		}

		if ok {
			fmt.Printf("(i) Site %s is due. Crawling\n", crawlConfig.Root)
			URLObjectList, err := fawnbot.GoWild(crawlConfig, programConfig)
			if err != nil {
				fmt.Println("[!] Error crawling root URL, aborting:", err)
				continue
			}
			analysis := fawnbot.AnalyseCrawl(URLObjectList)

			fawnbot.WriteWild(URLObjectList, analysis, crawlConfig)
		} else {
			fmt.Printf("(i) Site %s is not due\n", crawlConfig.Root)
		}

	}
}
