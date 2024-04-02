package database

import (
	"database/sql"
	"hagavi-blog-go/config"
	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connectString := config.Config("SQL_CONNECT")
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		panic(err)
	}
	return db
}

