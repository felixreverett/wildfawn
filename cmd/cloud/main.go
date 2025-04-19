package main

import (
	"fmt"
	"os"

	"github.com/felixreverett/wildfawn/fawnbot"
)

func main() {
	// temporary bypass of imports
	data, err := os.ReadFile("rooturl.txt")
	if err != nil {
		fmt.Println("[!] Error loading file. Aborting crawl: ", err)
		return
	}
	url := string(data)

	// b. Load Config
	config, err := fawnbot.LoadConfig("programConfig.json")
	if err != nil {
		fmt.Println("[!] Error loading config. Using default: ", err)
	} else {
		fmt.Println("(i) Successfully loaded config")
	}

	// 2. Crawl
	URLObjectList, err := fawnbot.GoWild(url, config)
	if err != nil {
		fmt.Println("[!] Error crawling root URL. Aborting crawl. Error: ", err)
		return
	}

	// 3. Export
	fawnbot.WriteWild(URLObjectList)
}
