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
	Inlinks               int
	Outlinks              int
	PageStatus            int
	CrawlDepth            int
	NoIndex               bool
	Indexability          bool
	Canonical             string
	MetaTitle             string
	MetaTitleLength       int
	MetaDescription       string
	MetaDescriptionLength int
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

func parseHTML(htmlString string) (bool, bool, string) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return true, false, ""
	}

	indexable := true
	noIndex := false
	canonical := ""

	// 2. recursive function to parse html
	var traverseHTML func(*html.Node)
	traverseHTML = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// a. check for <meta name="robots" content="noindex">
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
					noIndex = true
				}
			}

			// b. check for canonical
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
	return indexable, noIndex, canonical
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

// returns a URL's preferred www. config by checking for redirects
func SetWWWPreference(root string) (string, error) {
	_, status, redirectTo, err := fetchURL(root)
	if err != nil {
		return root, err
	}

	if status >= 300 && status < 400 && redirectTo != "" {
		rootHost := extractHost(root)
		redirectHost := extractHost(redirectTo)

		if rootHost != redirectHost {
			fmt.Printf("Detected www preference: %s -> %s\n", root, redirectTo)
			return redirectTo, nil
		}
	}

	return root, nil
}

// extract host of URL
func extractHost(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	parts := strings.Split(url, "/")
	return parts[0]
}

func normaliseWWW(url, preferredRoot string) string {
	urlHost := extractHost(url)
	preferredHost := extractHost(preferredRoot)

	if strings.Contains(preferredHost, "www.") && !strings.Contains(urlHost, "www.") {
		return strings.Replace(url, urlHost, preferredHost, 1)
	} else if !strings.Contains(preferredHost, "www.") && strings.Contains(urlHost, "www.") {
		return strings.Replace(url, urlHost, preferredHost, 1)
	}
	return url
}

func Crawl(root string) (map[string]*URLObject, error) {

	// 1. prepare regex to only crawl same-site URLs
	host := extractHost(root)
	regexPattern := fmt.Sprintf("^https?://%s.*", host)
	rootRegex, err := regexp.Compile(regexPattern)

	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return nil, err
	}

	// 2. prepare data structures
	URLObjects := make(map[string]*URLObject)

	visitedURLs := make(map[string]bool)

	// 3. crawl every URL in a queue
	var URLQueue []QueueEntry
	URLQueue = append(URLQueue, QueueEntry{root, 0})
	visitedURLs[root] = true

	for len(URLQueue) > 0 {
		// a. pop
		url := URLQueue[0].url
		//fmt.Printf(". . . Queue size: %d | Crawling %s\n", len(URLQueue), url) //debug
		depth := URLQueue[0].crawlDepth

		URLQueue = URLQueue[1:]

		// b. fetch
		html, status, redirectTo, err := fetchURL(url)
		if err != nil {
			fmt.Println("> Error fetching URL:", err)
			return nil, err
		}

		indexable := false
		noIndex := false
		canonical := ""

		// c. check for redirect status
		if status >= 300 && status < 400 {
			if redirectTo != "" && !visitedURLs[redirectTo] {
				URLQueue = append(URLQueue, QueueEntry{redirectTo, depth})
				visitedURLs[redirectTo] = true
				fmt.Printf("> Redirect: %s → %s\n", url, redirectTo)
			}
		} else if status == 200 {
			indexable, noIndex, canonical = parseHTML(html)
		}

		// d. collect every link on current URL
		links := extractLinks(html)

		// e. add current URL results to URLObject
		URLObjects[url] = &URLObject{Inlinks: 1, Outlinks: len(links), PageStatus: status, CrawlDepth: depth,
			Indexability: indexable, NoIndex: noIndex, Canonical: canonical}

		// f. iterate through all links of current URL
		for _, link := range links {

			// i. ignore 0-length URLs
			if len(link) == 0 {
				continue
			}

			// ii. resolve relative URLs
			if link[0] == '/' {
				if root[len(root)-1] == '/' {
					link = root[:len(root)-1] + link
				} else {
					link = root + link
				}
			}

			// iii. ignore external urls
			if !rootRegex.MatchString(link) {
				continue
			}

			// iv. normalise URLs to WWW preference
			link = normaliseWWW(link, root)

			// v. check if URL already processed, else add to queue (if not current URL)
			if obj, ok := URLObjects[link]; ok {
				obj.Inlinks++
			} else if !visitedURLs[link] {
				URLQueue = append(URLQueue, QueueEntry{link, depth + 1})
				visitedURLs[link] = true
			}
		}
	}
	return URLObjects, nil
}

// Crawl all URLs on a site
func GoWild(root string) map[string]*URLObject {
	start := time.Now()
	fmt.Printf("= = = Starting new crawl of %s = = =\n", root)

	// 1. detect and set preference for www or non www
	root, err := SetWWWPreference(root)
	if err != nil {
		fmt.Println("Error detecting www preference:", err)
		return nil
	}
	fmt.Println("> Normalising all URLs to:", root) //debug

	// 2. Get robots
	robots, err := GetRobots(root)
	if err != nil {
		fmt.Println(err)
	} else {
		PrintSiteMap(robots)
	}

	return nil //debug

	// 2. Crawl site
	URLObjects, err := Crawl(root)
	if err != nil {
		fmt.Println("failed to crawl root: ", err)
		return nil
	}

	// 3. Return (and print) results
	/*
		for key, value := range URLObjects {
			fmt.Printf("URL: %s\n ↳ Inlinks: %d | pageStatus: %d | outlinks: %d | crawl depth: %d | indexable: %v | canonical: %s\n",
				key, value.Inlinks, value.PageStatus, value.Outlinks, value.CrawlDepth, value.Indexability, value.Canonical)
		}
	*/
	fmt.Printf("Successfully crawled %s\n", root)
	fmt.Printf(" ↳ Total URLs crawled: %d\n", len(URLObjects))
	fmt.Printf(" ↳ Total crawl time: %s\n", time.Since(start))

	return URLObjects
}
