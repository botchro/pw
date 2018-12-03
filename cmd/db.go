package cmd

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DbHelper database
type DbHelper struct {
	Name   string
	Driver string
	DB     *sql.DB
}

// ExistsDB check exists the database file
func (dbHelper *DbHelper) ExistsDB() bool {
	_, err := os.Stat(dbHelper.GetDbFilePath(""))
	return err == nil
}

// Open open a database connection
func (dbHelper *DbHelper) Open(masterPassword string) error {
	db, err := sql.Open(dbHelper.Driver, dbHelper.Name)
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

// Execute execute query
func (dbHelper *DbHelper) Execute(sql string) (sql.Result, error) {
	return dbHelper.DB.Exec(sql)
}

// GetRow get a single row
func (dbHelper *DbHelper) GetRow(sql string, args ...interface{}) *sql.Row {
	return dbHelper.DB.QueryRow(sql, args...)
}

// GetRows get rows
func (dbHelper *DbHelper) GetRows(sql string, args ...interface{}) (*sql.Rows, error) {
	return dbHelper.DB.Query(sql, args...)
}

// GetDbFilePath get the file path to the db file
func (dbHelper *DbHelper) GetDbFilePath(connectionString string) string {
	if len(connectionString) > 0 {
		return fmt.Sprintf("./%s.db?%s", dbHelper.Name, connectionString)
	}
	return fmt.Sprintf("./%s.db", dbHelper.Name)
}
