package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strings"
	"time"

	baseDb "reportly/db"
	"reportly/model"
)

const (
	fileTypeCSV   = 1     // Enum for CSV format
	fileTypeExcel = 2     // Enum for Excel format
	reportCreated = true  // Flag indicating report is created
	reportPending = false // Flag indicating report is pending
)

// Run continuously checks and generates reports every 3 seconds
func Run() {
	for {
		fmt.Println("Checking for reports to generate...")
		PrepareReports()
		time.Sleep(3 * time.Second) // Wait before checking again
	}
}

// PrepareReports fetches all pending reports and processes them concurrently
func PrepareReports() bool {
	db := baseDb.GetDb()

	// Ensure ReportRequest table is migrated
	if err := db.AutoMigrate(&model.ReportRequest{}); err != nil {
		log.Fatalf("Migration error: %v", err)
		return false
	}

	var reports []model.ReportRequest
	// Fetch all reports where is_created_report is false
	err := db.Model(&model.ReportRequest{}).
		Where("is_created_report = ?", reportPending).
		Find(&reports).Error

	if err != nil {
		log.Fatalf("Query error: %v", err)
		return false
	}

	if len(reports) == 0 {
		return true
	}

	// Process each report in a separate goroutine
	for _, report := range reports {
		go createReport(report)
	}

	return true
}

// createReport handles generating the actual report file and updating status
func createReport(rq model.ReportRequest) {

	// Run SQL query and format data
	data, err := runQuery(rq)
	if err != nil {
		log.Printf("query execution failed: %s", err)
		return
	}

	// Generate appropriate file type
	switch rq.ReportFileType {
	case fileTypeCSV:
		generateCSV(data)
	case fileTypeExcel:
		generateExcel(data)
	default:
		log.Printf("unsupported file type: %d", rq.ReportFileType)
		return
	}

	// Mark report as created in database
	setReportCreated(rq.ID)
}

// runQuery executes the SQL and returns formatted data (rows as string slices)
func runQuery(rq model.ReportRequest) ([][]string, error) {
	db := baseDb.GetDb()

	// Parse field mappings from JSON (column -> header name)
	var fields map[string]string
	if err := json.Unmarshal([]byte(rq.Fields), &fields); err != nil {
		return nil, fmt.Errorf("invalid JSON fields: %w", err)
	}

	// Prepare SELECT columns and header row
	columns := make([]string, 0, len(fields))
	headers := make([]string, 0, len(fields))
	for col, header := range fields {
		columns = append(columns, col)
		headers = append(headers, header)
	}

	// Construct SQL query
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(columns, ", "), rq.TableName, rq.Query)

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("error executing raw query: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %s", err)
		}
	}()

	var data [][]string
	data = append(data, headers) // add header row

	// Reusable buffer for scanning values
	colsCount := len(columns)
	values := make([]interface{}, colsCount)
	valuePtrs := make([]interface{}, colsCount)

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Iterate over rows and convert to string slice
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("row scan error: %s", err)
			continue
		}

		row := make([]string, colsCount)
		for i, val := range values {
			row[i] = toString(val)
		}
		data = append(data, row)
	}

	return data, nil
}

// generateCSV writes the data to a .csv file
func generateCSV(data [][]string) {
	fileName := fmt.Sprintf("%s/%d.csv", os.Getenv("REPORT_PATH"), time.Now().Unix())
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("CSV file creation failed: %s", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("CSV file close failed: %s", err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write each row to the CSV
	for _, row := range data {
		if err := writer.Write(row); err != nil {
			log.Printf("CSV write error: %s", err)
		}
	}

	log.Println("CSV report generated: " + fileName)
}

// generateExcel writes the data to an .xlsx file using excelize
func generateExcel(data [][]string) {
	fileName := fmt.Sprintf("%s/%d.xlsx", os.Getenv("REPORT_PATH"), time.Now().Unix())
	xl := excelize.NewFile()
	sheet := "Sheet1"

	// Set each cell value
	for i, row := range data {
		for j, cell := range row {
			cellName, _ := excelize.CoordinatesToCellName(j+1, i+1)
			if err := xl.SetCellValue(sheet, cellName, cell); err != nil {
				log.Printf("Excel cell write error: %s", err)
			}
		}
	}

	// Save file to disk
	if err := xl.SaveAs(fileName); err != nil {
		log.Printf("Excel save failed: %s", err)
	} else {
		log.Println("Excel report generated: " + fileName)
	}
}

// setReportCreated updates the DB record to mark the report as generated
func setReportCreated(id uint) {
	db := baseDb.GetDb()
	if err := db.Model(&model.ReportRequest{}).
		Where("id = ?", id).
		Update("is_created_report", reportCreated).Error; err != nil {
		log.Printf("update report status failed: %s", err)
	}
}

// toString safely converts DB values to string
func toString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
