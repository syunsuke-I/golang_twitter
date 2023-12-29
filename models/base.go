package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	LoginError          string `json:"loginError"`
	InactiveAccount     string `json:"inactiveAccount"`
	ServerError         string `json:"serverError"`
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
func Encrypt(password string) (cryptext string) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

// 暗号(Hash)と入力された平パスワードの比較
func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// メールアドレスからtokenを作成
func GenerateTokenFromEmail(email string) string {
	hash := sha256.Sum256([]byte(email))
	return hex.EncodeToString(hash[:])
}
