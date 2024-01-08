package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Tweet struct {
	ID        uint64    `gorm:"primarykey"`
	UserID    uint64    `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (p *Repository) CreateTweet(t *Tweet) (*Tweet, error) {

	// ツイートのバリデーション
	if err := t.Validate(); err != nil {
		return nil, err
	}

	entry := Tweet{
		UserID:  t.UserID,
		Content: t.Content,
	}

	result := p.DB.Create(&entry)
	if result.Error != nil {
		return nil, TranslateErrors(result)
	}

	return &entry, nil
}

func (t Tweet) Validate() error {
	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	return validation.ValidateStruct(&t,
		validation.Field(
			&t.Content,
			validation.Required.Error(errMsg.TweetRequired),
			validation.RuneLength(1, 140).Error(errMsg.TweetLength),
		),
	)
}
