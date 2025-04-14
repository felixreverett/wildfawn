# 🦌 wildfawn SEO Crawler

A small project to develop a propietary lightweight web crawler to monitor technical hygiene metrics of any given site, without the costly drawbacks of larger software solutions. As a bonus, it's a great opportunity to further my knowledge of Go, too!

## Commit messages
- ✨ New feature
- 🪲 Bug fix
- 🧹 Code cleanup
- 📖 Other (e.g. documentation)

## Why is it called wild fawn?
- Let me know if you figure it out!

## Documentation
1. Each file is designed to contain types and functions for specific purposes:

| File             | Functionality                                                                              |
| ---------------- | ------------------------------------------------------------------------------------------ |
| analysis.go      | Handles the preparation of the CrawlAnalysis object for crawl summary.                     |
| crawler.go       | Anything involving the actual HTML data collection.                                        |
| debug.go         | Place for miscellaneous helper functions as part of the development process.               |
| export.go        | Code for exporting CrawlAnalysis and URLObject data. Currently only configured for Sheets. |
| import.go        | Handles import of any API keys and crawl instructions.                                     |
| main.go          | Entry point.                                                                               |
| postcrawl.go     | Post-crawl analysis for the more complicated metrics of URLObjects.                        |
| robotsManager.go | Handles all functionality for parsing the root's robots.txt file.                          |

## Fun technical features in this project
- receiver functions (see postcrawl.go)
    - I'm really just forcing a use of these to improve my familiarity with the syntax.
- pointers - an egregious disregard for functional programming practices :)
    - I've been using these a bit to improve my understanding of pointers, but I've pulled back a little because I prefer pure functions.
- clean method signatures (see crawler.go)