package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

/*
type config struct {
	respectRobots bool
	maxCrawlDepth int
	maxCrawlsPerSecond int
}*/

type URLObject struct {
	inlinks    int
	outlinks   int
	pageStatus int
	crawlDepth int
}

type QueueEntry struct {
	url        string
	crawlDepth int
}

// Send HTTP request to URL, returning HTML, response code, and any errors
func fetchURL(url string) (string, int, error) {
	//time.Sleep(time.Second * time.Duration(rand.Intn(5))) // Wait 1-5 seconds
	/*
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		response, err := client.Do(req)
		if err != nil {
			fmt.Println("Error:", err)
			return "", response.StatusCode, err
		}*/

	response, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", response.StatusCode, err
	}

	return string(body), response.StatusCode, nil
}

// returns a list of hrefs from an html string
func extractLinks(htmlString string) []string {
	var links []string
	tokenizer := html.NewTokenizer(strings.NewReader(htmlString))

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}

// Crawl all URLs on a site
func GoWild(root string) {
	start := time.Now()
	fmt.Printf("= = = Starting new crawl of %s = = =\n", root)

	// 1. Prepare regex to only crawl same-site URLs
	regexPattern := fmt.Sprintf("%s.*", root)

	re, err := regexp.Compile(regexPattern)

	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	// 2. Prepare objects data structure
	URLObjects := make(map[string]*URLObject)

	visitedURLs := make(map[string]bool)

	//todo: track depth info by tagging it onto enqueued urls

	// 3. Crawl site
	var URLQueue []QueueEntry
	URLQueue = append(URLQueue, QueueEntry{root, 0})
	visitedURLs[root] = true

	for len(URLQueue) > 0 {
		url := URLQueue[0].url
		fmt.Printf(". . . Queue size: %d | Crawling %s\n", len(URLQueue), url) //debug
		depth := URLQueue[0].crawlDepth

		URLQueue = URLQueue[1:]

		html, status, err := fetchURL(url)
		if err != nil {
			fmt.Println("> Error fetching URL:", err)
			return
		}

		// 3b. Iterate through every found link on url
		links := extractLinks(html)
		URLObjects[url] = &URLObject{inlinks: 1, outlinks: len(links), pageStatus: status, crawlDepth: depth}

		for _, link := range links {
			// resolve relative URLs
			if link[0] == '/' {
				if root[len(root)-1] == '/' {
					link = root[:len(root)-1] + link
				} else {
					link = root + link
				}
			}

			// check if already processed, else add to queue (if not current URL)
			if obj, ok := URLObjects[link]; ok {
				obj.inlinks++
			} else if re.MatchString(link) && !visitedURLs[link] {
				URLQueue = append(URLQueue, QueueEntry{link, depth + 1})
				visitedURLs[link] = true
			} else {
				//fmt.Println("URL invalid:", link) //debug
			}
		}
	}

	// 4. Results
	for key, value := range URLObjects {
		fmt.Printf("URL: %s\n ↳ Inlinks: %d | pageStatus: %d | outlinks: %d | crawl depth: %d\n", key, value.inlinks, value.pageStatus, value.outlinks, value.crawlDepth)
	}
	fmt.Printf("Successfully crawled %s\n", root)
	fmt.Printf(" ↳ Total URLs crawled: %d\n", len(URLObjects))
	fmt.Printf(" ↳ Total crawl time: %s\n", time.Since(start))
}

func main() {
	GoWild("https://web-scraping.dev/")
	//GoWild("https://yavsrg.net")
}
