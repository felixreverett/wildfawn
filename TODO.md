# 🦌 Wild Fawn todo:

## main.go:
- ⬜️ Oversee full import, crawl, and export
- ⬜️ Add gocurrency

## crawler.go:

### data collection:
- ✅ Track crawl depth
- ✅ Track page outlinks
- ✅ Track 300's and redirectee URL
- ⬜️ Track OnSitemap bool
- ⬜️ Track IsOrphan bool
- ⬜️ Track IsAllopatric bool (separated cluster)

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