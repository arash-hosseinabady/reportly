# Reportly

**Reportly** is a lightweight Go service for generating dynamic reports (Excel/PDF) based on user-defined SQL templates.  
It listens for report generation requests, queries the database, and produces downloadable files asynchronously.

---

## Features

- Dynamic SQL report generation
- Excel / PDF export support
- Auto-migration with GORM
- Queue-based report processing
- Easy-to-extend model structure

---
## Installation
run `go buil` to generate compiled file then run that as a service.
