package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
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

var maxLenMasterPassword = 32

// GetMasterPassword show prompt and get the entered password
func (ps *Password) GetMasterPassword() (string, error) {
	fmt.Print("Please enter the master password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if len(string(password)) > maxLenMasterPassword {
		return "", fmt.Errorf("The length of master password shold be less than %d", maxLenMasterPassword)
	}

	return string(password), err
}

// GetDbHelper get a new DbHelper instance
func (ps *Password) GetDbHelper() *DbHelper {
	return &DbHelper{Name: "pw", Driver: "sqlite3"}
}

// Register register a new password
func (ps *Password) Register(masterPassword string) error {

	dbHelper := DbHelper{Name: "pw", Driver: "sqlite3"}
	if err := dbHelper.Open(masterPassword); err != nil {
		return err
	}
	defer dbHelper.Close()

	if !dbHelper.ExistsTable("password") {
		if _, err := createTable(&dbHelper); err != nil {
			return err
		}
	}

	if err := ps.Encrypt(masterPassword); err != nil {
		return err
	}

	_, err := dbHelper.Execute(fmt.Sprintf("insert into password (tag, password) values ('%s', '%s')", ps.Tag, ps.Password))
	if err != nil {
		return err
	}

	return nil
}

// GetPassword get a registered password with the tag.
func (ps *Password) GetPassword(masterPassword string) (string, error) {

	dbHelper := DbHelper{Name: "pw", Driver: "sqlite3"}
	if err := dbHelper.Open(masterPassword); err != nil {
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

	ps.Password = password
	if err := ps.Decrypt(masterPassword); err != nil {
		return "", err
	}

	return ps.Password, nil
}

// GetTagList get the list of tags
func (ps *Password) GetTagList(masterPassword string) ([]string, error) {
	dbHelper := DbHelper{Name: "pw", Driver: "sqlite3"}
	if err := dbHelper.Open(masterPassword); err != nil {
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

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

// Encrypt encrypt password
func (ps *Password) Encrypt(masterPassword string) error {
	block, err := aes.NewCipher([]byte(ps.CreateHash(masterPassword)))
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ps.Password = string(gcm.Seal(nonce, nonce, []byte(ps.Password), nil))

	return nil
}

// Decrypt decrypt password
func (ps *Password) Decrypt(masterPassword string) error {
	key := []byte(ps.CreateHash(masterPassword))
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonceSize := gcm.NonceSize()

	data := []byte(ps.Password)
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}
	ps.Password = string(plaintext)

	return nil
}

// CreateHash create a hash
func (ps *Password) CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
