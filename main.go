package main

import (
	"fmt"
	"os"
)

func main() {
	// temporary bypass of imports
	data, err := os.ReadFile("rooturl.txt")
	if err != nil {
		fmt.Println("Error loading File: ", err)
	}
	url := string(data)

	// 2. Crawl
	URLObjects := GoWild(url)

	// 3. Export
	WriteWild(URLObjects)
	//GoTame() //test method which bypasses a crawl and exports data
}
