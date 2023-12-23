package models

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
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
	PasswordMixedCase   string `json:"passwordMixedCase"`
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
