package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type User struct {
	ID              uint64 `gorm:"primarykey"`
	Email           string `json:"email" gorm:"unique;size:255"`
	Password        string `json:"password" gorm:"type:text"`
	ActivationToken string `gorm:"size:64"`
	IsActive        bool   `gorm:"default:false"`
	CreatedAt       time.Time
}

func (p *Repository) CreateUser(u *User) (*User, error) {
	// ユーザーのバリデーション
	if err := u.Validate(); err != nil {
		return nil, err
	}

	entry := User{
		Email:           u.Email,
		Password:        Encrypt(u.Password),
		ActivationToken: u.ActivationToken,
	}

	result := p.DB.Create(&entry)
	if result.Error != nil {
		return nil, TranslateErrors(result)
	}

	return &entry, nil
}

func (p *Repository) FindUserByActivationToken(token string) (*User, error) {
	var user User
	result := p.DB.Where("activation_token = ?", token).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // トークンに該当するユーザーが見つからない場合
		}
		return nil, result.Error
	}
	return &user, nil
}

func (p *Repository) FindUserByEmail(email string) (*User, error) {
	var user User
	result := p.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // メールアドレスに該当するユーザーが見つからない場合
		}
		return nil, result.Error
	}
	return &user, nil
}

func (p *Repository) ActivateUser(u *User) error {
	if u == nil {
		return errors.New("provided user is nil")
	}
	u.IsActive = true // ユーザーをアクティブに設定
	return p.DB.Save(u).Error
}

func (u User) Validate() error {
	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	return validation.ValidateStruct(&u,
		validation.Field(
			&u.Email,
			validation.Required.Error(errMsg.EmailRequired),
			is.Email.Error(errMsg.EmailFormat),
		),
		validation.Field(
			&u.Password,
			validation.Required.Error(errMsg.PasswordRequired),
			validation.RuneLength(8, 20).Error(errMsg.PasswordLength),
			validation.Match(regexp.MustCompile(`[A-Za-z]`)).Error(errMsg.PasswordAlphabet),
			validation.Match(regexp.MustCompile(`\d`)).Error(errMsg.PasswordNumber),
			validation.Match(regexp.MustCompile(`[!?\\-_]`)).Error(errMsg.PasswordSpecialChar),
			validation.Match(regexp.MustCompile(`[A-Z]`)).Error(errMsg.PasswordMixedCase),
			validation.Match(regexp.MustCompile(`[a-z]`)).Error(errMsg.PasswordMixedCase),
		),
	)
}

func TranslateErrors(value *gorm.DB) error {
	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	if err, ok := value.Error.(*pgconn.PgError); ok {
		switch err.Code {
		case pgerrcode.UniqueViolation:
			return errors.New(errMsg.EmailInUse)
		}
	}
	return value.Error
}
