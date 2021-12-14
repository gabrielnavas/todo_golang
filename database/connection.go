package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const postgresDriver = "postgres"

func MakeConnection(user, host, port, password, dbname, sslmode string) (*sql.DB, error) {
	var DataSourceName = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	db, err := sql.Open(postgresDriver, DataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("[ * ] DB base connected")
	return db, nil
}
