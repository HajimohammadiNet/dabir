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

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("excel file does not contain any sheets")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read excel rows: %w", err)
	}

	if len(rows) < 2 {
		return nil, errors.New("excel file must contain header row and at least one data row")
	}

	header := rows[0]
	columnIndex, detectedColumns := detectColumns(header)

	result := &ParseLettersResult{
		DetectedColumns: detectedColumns,
		Rows:            make([]ImportedLetterRow, 0),
		Errors:          make([]ImportErrorDTO, 0),
	}

	required := []string{"letter_number", "title", "letter_date", "sender", "receiver"}
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

	var maxNumber int64

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

		if parsedRow.LetterNumber > maxNumber {
			maxNumber = parsedRow.LetterNumber
		}

		result.ValidRows++
		result.Rows = append(result.Rows, parsedRow)
	}

	if maxNumber > 0 {
		result.MaxLetterNumber = &maxNumber
	}

	return result, nil
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
	return value
}

func parseLetterRow(rowNumber int, row []string, columnIndex map[string]int) (ImportedLetterRow, []ImportErrorDTO) {
	errorsList := make([]ImportErrorDTO, 0)

	letterNumberRaw := getCell(row, columnIndex["letter_number"])
	title := strings.TrimSpace(getCell(row, columnIndex["title"]))
	dateRaw := strings.TrimSpace(getCell(row, columnIndex["letter_date"]))
	sender := strings.TrimSpace(getCell(row, columnIndex["sender"]))
	receiver := strings.TrimSpace(getCell(row, columnIndex["receiver"]))

	letterNumber, err := parseLetterNumber(letterNumberRaw)
	if err != nil {
		errorsList = append(errorsList, ImportErrorDTO{
			Row:     rowNumber,
			Field:   "letter_number",
			Message: err.Error(),
		})
	}

	letterDate, err := parseExcelDate(dateRaw)
	parsedGregorianDate, _ := time.Parse("2006-01-02", letterDate)
	letterDateJalali := dateutil.ToJalaliString(parsedGregorianDate)
	if err != nil {
		errorsList = append(errorsList, ImportErrorDTO{
			Row:     rowNumber,
			Field:   "letter_date",
			Message: err.Error(),
		})
	}

	if title == "" {
		errorsList = append(errorsList, ImportErrorDTO{
			Row:     rowNumber,
			Field:   "title",
			Message: "title is required",
		})
	}

	if sender == "" {
		errorsList = append(errorsList, ImportErrorDTO{
			Row:     rowNumber,
			Field:   "sender",
			Message: "sender is required",
		})
	}

	if receiver == "" {
		errorsList = append(errorsList, ImportErrorDTO{
			Row:     rowNumber,
			Field:   "receiver",
			Message: "receiver is required",
		})
	}

	if len(errorsList) > 0 {
		return ImportedLetterRow{}, errorsList
	}

	return ImportedLetterRow{
		RowNumber:        rowNumber,
		LetterNumber:     letterNumber,
		Title:            title,
		LetterDate:       letterDate,
		LetterDateJalali: letterDateJalali,
		Sender:           sender,
		Receiver:         receiver,
	}, nil
}

func getCell(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}

	return strings.TrimSpace(row[index])
}

func parseLetterNumber(value string) (int64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, errors.New("letter number is required")
	}

	value = strings.ReplaceAll(value, ",", "")
	value = strings.ReplaceAll(value, " ", "")

	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.New("letter number must be numeric")
	}

	if number <= 0 {
		return 0, errors.New("letter number must be greater than zero")
	}

	return number, nil
}

func parseExcelDate(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("letter date is required")
	}

	if parsed, err := dateutil.ParseOfficialDate(value); err == nil {
		return parsed.Format("2006-01-02"), nil
	}

	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"02/01/2006",
		"02-01-2006",
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

	return "", errors.New("invalid date format, expected Jalali YYYY/MM/DD")
}

func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}

	return true
}
