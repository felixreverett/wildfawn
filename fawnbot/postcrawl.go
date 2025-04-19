package main

/*
| - - postcrawl.go - -
| Post-crawl analysis for the more complicated metrics of URLObjects.
*/

// receiver function to fill in remaining data at end of crawl
func (u URLObjectList) RunPostCrawl() {

	for url, obj := range u.URLObjects {
		if obj.Inlinks == 0 {
			obj.IsOrphan = true
		}
		// 2. IsOnSitemap
		// todo: code here

		// 3. IsCanonicalIndexable, IsSelfCanonicalising
		if _, ok := u.URLObjects[obj.Canonical]; ok {
			if url == obj.Canonical {
				obj.IsSelfCanonicalising = true
			}
			if u.URLObjects[obj.Canonical].Indexability {
				obj.IsCanonicalIndexable = true
			}
		}
	}
}
