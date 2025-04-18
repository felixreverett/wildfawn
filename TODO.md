# 🦌 Wild Fawn todo:

## main.go:
- ⬜️ Oversee full import, crawl, and export
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
  - ⬜️ PageTitle
  - ⬜️ MetaDescription
  - ⬜️ H1*
- Post Parse metrics:
  - ⬜️ Canonical Indexability*
- Analysis metrics:
  - ✅ total URLs
  - ✅ total internal 200's
  - ✅ total internal 3xx's (and chains)
  - ✅ total internal 4xx's
  - ✅ total internal 5xx's
  - ⬜️ Pages with high crawl depth
  - ⬜️ non-indexable URLs in sitemap

### QOL:
- ✅ Parse relative & absolute URLs
- ✅ Evaluate site preference for naked or www. URLs
- ✅ Read robots.txt
  - ✅ Add option to respect robots.txt
  - ⬜️ check sitemaps

## crawler.go:
- ✅ Add respect disallows
- ✅ Add respect crawl delay

## import.go:
- ✅ Implement config
- ⬜️ Import from Google Sheet

## export.go:
- ✅ Write crawl to existing sheet in Google Sheets
- ⬜️ Write crawl to new sheet
- ⬜️ Write overview to existing sheet
- ⬜️ Write overview to new sheet
- ⬜️ Enable 'write copy' to write to both "latest" and "dated" sheets (based on config)