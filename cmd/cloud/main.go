package main

import (
	"fmt"

	"github.com/felixreverett/wildfawn/fawnbot"
)

func main() {
	// a. Load crawl configs
	crawlConfigs, err := fawnbot.LoadCrawlConfigs()
	if err != nil {
		fmt.Println("[!] Error loading crawl configs:", err)
	}

	// b. Load program config
	programConfig, err := fawnbot.LoadProgramConfig("programConfig.json")
	if err != nil {
		fmt.Println("[!] Error loading program config. Using default:", err)
	} else {
		fmt.Println("(i) Successfully loaded program config")
	}

	// 2. Crawl and export all
	for _, crawlConfig := range crawlConfigs {
		URLObjectList, err := fawnbot.GoWild(crawlConfig, programConfig)
		if err != nil {
			fmt.Println("[!] Error crawling root URL, aborting:", err)
			continue
		}
		fawnbot.WriteWild(URLObjectList, crawlConfig)
	}
}
