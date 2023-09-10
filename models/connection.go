package models

import "gorm.io/gorm"

type Connection struct {
	gorm.Model
	UserID     uint   `json:"userId"`
	Label      string `json:"label"`
	ConnString string `json:"connectionString"`
}

func DeleteConnection(db *gorm.DB, id uint) error {
	return db.Delete(&Connection{}, id).Error
}
