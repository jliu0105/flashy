package common

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewMysqlConn() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:flashy@tcp(127.0.0.1:3306)/flashy?charset=utf8")
	return
}

// get ONE result
func GetResultRow(rows *sql.Rows) map[string]string {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]string)
	for rows.Next() {
		// store the row data into record dict
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				record[columns[i]] = string(v)
			}
		}
	}
	return record
}

// get all results
func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	// return all columns
	columns, _ := rows.Columns()
	//This represents the values of all columns in a row, represented by []byte
	vals := make([][]byte, len(columns))
	//Here represents a line of filling data
	scans := make([]interface{}, len(columns))
	//scans here use vals, put data into []bytes
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		// input data
		rows.Scan(scans...)
		row := make(map[string]string)
		// put data inside vals into row
		for k, v := range vals {
			key := columns[k]
			// convert []byte into string
			row[key] = string(v)
		}
		result[i] = row
		i++
	}
	return result
}
