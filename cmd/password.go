package cmd

import "fmt"

// Password an password data set
type Password struct {
	Tag      string
	Password string
}

func createTable(dbHelper *DbHelper) error {
	return dbHelper.Execute(`
	create table password (
		id integer autoincrement primary key,
		tag varchar not null,
		password varchar not null,
		key tag (tag)
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
		if err := createTable(&dbHelper); err != nil {
			return err
		}
	}

	err := dbHelper.Execute(fmt.Sprintf("insert into password (tag, password) values ('%s', '%s')", ps.Tag, ps.Password))
	if err != nil {
		return err
	}

	return nil
}
