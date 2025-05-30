package fawnbot

import "fmt"

func goTame() {
	// substitutes a crawl (so I don't recrawl every time)
	/*URLObjects := make(map[string]*URLObject)
	URLObjects["https://example.com/page1"] = &URLObject{Inlinks: 5, Outlinks: 3, PageStatus: 200, CrawlDepth: 2}
	URLObjects["https://example.com/page2"] = &URLObject{Inlinks: 8, Outlinks: 2, PageStatus: 301, CrawlDepth: 3}
	URLObjects["https://example.com/page3"] = &URLObject{Inlinks: 2, Outlinks: 7, PageStatus: 404, CrawlDepth: 1}

	secrets, err := LoadSecrets("secrets.json")
	if err != nil {
		fmt.Println("Error loading secrets;", err)
		return
	}

	// Debug
	fmt.Println("Spreadsheet ID:", secrets.SheetID)
	fmt.Println("Sheet Name:", secrets.SheetName)

	if err := WriteToSheet(secrets.SheetID, secrets.SheetName, URLObjects); err != nil {
		fmt.Printf("Error: %v", err)
	}*/
}

func printSiteMap(robots Robots) {
	fmt.Println("(i) Debug - logging robots file:")
	fmt.Println(">   Agents:")
	for _, agent := range robots.Agents {
		fmt.Printf(">   Agent name: %s\n>   Agent disallows: %v\n>   Agent allows: %v\n", agent.Name, agent.Disallow, agent.Allow)
	}
	fmt.Println(">   Sitemaps:")
	for _, sitemap := range robots.Sitemaps {
		fmt.Println(">   ", sitemap)
	}
	fmt.Printf(">   CrawlDelay: %d\n", robots.CrawlDelay)
}

func printCrawlConfig(crawlConfig CrawlConfig) {
	fmt.Println("(i) Debug - logging crawl config")
	fmt.Printf("     Root: %s\n", crawlConfig.Root)
	fmt.Printf("     Start: %s\n", crawlConfig.CrawlStart)
	fmt.Printf("     Frequency: %s\n", crawlConfig.CrawlFrequency)
	fmt.Printf("     KeepOldCrawls: %t\n", crawlConfig.KeepOldCrawls)
}

func VerifyModuleImport() {
	fmt.Println("Successfully accessed function in wildfawn module")
}
