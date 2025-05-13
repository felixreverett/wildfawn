# 🦌 Wild Fawn todo:

## main.go:
- ✅ Oversee full import, crawl, and export
- ⬜️ Add gocurrency

### data collection:
- Crawl metrics:
  - ✅ Inlinks
  - ✅ Outlinks
  - ✅ Page Status (including 3xxs)
  - ✅ Crawl Depth
  - ✅ No Index
  - ✅ Indexability
  - ✅ Canonical
  - ⬜️ OnSitemap bool
  - ✅ IsOrphan bool
  - ✅ Is Canonical Indexable bool
  - ✅ Is Self Canonicalising bool
  - ⬜️ IsAllopatric bool (separated cluster)
  - ✅ BlockedByRobots bool
  - ✅ PageTitle
  - ✅ PageTitleLength
  - ✅ MetaDescription
  - ✅ MetaDescriptionLength
  - ✅ H1
  - ✅ H1Length
- Post Parse metrics:
  - ⬜️ Canonical Indexability*
- Analysis metrics:
  - General:
    - ✅ total URLs
    - ✅ total internal 200's
    - ✅ total internal 3xx's (and chains)
    - ✅ total internal 4xx's
    - ✅ total internal 5xx's
    - Meta:
      - ⬜️ total missing page titles
      - ⬜️ total missing meta descriptions
      - ⬜️ total missing h1s
      - ⬜️ total multiple page titles
      - ⬜️ total multiple meta descriptions
      - ⬜️ total multiple h1s
  - ⬜️ Pages with high crawl depth
  - ⬜️ non-indexable URLs in sitemap

### QOL:
- ✅ Parse relative & absolute URLs
- ✅ Evaluate site preference for naked or www. URLs
- ✅ Read robots.txt
  - ✅ Add option to respect robots.txt
  - ⬜️ check sitemaps
- ✅ Crawl scheduling

## crawler.go:
- ✅ Add respect disallows
- ✅ Add respect crawl delay

## import.go:
- ✅ Implement crawl config (JSON)
- ✅ Implement program config (JSON)
- ✅ Import crawl and program configs (JSON)
- ✅ Import from Google Sheet
- ✅ Calculate if site is due based on daily/weekly/fortnightly startdate
- ✅ Calculate if site is due based on monthly startdate

## export.go:
- ✅ Write crawl to existing sheet in Google Sheets
- ✅ Write crawl to new sheet
- ⬜️ Write overview to existing sheet
- ⬜️ Write overview to new sheet
- ✅ Enable 'keepOldCrawls' to write to both "latest" and "dated" sheets (based on crawl config)