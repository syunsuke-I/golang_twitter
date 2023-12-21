package models

import (
	"errors"
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
			validation.RuneLength(8, 20).Error("パスワードは 8=20 文字です"),
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
