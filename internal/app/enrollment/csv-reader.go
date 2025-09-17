package enrollment

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"shedoo-backend/internal/models"

	"github.com/xuri/excelize/v2"
)

func readCSV(filePath string) ([]models.Enrollment, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fmt.Println("File raw content:")
	fmt.Println(string(data))
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read()

	var enrollments []models.Enrollment
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		enrollments = append(enrollments, models.Enrollment{
			StudentCode: record[0],
			CourseCode:  record[4],
			LecSection:  record[7],
			LabSection:  record[8],
			Semester:    record[9],
			Year:        record[10],
		})
	}

	return enrollments, nil
}

func readXLSX(filePath string) ([]models.Enrollment, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var enrollments []models.Enrollment
	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 11 {
			continue
		}

		for i := range row {
			row[i] = strings.TrimSpace(row[i])
		}

		enrollments = append(enrollments, models.Enrollment{
			StudentCode: row[0],
			CourseCode:  row[4],
			LecSection:  row[7],
			LabSection:  row[8],
			Semester:    row[9],
			Year:        row[10],
		})
	}

	return enrollments, nil
}
