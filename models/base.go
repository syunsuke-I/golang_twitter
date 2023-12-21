package models

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/driver/postgres"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ErrorMsg struct {
	EmailRequired       string `json:"emailRequired"`
	EmailLength         string `json:"emailLength"`
	EmailFormat         string `json:"emailFormat"`
	PasswordRequired    string `json:"passwordRequired"`
	PasswordLength      string `json:"passwordLength"`
	PasswordAlphabet    string `json:"passwordAlphabet"`
	PasswordNumber      string `json:"passwordNumber"`
	PasswordSpecialChar string `json:"passwordSpecialChar"`
	EmailInUse          string `json:"emailInUse"`
}

var DB *gorm.DB

const (
	host          = "postgres"
	port          = 5432
	user          = "postgres"
	password      = "password"
	dbname        = "postgres"
	tableNameUser = "users"
)

func Init() {
	ConnectionDatabase()
	CreateTables()
}

func ConnectionDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s "+
		"dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connected to database!")
	}
	DB = db
	fmt.Println("connected to db is seceded!")
}

func CreateTables() {
	cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    password TEXT,
    created_at TIMESTAMP
	)`, tableNameUser)

	DB.Exec(cmdU)
}

func LoadConfig(filename string) (ErrorMsg, error) {
	var errMsg ErrorMsg

	file, err := os.Open(filename)
	if err != nil {
		return errMsg, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&errMsg)
	if err != nil {
		return errMsg, err
	}

	return errMsg, nil
}

// 暗号(Hash)化
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}

// 暗号(Hash)と入力された平パスワードの比較
func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
