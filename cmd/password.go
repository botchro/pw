package cmd

// Password an password data set
type Password struct {
	Tag      string
	Password string
}

// Register register a new password
func (ps *Password) Register() error {
	return nil
}
