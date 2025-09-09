package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"api-service/internal/constants"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// Export data limits for different formats
const (
	MaxExportRecordsExcel = 100000 // Excel format max 100,000 records
	MaxExportRecordsOther = 500000 // Other formats max 500,000 records
)

// ExportFormat defines the export format type
type ExportFormat string

const (
	FormatJSON  ExportFormat = "json"
	FormatCSV   ExportFormat = "csv"
	FormatExcel ExportFormat = "excel"
)

// QueryBuilder is a function type that applies conditions to a GORM query
type QueryBuilder func(*gorm.DB) *gorm.DB

// ExportConfig represents the export configuration
type ExportConfig struct {
	TableName    string       `json:"table_name"`
	Fields       []string     `json:"fields"`
	QueryBuilder QueryBuilder `json:"-"` // Function to build query conditions
	Format       ExportFormat `json:"format"`
	Limit        int          `json:"limit"`
	OrderBy      string       `json:"order_by"`
}

// DBExporter is the database export utility
type DBExporter struct {
	db *gorm.DB
}

// NewDBExporter creates a new database exporter instance
func NewDBExporter(db *gorm.DB) *DBExporter {
	return &DBExporter{db: db}
}

// Export exports data according to the given configuration
func (e *DBExporter) Export(config *ExportConfig) ([]byte, error) {
	if err := e.validateConfig(config); err != nil {
		return nil, err
	}

	e.setDefaultLimits(config)

	results, err := e.executeQuery(config)
	if err != nil {
		return nil, err
	}

	return e.exportData(results, config)
}

// validateConfig validates the export configuration
func (e *DBExporter) validateConfig(config *ExportConfig) error {
	if config == nil {
		return errors.New("export config cannot be nil")
	}

	if config.TableName == "" {
		return errors.New("table name cannot be empty")
	}

	if config.Format == "" {
		config.Format = FormatJSON
	}

	return nil
}

// setDefaultLimits sets default export limits based on format
func (e *DBExporter) setDefaultLimits(config *ExportConfig) {
	switch config.Format {
	case FormatExcel:
		if config.Limit == 0 || config.Limit > MaxExportRecordsExcel {
			config.Limit = MaxExportRecordsExcel
		}
	default:
		if config.Limit == 0 || config.Limit > MaxExportRecordsOther {
			config.Limit = MaxExportRecordsOther
		}
	}
}

// executeQuery builds and executes the database query
func (e *DBExporter) executeQuery(config *ExportConfig) ([]map[string]interface{}, error) {
	query := e.db.Table(config.TableName)

	// Apply query conditions using the QueryBuilder function
	if config.QueryBuilder != nil {
		query = config.QueryBuilder(query)
	}

	// Apply field selection
	if len(config.Fields) > 0 && !(len(config.Fields) == 1 && config.Fields[0] == "*") {
		query = query.Select(config.Fields)
	}

	// Apply ordering
	if config.OrderBy != "" {
		query = query.Order(config.OrderBy)
	}

	// Apply limit
	if config.Limit > 0 {
		query = query.Limit(config.Limit)
	}

	// Execute query and get results
	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	return results, nil
}

// exportData exports the results in the specified format
func (e *DBExporter) exportData(results []map[string]interface{}, config *ExportConfig) ([]byte, error) {
	switch config.Format {
	case FormatJSON:
		return e.exportToJSON(results)
	case FormatCSV:
		return e.exportToCSV(results, config.Fields)
	case FormatExcel:
		return e.exportToExcel(results, config.Fields, config.TableName)
	default:
		return nil, errors.New("unsupported export format")
	}
}

// exportToJSON exports data to JSON format
func (e *DBExporter) exportToJSON(data []map[string]interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

// exportToCSV exports data to CSV format
func (e *DBExporter) exportToCSV(data []map[string]interface{}, fields []string) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	headers := e.prepareCSVHeaders(data, fields)
	if err := e.writeCSVHeaders(writer, headers); err != nil {
		return nil, err
	}

	if err := e.writeCSVData(writer, data, headers); err != nil {
		return nil, err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buffer.Bytes(), nil
}

// prepareCSVHeaders prepares headers for CSV export
func (e *DBExporter) prepareCSVHeaders(data []map[string]interface{}, fields []string) []string {
	headers := fields
	if len(headers) == 0 || (len(headers) == 1 && headers[0] == "*") {
		return e.extractSortedHeaders(data[0])
	}
	return headers
}

// writeCSVHeaders writes headers to CSV
func (e *DBExporter) writeCSVHeaders(writer *csv.Writer, headers []string) error {
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}
	return nil
}

// writeCSVData writes data rows to CSV
func (e *DBExporter) writeCSVData(writer *csv.Writer, data []map[string]interface{}, headers []string) error {
	for _, record := range data {
		row := make([]string, len(headers))
		for i, header := range headers {
			if value, exists := record[header]; exists && value != nil {
				row[i] = e.formatValue(value)
			}
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	return nil
}

// exportToExcel exports data to Excel format
func (e *DBExporter) exportToExcel(data []map[string]interface{}, fields []string, sheetName string) ([]byte, error) {
	file := excelize.NewFile()
	defer file.Close()

	if sheetName == "" {
		sheetName = "Sheet1"
	}

	// Set sheet name
	if err := file.SetSheetName("Sheet1", sheetName); err != nil {
		return nil, fmt.Errorf("failed to set sheet name: %w", err)
	}

	if len(data) == 0 {
		return e.writeEmptyExcelFile(file)
	}

	headers := e.prepareExcelHeaders(data, fields)
	if err := e.writeExcelHeaders(file, sheetName, headers); err != nil {
		return nil, err
	}

	if err := e.writeExcelData(file, sheetName, data, headers); err != nil {
		return nil, err
	}

	return e.saveExcelToBuffer(file)
}

// writeEmptyExcelFile creates an empty Excel file
func (e *DBExporter) writeEmptyExcelFile(file *excelize.File) ([]byte, error) {
	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write empty Excel file: %w", err)
	}
	return buffer.Bytes(), nil
}

// prepareExcelHeaders prepares headers for Excel export
func (e *DBExporter) prepareExcelHeaders(data []map[string]interface{}, fields []string) []string {
	headers := fields
	if len(headers) == 0 || (len(headers) == 1 && headers[0] == "*") {
		headers = e.extractSortedHeaders(data[0])
	}
	return headers
}

// extractSortedHeaders extracts and sorts headers from data
func (e *DBExporter) extractSortedHeaders(record map[string]interface{}) []string {
	headers := make([]string, 0, len(record))
	for key := range record {
		headers = append(headers, key)
	}
	// Sort headers for consistent output
	for i := 0; i < len(headers); i++ {
		for j := i + 1; j < len(headers); j++ {
			if headers[i] > headers[j] {
				headers[i], headers[j] = headers[j], headers[i]
			}
		}
	}
	return headers
}

// writeExcelHeaders writes headers to Excel file
func (e *DBExporter) writeExcelHeaders(file *excelize.File, sheetName string, headers []string) error {
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", e.columnName(i))
		if err := file.SetCellValue(sheetName, cell, header); err != nil {
			return fmt.Errorf("failed to set header cell %s: %w", cell, err)
		}
	}
	return nil
}

// writeExcelData writes data rows to Excel file
func (e *DBExporter) writeExcelData(file *excelize.File, sheetName string, data []map[string]interface{}, headers []string) error {
	for rowIndex, record := range data {
		for colIndex, header := range headers {
			cell := fmt.Sprintf("%s%d", e.columnName(colIndex), rowIndex+constants.ExcelRowOffset)
			value := ""
			if val, exists := record[header]; exists && val != nil {
				value = e.formatValue(val)
			}
			if err := file.SetCellValue(sheetName, cell, value); err != nil {
				return fmt.Errorf("failed to set data cell %s: %w", cell, err)
			}
		}
	}
	return nil
}

// saveExcelToBuffer saves Excel file to buffer
func (e *DBExporter) saveExcelToBuffer(file *excelize.File) ([]byte, error) {
	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write Excel file to buffer: %w", err)
	}
	return buffer.Bytes(), nil
}

// formatValue formats interface{} value to string for export
func (e *DBExporter) formatValue(value interface{}) string {
	if value == nil {
		return ""
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return value.(string)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	default:
		return fmt.Sprintf("%v", value)
	}
}

// columnName converts column index to Excel column name (A, B, C, ..., AA, AB, ...)
func (e *DBExporter) columnName(index int) string {
	name := ""
	for index >= 0 {
		name = string(rune('A'+index%constants.ExcelColumnDivisor)) + name
		index = index/constants.ExcelColumnDivisor - 1
	}
	return name
}

// GetRowCount returns the count of rows for given table and query builder
func (e *DBExporter) GetRowCount(tableName string, queryBuilder QueryBuilder) (int64, error) {
	if tableName == "" {
		return 0, errors.New("table name cannot be empty")
	}

	query := e.db.Table(tableName)

	if queryBuilder != nil {
		query = queryBuilder(query)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count rows: %w", err)
	}

	return count, nil
}

// ExportConfigBuilder helps build ExportConfig with method chaining
type ExportConfigBuilder struct {
	config *ExportConfig
}

// NewExportConfigBuilder creates a new ExportConfigBuilder
func NewExportConfigBuilder() *ExportConfigBuilder {
	return &ExportConfigBuilder{
		config: &ExportConfig{
			Fields: make([]string, 0),
			Format: FormatJSON,
		},
	}
}

// TableName sets the table name
func (ecb *ExportConfigBuilder) TableName(tableName string) *ExportConfigBuilder {
	ecb.config.TableName = tableName
	return ecb
}

// Fields sets the fields to export
func (ecb *ExportConfigBuilder) Fields(fields ...string) *ExportConfigBuilder {
	ecb.config.Fields = fields
	return ecb
}

// QueryBuilder sets the query builder function
func (ecb *ExportConfigBuilder) QueryBuilder(builder QueryBuilder) *ExportConfigBuilder {
	ecb.config.QueryBuilder = builder
	return ecb
}

// Format sets the export format
func (ecb *ExportConfigBuilder) Format(format ExportFormat) *ExportConfigBuilder {
	ecb.config.Format = format
	return ecb
}

// Limit sets the record limit
func (ecb *ExportConfigBuilder) Limit(limit int) *ExportConfigBuilder {
	ecb.config.Limit = limit
	return ecb
}

// OrderBy sets the ordering
func (ecb *ExportConfigBuilder) OrderBy(orderBy string) *ExportConfigBuilder {
	ecb.config.OrderBy = orderBy
	return ecb
}

// Build creates the final ExportConfig
func (ecb *ExportConfigBuilder) Build() *ExportConfig {
	// Create a copy to avoid modifications after build
	result := &ExportConfig{
		TableName:    ecb.config.TableName,
		Fields:       make([]string, len(ecb.config.Fields)),
		QueryBuilder: ecb.config.QueryBuilder,
		Format:       ecb.config.Format,
		Limit:        ecb.config.Limit,
		OrderBy:      ecb.config.OrderBy,
	}

	copy(result.Fields, ecb.config.Fields)
	return result
}

// WhereEqual creates a QueryBuilder for equality condition
func WhereEqual(field string, value interface{}) QueryBuilder {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", field), value)
	}
}

// WhereBetween creates a QueryBuilder for between condition
func WhereBetween(field string, start, end interface{}) QueryBuilder {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), start, end)
	}
}

// WhereLike creates a QueryBuilder for like condition
func WhereLike(field, pattern string) QueryBuilder {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s LIKE ?", field), pattern)
	}
}

// WhereIn creates a QueryBuilder for in condition
func WhereIn(field string, values []interface{}) QueryBuilder {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN ?", field), values)
	}
}

// WhereNotNull creates a QueryBuilder for not null condition
func WhereNotNull(field string) QueryBuilder {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IS NOT NULL", field))
	}
}

// WhereIsNull creates a QueryBuilder for is null condition
func WhereIsNull(field string) QueryBuilder {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IS NULL", field))
	}
}

// CombineQueryBuilders combines multiple QueryBuilders with AND logic
func CombineQueryBuilders(builders ...QueryBuilder) QueryBuilder {
	return func(db *gorm.DB) *gorm.DB {
		query := db
		for _, builder := range builders {
			if builder != nil {
				query = builder(query)
			}
		}
		return query
	}
}
