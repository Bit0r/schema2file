package main

type ColumnsWriter interface {
	WriteColumns(table string) error
	Close() error
}
