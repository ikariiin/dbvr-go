package models

import "gorm.io/gorm"

type Connection struct {
	gorm.Model
	UserID     uint   `json:"user-id"`
	ConnString string `json:"connection-string"`
}
