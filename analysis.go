package main

/*
| - - analysis.go - -
| Summary analysis for aggregated metrics
*/

type CrawlAnalysis struct {
	TotalInternalURLs          int
	Total200s                  int
	Total300s                  int
	Total400s                  int
	Total500s                  int
	TotalEmptyMetaTitles       int
	TotalEmptyMetaDescriptions int
	TotalMissingCanonicals     int
	TotalNoIndexes             int
	TotalNotInSitemap          int
	TotalNonIndexableInSitemap int
	TotalOrphans               int
}

func AnalyseCrawl(URLObjects map[string]*URLObject) CrawlAnalysis {
	var analysis CrawlAnalysis

	for _, URLObject := range URLObjects {
		// 1. Core data
		analysis.TotalInternalURLs = len(URLObjects)

		// 2. Status
		if URLObject.PageStatus >= 500 {
			analysis.Total500s++
		} else if URLObject.PageStatus >= 400 {
			analysis.Total400s++
		} else if URLObject.PageStatus >= 300 {
			analysis.Total300s++
		} else if URLObject.PageStatus >= 200 {
			analysis.Total200s++
		}

		// ?? meta details (indexable and non-indexable)
		if URLObject.MetaTitleLength == 0 {
			analysis.TotalEmptyMetaTitles++
		}
		if URLObject.MetaDescriptionLength == 0 {
			analysis.TotalEmptyMetaDescriptions++
		}

		// 4. canonicals
		if len(URLObject.Canonical) == 0 {
			analysis.TotalMissingCanonicals++
		}

		// 5. Directives
		if URLObject.NoIndex {
			analysis.TotalNoIndexes++
		}
	}

	return analysis
}
