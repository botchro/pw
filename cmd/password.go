package cmd

// Password an password data set
type Password struct {
	Tag      string
	Password string
}

func createTable(dbHelper *DbHelper) error {
	return dbHelper.CreateTable("password", map[string]string{
		"id": "integer primary key autoincrement",
	})
}

// Register register a new password
func (ps *Password) Register() error {

	dbHelper := DbHelper{Name: "pw", Driver: "sqlite3"}
	if err := dbHelper.Init(); err != nil {
		return err
	}
	defer dbHelper.DB.Close()

	if !dbHelper.ExistsTable("password") {
		if err := createTable(&dbHelper); err != nil {
			return err
		}
	}

	return nil
}
