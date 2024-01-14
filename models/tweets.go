package models

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Tweet struct {
	ID        uint64    `gorm:"primarykey"`
	UserID    uint64    `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type TweetsDto struct {
	Page   Page    `json:"page"`
	Tweets []Tweet `json:"tweets"`
}

var query_pagination Pagination

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

func (p *Repository) TweetsFind(id uint64) *[]Tweet {
	var tweets []Tweet
	p.DB.Order("updated_at desc").Where("user_id = ?", id).Find(&tweets)
	return &tweets
}

func (p *Repository) GetTweets(context *gin.Context, uid uint64) (TweetsDto, error) {
	var tweets []Tweet

	totalElements := p.DB.Find(&tweets).RowsAffected
	var page Page = ConvertContextAndTotalElementsToPage(context, int(totalElements))

	if err := p.DB.Scopes(query_pagination.Pagination(page)).Order("updated_at desc").Where("user_id = ?", uid).Find(&tweets).Error; err != nil {
		return TweetsDto{}, err
	}

	return TweetsDto{Page: page, Tweets: tweets}, nil
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
