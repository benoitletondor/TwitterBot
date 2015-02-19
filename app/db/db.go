package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func Init(user string, pass string, schema string) (*sql.DB, error) {
	// Init Mysql DB
	dbLink, err := sql.Open("mysql", user+":"+pass+"@/"+schema+"?parseTime=True")
	if err != nil {
		return nil, err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = dbLink.Ping()
	if err != nil {
		return nil, err
	}

	// Set up global var
	database = dbLink

	return database, nil
}
