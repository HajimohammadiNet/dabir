package imports

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hajimohammadinet/dabir/internal/shared/dateutil"
	"github.com/xuri/excelize/v2"
)

//const incomingSheetName = "دریافتی"

type LetterExcelParser struct{}

type ParseLettersResult struct {
	DetectedColumns map[string]string
	Rows            []ImportedLetterRow
	Errors          []ImportErrorDTO
	MaxLetterNumber *int64
	TotalRows       int
	ValidRows       int
	InvalidRows     int
}

var columnAliases = map[string][]string{
	"letter_number": {
		"letter_number",
		"number",
		"no",
		"شماره نامه",
		"شماره",
	},
	"title": {
		"title",
		"subject",
		"عنوان",
		"عنوان نامه",
		"موضوع",
	},
	"letter_date": {
		"letter_date",
		"date",
		"تاریخ",
		"تاریخ نامه",
	},
	"sender": {
		"sender",
		"from",
		"فرستنده",
		"ارسال کننده",
	},
	"receiver": {
		"receiver",
		"to",
		"گیرنده",
		"دریافت کننده",
		"مقصد",
	},
}

func NewLetterExcelParser() *LetterExcelParser {
	return &LetterExcelParser{}
}

func (p *LetterExcelParser) Parse(fileName string, reader io.Reader) (*ParseLettersResult, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".xlsx" {
		return nil, errors.New("only .xlsx files are supported")
	}

	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer f.Close()

	sheetName, rows, columnIndex, detectedColumns, err := findImportableSheet(f)
	if err != nil {
		return nil, err
	}

	_ = sheetName

	result := &ParseLettersResult{
		DetectedColumns: detectedColumns,
		Rows:            make([]ImportedLetterRow, 0),
		Errors:          make([]ImportErrorDTO, 0),
	}

	required := []string{"letter_number", "letter_date"}
	for _, field := range required {
		if _, ok := columnIndex[field]; !ok {
			result.Errors = append(result.Errors, ImportErrorDTO{
				Row:     1,
				Field:   field,
				Message: "required column not found",
			})
		}
	}

	if len(result.Errors) > 0 {
		result.InvalidRows = len(rows) - 1
		result.TotalRows = len(rows) - 1
		return result, nil
	}

	for i := 1; i < len(rows); i++ {
		excelRowNumber := i + 1
		row := rows[i]

		if isEmptyRow(row) {
			continue
		}

		result.TotalRows++

		parsedRow, rowErrors := parseLetterRow(excelRowNumber, row, columnIndex)
		if len(rowErrors) > 0 {
			result.InvalidRows++
			result.Errors = append(result.Errors, rowErrors...)
			continue
		}

		result.ValidRows++
		result.Rows = append(result.Rows, parsedRow)
	}

	return result, nil
}

func findImportableSheet(f *excelize.File) (string, [][]string, map[string]int, map[string]string, error) {
	sheets := f.GetSheetList()

	if len(sheets) == 0 {
		return "", nil, nil, nil, errors.New("excel file has no sheets")
	}

	var lastError error

	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			lastError = fmt.Errorf("failed to read excel rows from sheet %q: %w", sheet, err)
			continue
		}

		if len(rows) < 2 {
			continue
		}

		header := rows[0]
		columnIndex, detectedColumns := detectColumns(header)

		if hasRequiredImportColumns(columnIndex) {
			return sheet, rows, columnIndex, detectedColumns, nil
		}
	}

	if lastError != nil {
		return "", nil, nil, nil, lastError
	}

	return "", nil, nil, nil, errors.New("no importable sheet found: required columns are letter number and letter date")
}

func hasRequiredImportColumns(columnIndex map[string]int) bool {
	_, hasNumber := columnIndex["letter_number"]
	_, hasDate := columnIndex["letter_date"]

	return hasNumber && hasDate
}

func detectColumns(header []string) (map[string]int, map[string]string) {
	columnIndex := make(map[string]int)
	detectedColumns := make(map[string]string)

	normalizedHeader := make([]string, len(header))
	for i, h := range header {
		normalizedHeader[i] = normalizeHeader(h)
	}

	for field, aliases := range columnAliases {
		for i, normalized := range normalizedHeader {
			for _, alias := range aliases {
				if normalized == normalizeHeader(alias) {
					columnIndex[field] = i
					detectedColumns[field] = strings.TrimSpace(header[i])
					break
				}
			}

			if _, ok := columnIndex[field]; ok {
				break
			}
		}
	}

	return columnIndex, detectedColumns
}

func normalizeHeader(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "‌", "")
	value = strings.ReplaceAll(value, "\u200c", "")
	return value
}

func parseLetterRow(rowNumber int, row []string, columnIndex map[string]int) (ImportedLetterRow, []ImportErrorDTO) {
	errorsList := make([]ImportErrorDTO, 0)

	displayLetterNumber := strings.TrimSpace(getCell(row, columnIndex["letter_number"]))
	dateRaw := strings.TrimSpace(getCell(row, columnIndex["letter_date"]))

	title := getOptionalCell(row, columnIndex, "title")
	sender := getOptionalCell(row, columnIndex, "sender")
	receiver := getOptionalCell(row, columnIndex, "receiver")

	if displayLetterNumber == "" {
		errorsList = append(errorsList, ImportErrorDTO{
			Row:     rowNumber,
			Field:   "display_letter_number",
			Message: "letter number is required",
		})
	}

	letterDate, err := parseExcelDate(dateRaw)
	if err != nil {
		errorsList = append(errorsList, ImportErrorDTO{
			Row:     rowNumber,
			Field:   "letter_date",
			Message: err.Error(),
		})
	}

	if title == "" {
		title = "بدون موضوع"
	}

	if sender == "" {
		sender = "-"
	}

	if receiver == "" {
		receiver = "-"
	}

	if len(errorsList) > 0 {
		return ImportedLetterRow{}, errorsList
	}

	parsedGregorianDate, err := time.Parse("2006-01-02", letterDate)
	letterDateJalali := ""
	if err == nil {
		letterDateJalali = dateutil.ToJalaliString(parsedGregorianDate)
	}

	return ImportedLetterRow{
		RowNumber:           rowNumber,
		LetterNumber:        int64(rowNumber),
		DisplayLetterNumber: displayLetterNumber,
		Title:               title,
		LetterDate:          letterDate,
		LetterDateJalali:    letterDateJalali,
		Sender:              sender,
		Receiver:            receiver,
	}, errorsList
}

func getOptionalCell(row []string, columnIndex map[string]int, field string) string {
	index, ok := columnIndex[field]
	if !ok {
		return ""
	}

	return getCell(row, index)
}

func getCell(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}

	return strings.TrimSpace(row[index])
}

func parseExcelDate(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("letter date is required")
	}

	value = dateutil.NormalizeDigits(value)
	value = strings.ReplaceAll(value, "\\", "/")
	value = strings.ReplaceAll(value, "-", "/")
	value = strings.ReplaceAll(value, ".", "/")

	if parsed, err := dateutil.ParseOfficialDate(value); err == nil {
		return parsed.Format("2006-01-02"), nil
	}

	if parsed, ok := parseCompactJalaliDate(value); ok {
		return parsed.Format("2006-01-02"), nil
	}

	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"02/01/2006",
		"02-01-2006",
		"1/2/2006",
		"2006/1/2",
	}

	for _, layout := range formats {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed.Format("2006-01-02"), nil
		}
	}

	if serial, err := strconv.ParseFloat(value, 64); err == nil {
		parsed, err := excelize.ExcelDateToTime(serial, false)
		if err == nil {
			return parsed.Format("2006-01-02"), nil
		}
	}

	return "", errors.New("invalid date format")
}

func parseCompactJalaliDate(value string) (time.Time, bool) {
	cleaned := strings.TrimSpace(value)
	cleaned = strings.ReplaceAll(cleaned, "/", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	if len(cleaned) != 6 && len(cleaned) != 8 {
		return time.Time{}, false
	}

	if _, err := strconv.Atoi(cleaned); err != nil {
		return time.Time{}, false
	}

	var year int
	var month int
	var day int
	var err error

	if len(cleaned) == 6 {
		year, err = strconv.Atoi(cleaned[0:2])
		if err != nil {
			return time.Time{}, false
		}

		if year >= 90 {
			year += 1300
		} else {
			year += 1400
		}

		month, err = strconv.Atoi(cleaned[2:4])
		if err != nil {
			return time.Time{}, false
		}

		day, err = strconv.Atoi(cleaned[4:6])
		if err != nil {
			return time.Time{}, false
		}
	} else {
		year, err = strconv.Atoi(cleaned[0:4])
		if err != nil {
			return time.Time{}, false
		}

		month, err = strconv.Atoi(cleaned[4:6])
		if err != nil {
			return time.Time{}, false
		}

		day, err = strconv.Atoi(cleaned[6:8])
		if err != nil {
			return time.Time{}, false
		}
	}

	if month < 1 || month > 12 || day < 1 || day > 31 {
		return time.Time{}, false
	}

	gYear, gMonth, gDay, err := dateutil.JalaliToGregorian(year, month, day)
	if err != nil {
		return time.Time{}, false
	}

	return time.Date(gYear, time.Month(gMonth), gDay, 0, 0, 0, 0, time.UTC), true
}

func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}

	return true
}
