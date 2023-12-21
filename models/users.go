package models

import (
	"errors"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type User struct {
	ID        uint64 `gorm:"primarykey"`
	Email     string `json:"email" gorm:"unique;size:255"`
	Password  string `json:"password" gorm:"type:text"`
	CreatedAt time.Time
}

func CreateUser(u *User) (*User, error) {
	// ユーザーのバリデーション
	if err := u.Validate(); err != nil {
		return nil, err
	}

	entry := User{
		Email:    u.Email,
		Password: Encrypt(u.Password),
	}
	result := DB.Create(&entry)

	if result.Error != nil {
		return nil, TranslateErrors(result)
	}

	return &entry, nil
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(
			&u.Email,
			validation.Required.Error("メールアドレスは必須入力です"),
			validation.RuneLength(5, 254).Error("メールアドレスは 5~254 文字です"), // RFC 5321 に準拠
			is.Email.Error("メールアドレスの形式を確認してください"),
		),
		validation.Field(
			&u.Password,
			validation.Required.Error("パスワードは必須入力です"),
			validation.RuneLength(8, 20).Error("パスワードは 8~20 文字です"),
			validation.Match(regexp.MustCompile(`[A-Za-z]`)).Error("パスワードには 半角英字 を少なくとも1つ含んで下さい"),
			validation.Match(regexp.MustCompile(`\d`)).Error("パスワードには 半角数字 を少なくとも1つ含んで下さい"),
			validation.Match(regexp.MustCompile(`[!?\\-_]`)).Error("パスワードには !?-_ を少なくとも1つ含んで下さい"),
		),
	)
}

func TranslateErrors(value *gorm.DB) error {
	if err, ok := value.Error.(*pgconn.PgError); ok {
		switch err.Code {
		case pgerrcode.UniqueViolation:
			return errors.New("そのメールアドレスは既に使われています")
		}
	}
	return value.Error
}
