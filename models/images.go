package models

import (
	"time"
)

type Image struct {
	ID        uint64    `gorm:"primarykey"`
	TweetID   uint64    `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey"`
	ImgUrl    string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (p *Repository) CreateImage(t *Image) (*Image, error) {

	entry := Image{
		TweetID: t.TweetID,
		ImgUrl:  t.ImgUrl,
	}

	result := p.DB.Create(&entry)
	if result.Error != nil {
		return nil, TranslateErrors(result)
	}

	return &entry, nil
}
