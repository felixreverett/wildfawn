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
	Inlinks      int
	Outlinks     int
	PageStatus   int
	CrawlDepth   int
	Indexability bool
	Canonical    string
}

type QueueEntry struct {
	url        string
	crawlDepth int
}

// Send HTTP request to URL, returning HTML, response code, and any errors
func fetchURL(url string) (string, int, string, error) {
	//time.Sleep(time.Second * time.Duration(rand.Intn(2))) // Wait 1-2 seconds

	// Be respectful to the server by setting a user-agent 🙇🙇🙇
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, "", err
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Do(request)
	if err != nil {
		return "", 0, "", nil
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", response.StatusCode, "", err
	}

	redirectTo := ""
	if response.StatusCode >= 300 && response.StatusCode < 400 {
		redirectTo = response.Header.Get("Location")
	}

	return string(body), response.StatusCode, redirectTo, nil
}

func parseHTML(htmlString string) (bool, string) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return true, ""
	}

	indexable := true
	canonical := ""

	var traverseHTML func(*html.Node)
	traverseHTML = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// 2a. check for <meta name="robots" content="noindex">
			if n.Data == "meta" {
				var name, content string
				for _, attr := range n.Attr {
					if attr.Key == "name" && strings.ToLower(attr.Val) == "robots" {
						name = attr.Val
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				if name == "robots" && strings.Contains(strings.ToLower(content), "noindex") {
					indexable = false
				}
			}

			// 2b. check for canonical
			if n.Data == "link" {
				var rel, href string
				for _, attr := range n.Attr {
					if attr.Key == "rel" && attr.Val == "canonical" {
						rel = attr.Val
					}
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				if rel == "canonical" && href != "" {
					canonical = href
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverseHTML(c)
		}
	}

	traverseHTML(doc)
	return indexable, canonical
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

func Crawl(root string) (map[string]*URLObject, error) {

	// 1. Prepare regex to only crawl same-site URLs
	regexPattern := fmt.Sprintf("%s.*", root)

	re, err := regexp.Compile(regexPattern)

	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return nil, err
	}

	// 2. Prepare objects data structure
	URLObjects := make(map[string]*URLObject)

	visitedURLs := make(map[string]bool)

	// 3. Crawl site
	var URLQueue []QueueEntry
	URLQueue = append(URLQueue, QueueEntry{root, 0})
	visitedURLs[root] = true

	for len(URLQueue) > 0 {
		url := URLQueue[0].url
		fmt.Printf(". . . Queue size: %d | Crawling %s\n", len(URLQueue), url) //debug
		depth := URLQueue[0].crawlDepth

		URLQueue = URLQueue[1:]

		html, status, redirectTo, err := fetchURL(url)
		if err != nil {
			fmt.Println("> Error fetching URL:", err)
			return nil, err
		}

		indexable := false
		canonical := ""

		// Check for redirect status
		if status >= 300 && status < 400 {
			if redirectTo != "" && !visitedURLs[redirectTo] {
				URLQueue = append(URLQueue, QueueEntry{redirectTo, depth})
				visitedURLs[redirectTo] = true
				fmt.Printf("> Redirect: %s → %s\n", url, redirectTo)
			}
		} else if status == 200 {
			indexable, canonical = parseHTML(html)
		}

		// 3b. Iterate through every found link on url
		links := extractLinks(html)

		// 3c. Write URL results to URLObject
		URLObjects[url] = &URLObject{Inlinks: 1, Outlinks: len(links), PageStatus: status, CrawlDepth: depth, Indexability: indexable, Canonical: canonical}

		// 3d.
		for _, link := range links {

			// ignore 0-length URLs
			if len(link) == 0 {
				continue
			}

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
				obj.Inlinks++
			} else if re.MatchString(link) && !visitedURLs[link] {
				URLQueue = append(URLQueue, QueueEntry{link, depth + 1})
				visitedURLs[link] = true
			} else {
				//fmt.Println("URL invalid:", link) //debug
			}
		}
	}
	return URLObjects, nil
}

// Crawl all URLs on a site
func GoWild(root string) {
	start := time.Now()
	fmt.Printf("= = = Starting new crawl of %s = = =\n", root)

	URLObjects, err := Crawl(root)
	if err != nil {
		fmt.Println("failed to crawl root: ", err)
		return
	}

	// 4. Results
	for key, value := range URLObjects {
		fmt.Printf("URL: %s\n ↳ Inlinks: %d | pageStatus: %d | outlinks: %d | crawl depth: %d | indexable: %v | canonical: %s\n",
			key, value.Inlinks, value.PageStatus, value.Outlinks, value.CrawlDepth, value.Indexability, value.Canonical)
	}
	fmt.Printf("Successfully crawled %s\n", root)
	fmt.Printf(" ↳ Total URLs crawled: %d\n", len(URLObjects))
	fmt.Printf(" ↳ Total crawl time: %s\n", time.Since(start))

	fmt.Printf("Writing export to sheets")

	secrets, err := LoadSecrets("secrets.json")
	if err != nil {
		fmt.Println("Error loading secrets;", err)
		return
	}

	if err := WriteToSheet(secrets.SheetID, secrets.SheetName, URLObjects); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

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
