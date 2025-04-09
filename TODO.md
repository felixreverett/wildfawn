# 🦌 Wild Fawn todo:

## main.go:
- ⬜️ Oversee full import, crawl, and export
- ⬜️ Add gocurrency

## crawler.go:

### data collection:
- Crawl metrics:
  - ✅ Inlinks
  - ✅ Outlinks
  - ✅ Page Status (including 3xxs)
  - ✅ Crawl Depth
  - ✅ Indexability
  - ✅ Canonical
  - ⬜️ OnSitemap bool
  - ⬜️ IsOrphan bool
  - ⬜️ IsAllopatric bool (separated cluster)
  - ⬜️ BlockedByRobots bool
  - ⬜️ PageTitle
  - ⬜️ MetaDescription
  - ⬜️ H1*
  - ⬜️ Soft 404
- Post Parse metrics:
  - ⬜️ Canonical Indexability*
- Analysis metrics:
  - ⬜️ total URLs
  - ⬜️ total internal 200's
  - ⬜️ total internal 3xx's (and chains)
  - ⬜️ total internal 4xx's
  - ⬜️ Pages with high crawl depth
  - ⬜️ non-indexable URLs in sitemap

### QOL:
- ✅ Parse relative & absolute URLs
- ✅ Evaluate site preference for naked or www. URLs
- ⬜️ Read robots.txt
  - ⬜️ Add option to respect robots.txt
  - ⬜️ Read sitemaps

## import.go:
- ⬜️ Implement config

## export.go:
- ✅ Write crawl to existing sheet in Google Sheets
- ⬜️ Write crawl to new sheet
- ⬜️ Write overview to existing sheet
- ⬜️ Write overview to new sheet
- ⬜️ Enable 'write copy' to write to both "latest" and "dated" sheets (based on config)