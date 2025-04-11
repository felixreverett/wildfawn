package main

import "fmt"

func GoTame() {
	// substitutes a crawl (so I don't recrawl every time)
	URLObjects := make(map[string]*URLObject)
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
	}
}

func PrintSiteMap(robots Robots) {
	fmt.Println("Debug - logging robots file:")
	fmt.Println("Agents:")
	for _, agent := range robots.Agents {
		fmt.Printf("Agent name: %s\nAgent disallows: %v\nAgent allows: %v\n", agent.Name, agent.Disallow, agent.Allow)
	}
	fmt.Println("Sitemaps:")
	for _, sitemap := range robots.Sitemaps {
		fmt.Println(sitemap)
	}
	fmt.Printf("CrawlDelay: %d\n", robots.CrawlDelay)
}
