package models

import (
	"time"
)

type User struct {
	ID        uint64 `gorm:"primarykey"`
	Email     string `gorm:"size:255"`
	Password  string `gorm:"type:text"`
	CreatedAt time.Time
}

func CreateUser(email string, password string) *User {
	pwd := Encrypt(password)
	entry := User{
		Email:    email,
		Password: pwd,
	}
	DB.Create(&entry)
	return &entry
}
