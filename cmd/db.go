package cmd

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// DbHelper database
type DbHelper struct {
	Name   string
	Driver string
	DB     *sql.DB
}

// Init initialize
func (dbHelper *DbHelper) Init() error {
	db, err := sql.Open(dbHelper.Driver, dbHelper.GetDbFilePath())
	if err != nil {
		return err
	}
	dbHelper.DB = db
	return nil
}

// Close close the database connection
func (dbHelper *DbHelper) Close() {
	if dbHelper.DB == nil {
		return
	}
	dbHelper.DB.Close()
}

// ExistsTable check exists the table
func (dbHelper *DbHelper) ExistsTable(tableName string) bool {
	sql := fmt.Sprintf("select * from %s limit 1", tableName)
	_, err := dbHelper.DB.Exec(sql)
	return err == nil
}

// CreateTable create a new table
func (dbHelper *DbHelper) CreateTable(tableName string, columns map[string]string) error {
	var sql string
	sql += fmt.Sprintf("create table `%s` (", tableName)
	for k, v := range columns {
		sql += fmt.Sprintf("%s %s", k, v)
	}
	sql += ")"

	_, err := dbHelper.DB.Exec(sql)
	return err
}

// GetDbFilePath get the file path to the db file
func (dbHelper *DbHelper) GetDbFilePath() string {
	return fmt.Sprintf("./%s.db", dbHelper.Name)
}
