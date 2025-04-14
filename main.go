package main

import (
	"fmt"
	"os"
)

func main() {
	// temporary bypass of imports
	data, err := os.ReadFile("rooturl.txt")
	if err != nil {
		fmt.Println("[!] Error loading file. Aborting crawl: ", err)
		return
	}
	url := string(data)

	// 2. Crawl
	URLObjectList, err := GoWild(url)
	if err != nil {
		fmt.Println("[!] Error crawling root URL. Aborting crawl. Error: ", err)
		return
	}

	// 3. Export
	WriteWild(URLObjectList)
}
