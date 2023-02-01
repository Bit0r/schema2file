package markdown

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Bit0r/schema2file/model"
	"github.com/jmoiron/sqlx"
)

type MarkdownWriter struct {
	db *sqlx.DB
	w  *bufio.Writer
	f  *os.File
}

func New(db *sqlx.DB, output string) (*MarkdownWriter, error) {
	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	w := bufio.NewWriter(f)
	return &MarkdownWriter{db, w, f}, nil
}

func (mw *MarkdownWriter) WriteColumns(table string) error {
	db, w := mw.db, mw.w

	w.WriteString(fmt.Sprintf("## %s\n\n", table))
	w.WriteString("|列名|类型|允许空值|键|默认值|额外信息|说明|\n")
	w.WriteString("|:---:|:---:|:---:|:---:|:---:|:---:|:---:|\n")

	cols, err := db.Queryx("SHOW FULL COLUMNS FROM " + table)
	if err != nil {
		return err
	}
	defer cols.Close()

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

		w.WriteString(
			fmt.Sprintf("|%s|%s|%s|%s|%s|%s|%s|\n", col.Field, col.Type, col.Null, col.Key, col.Default.String, col.Extra, col.Comment))
	}

	w.WriteString("\n\n")
	return nil
}

func (mw *MarkdownWriter) Close() error {
	w, f := mw.w, mw.f
	w.Flush()
	return f.Close()
}
