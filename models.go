package main

import "github.com/jinzhu/gorm"

// Notes ...
type Notes struct {
	gorm.Model
	Content       string `gorm:"type:text"`
	EncryptionKey string
}
