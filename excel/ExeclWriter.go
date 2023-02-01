package excel

import (
	"fmt"
	"strconv"

	"github.com/Bit0r/schema2file/model"
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
)

type ExcelWriter struct {
	db          *sqlx.DB
	f           *excelize.File
	out         string
	headerStyle int
	colStyle    int
}

func New(db *sqlx.DB, output string) (*ExcelWriter, error) {
	f := excelize.NewFile()
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return nil, err
	}

	colStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return nil, err
	}

	return &ExcelWriter{db, f, output, headerStyle, colStyle}, nil
}

func (ew *ExcelWriter) WriteColumns(table string) error {
	db, f := ew.db, ew.f

	sheetName := table
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	f.SetColStyle(sheetName, "A:G", ew.colStyle)
	f.SetRowStyle(sheetName, 1, 1, ew.headerStyle)
	f.SetSheetRow(sheetName, "A1", &[]any{"列名", "类型", "允许空值", "键", "默认值", "额外信息", "说明"})

	cols, err := db.Queryx("SHOW FULL COLUMNS FROM " + table)
	if err != nil {
		return err
	}
	defer cols.Close()

	row := 2
	for cols.Next() {
		col := model.Column{}
		if err := cols.StructScan(&col); err != nil {
			return err
		}

		col.Null = model.Translate[col.Null]
		col.Key = model.Translate[col.Key]

		if !col.Default.Valid {
			col.Default.String = "NULL"
			col.Default.Valid = true
		}

		f.SetSheetRow(sheetName, "A"+strconv.Itoa(row), &[]any{col.Field, col.Type, col.Null, col.Key, col.Default.String, col.Extra, col.Comment})
		row++
	}
	return nil
}

func (ew *ExcelWriter) Close() error {
	f, out := ew.f, ew.out
	f.DeleteSheet("Sheet1")

	err := f.SaveAs(out)

	if closeErr := f.Close(); closeErr != nil {
		fmt.Println(closeErr)
	}

	return err
}
