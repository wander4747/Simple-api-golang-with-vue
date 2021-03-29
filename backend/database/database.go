package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Connection open connection database
func Connection() (*sql.DB, error) {
	stringConnection := "root:@/users?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", stringConnection)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
