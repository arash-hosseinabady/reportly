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
run `go buil` to generate a compiled file then run that as a service.

## Request
To generate a report file, you should insert new record in `report_requests` table.
Fields that should be filled:
- table_name: name of table that you want to create a report from that. Also, you can set join with another table
- query: set where clause
- fields: as JSON set name of fields that you want to select from `query` and name of title of this field that you want
    to show in a report file.

    `{"field_1": "title_field_1", "field_2": "title_field_2", "field_3": "title_field_3"}`
- report_file_type: type of report file. `csv = 1` and `excel = 2`
- is_created_report: set 0 that this service recognizes should be generating report file based on this record.
  After generated report file, this field set to 1.