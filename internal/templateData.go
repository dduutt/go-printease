package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/xuri/excelize/v2"
)

var tpdColl = mongox.NewCollection[TemplateData](DBClient.NewDatabase("printease"), "template_datas")

type TemplateData map[string]string

func LoadTemplateDatas(id string, path string) (insertIDs []any, errs []error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, []error{err}
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("close file error:", err)
		}
	}()

	sheetName := f.GetSheetList()[0]
	if sheetName == "" {
		return nil, []error{fmt.Errorf("sheet name is empty")}
	}

	rows, err := f.Rows(sheetName)
	if err != nil {
		return nil, []error{err}
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("close rows error:", err)
		}
	}()

	var headers []string
	if rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, []error{fmt.Errorf("failed to read header row: %v", err)}
		}
		if len(row) == 0 {
			return nil, []error{fmt.Errorf("header row is empty")}
		}
		headers = row
	}

	batchSize := 1000
	var batch []*TemplateData
	rowIdx := 1

	for rows.Next() {
		rowIdx++
		row, err := rows.Columns()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read row %d: %v", rowIdx, err))
			continue
		}

		data := make(TemplateData)
		data["_template_id"] = id
		for i, value := range row {
			if i < len(headers) {
				data[headers[i]] = value
			}
		}
		batch = append(batch, &data)

		if len(batch) >= batchSize {
			inserted, batchErrs := insertBatch(batch, rowIdx-batchSize+1, rowIdx)
			insertIDs = append(insertIDs, inserted...)
			errs = append(errs, batchErrs...)
			batch = nil
		}
	}

	if len(batch) > 0 {
		inserted, batchErrs := insertBatch(batch, rowIdx-len(batch)+1, rowIdx)
		insertIDs = append(insertIDs, inserted...)
		errs = append(errs, batchErrs...)
	}

	return
}

func insertBatch(batch []*TemplateData, startRow, endRow int) (insertIDs []any, errs []error) {
	result, err := tpdColl.Creator().InsertMany(context.Background(), batch)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to save rows %d-%d: %v", startRow, endRow, err))
		return
	}
	insertIDs = append(insertIDs, result.InsertedIDs...)
	return
}
