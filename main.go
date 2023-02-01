package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/Bit0r/schema2file/excel"
	"github.com/Bit0r/schema2file/markdown"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Column struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
	Comment string
}

func main() {
	host := flag.String("h", "localhost", "host name")
	user := flag.String("u", "root", "user name")
	password := flag.String("p", "", "password")
	port := flag.Int("P", 3306, "port number")
	database := flag.String("B", "", "database name")
	outputfile := flag.String("o", "schema.md", "output file name")
	flag.Parse()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *password, *host, *port, *database)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	tables, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatalln(err)
	}

	w, err := getWriter(db, *outputfile)
	if err != nil {
		log.Fatalln(err)
	}
	defer w.Close()

	for tables.Next() {
		var table string
		if err := tables.Scan(&table); err != nil {
			log.Fatalln(err)
		}
		if err := w.WriteColumns(table); err != nil {
			log.Fatalln(err)
		}
	}
}

func getWriter(db *sqlx.DB, output string) (ColumnsWriter, error) {
	switch filepath.Ext(output) {
	case ".md":
		return markdown.New(db, output)
	case ".xlsx":
		return excel.New(db, output)
	default:
		return nil, fmt.Errorf("unsupported output file type: %s", output)
	}
}
