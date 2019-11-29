package mysqlToolkit

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// CreateTableIfNotExists create table automatically,and return table name
func CreateTableIfNotExists(dsn string, v interface{}) (string, error) {
	db, e := sql.Open("mysql", dsn)
	if e != nil {
		log.Println(e)
		return "", e
	}
	
}
