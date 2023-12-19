package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID        uint64 `gorm:"primarykey"`
	Email     string `json:"email" gorm:"size:255"`
	Password  string `json:"password" gorm:"type:text"`
	CreatedAt time.Time
}

func CreateUser(u *User) (*User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	entry := User{
		Email:    u.Email,
		Password: Encrypt(u.Password),
	}
	if err := DB.Create(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(
			&u.Email,
			validation.Required.Error("メールアドレスは必須入力です"),
			validation.RuneLength(5, 254).Error("メールアドレスは 5~254 文字です"), // RFC 5321 に準拠
			is.Email.Error("メールアドレスを入力して下さい"),
		),
		validation.Field(
			&u.Password,
			validation.Required.Error("パスワードは必須入力です"),
			validation.RuneLength(8, 20).Error("パスワードは 8=20 文字です"),
		),
	)
}
