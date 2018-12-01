package cmd

import (
	"database/sql"
	"fmt"
)

// Password an password data set
type Password struct {
	Tag      string
	Password string
}

func createTable(dbHelper *DbHelper) (sql.Result, error) {
	return dbHelper.Execute(`
	create table password (
		id integer primary key autoincrement,
		tag text not null,
		password text not null
	);
	`)
}

// Register register a new password
func (ps *Password) Register() error {

	dbHelper := DbHelper{Name: "pw", Driver: "sqlite3"}
	if err := dbHelper.Init(); err != nil {
		return err
	}
	defer dbHelper.Close()

	if !dbHelper.ExistsTable("password") {
		if _, err := createTable(&dbHelper); err != nil {
			return err
		}
	}

	_, err := dbHelper.Execute(fmt.Sprintf("insert into password (tag, password) values ('%s', '%s')", ps.Tag, ps.Password))
	if err != nil {
		return err
	}

	return nil
}

// GetPassword get a registered password with the tag.
func (ps *Password) GetPassword() (string, error) {

	dbHelper := DbHelper{Name: "pw", Driver: "sqlite3"}
	if err := dbHelper.Init(); err != nil {
		return "", err
	}
	defer dbHelper.Close()

	if !dbHelper.ExistsTable("password") {
		return "", fmt.Errorf("No password registed")
	}

	row := dbHelper.GetRow(`select password from password where tag = ?`, ps.Tag)
	var password string
	if err := row.Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

// GetTagList get the list of tags
func (ps *Password) GetTagList() ([]string, error) {
	dbHelper := DbHelper{Name: "pw", Driver: "sqlite3"}
	if err := dbHelper.Init(); err != nil {
		return nil, err
	}
	defer dbHelper.Close()

	if !dbHelper.ExistsTable("password") {
		return nil, fmt.Errorf("No password registed")
	}

	rows, err := dbHelper.GetRows(`select tag from password`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
