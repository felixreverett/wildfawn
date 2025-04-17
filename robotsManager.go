package main

/*
| - - robotsManager.go - -
| Contains functionality for parsing robots.txt
*/

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Robots struct {
	Agents     []UserAgent
	Sitemaps   []string
	CrawlDelay int
}

type UserAgent struct {
	Name     string
	Disallow []string
	Allow    []string
}

func ParseRobots(robotsFile string) Robots {
	var robots Robots
	var currentAgent *UserAgent

	lines := strings.Split(robotsFile, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		subline := strings.SplitN(line, ":", 2)
		if len(subline) != 2 {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(subline[0]))
		val := strings.TrimSpace(subline[1])

		switch key {
		case "user-agent":
			if currentAgent != nil {
				robots.Agents = append(robots.Agents, *currentAgent)
			}
			currentAgent = &UserAgent{Name: val}
		case "disallow":
			if currentAgent != nil {
				currentAgent.Disallow = append(currentAgent.Disallow, val)
			}
		case "allow":
			if currentAgent != nil {
				currentAgent.Allow = append(currentAgent.Allow, val)
			}
		case "sitemap":
			robots.Sitemaps = append(robots.Sitemaps, val)
		case "crawl-delay":
			delay, err := strconv.Atoi(val)
			if err == nil {
				robots.CrawlDelay = delay
			}
		}
	}

	if currentAgent != nil {
		robots.Agents = append(robots.Agents, *currentAgent)
	}

	return robots
}

func ExtractRootURL(inputURL string) (string, error) {
	parsed, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host), nil
}

func GetRobots(url string) (Robots, error) {
	pruned, err := ExtractRootURL(url)
	if err != nil {
		return Robots{}, err
	}

	robotsURL := pruned + "/robots.txt"
	robotsFile, status, _, err := fetchURLQuick(robotsURL)

	if status == 200 {
		fmt.Printf("(i) Found robots file at %s\n", robotsURL)
		return ParseRobots(robotsFile), nil
	} else {
		fmt.Printf("[!] Could not find robots file at %s\n", robotsURL)
		return Robots{}, err
	}
}

func IsURLBlockedByRobots(url string, robots Robots) bool {
	url = strings.ToLower(url)

	for _, agent := range robots.Agents {
		if agent.Name == "*" || strings.Contains("fawnbot", strings.ToLower(agent.Name)) {
			for _, allow := range agent.Allow {
				if strings.HasPrefix(url, allow) {
					return false
				}
			}

			for _, disallow := range agent.Disallow {
				if disallow == "/" || strings.HasPrefix(url, disallow) {
					return true
				}
			}
		}
	}

	return false
}
