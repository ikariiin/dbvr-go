package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() (*gorm.DB, error) {
	dbUrl := fmt.Sprint(os.Getenv("DB_URL"))

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	if err = db.AutoMigrate(&User{}, &Connection{}); err != nil {
		log.Println(err.Error())
	}

	return db, err
}
