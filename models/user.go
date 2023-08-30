package models

import (
	"errors"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string       `gorm:"size:255;not null;unique" json:"username"`
	Email       string       `gorm:"size:255;not null" json:"email"`
	Password    string       `gorm:"size:255,not null" json:"-"`
	Connections []Connection `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"connections"`
}

func GetUserById(uid uint, db *gorm.DB) (User, error) {
	var user User

	if err := db.Model(&User{}).Where("id=?", uid).Take(&user).Error; err != nil {
		return user, errors.New("User not found")
	}

	return user, nil
}

func GetUserConnections(db *gorm.DB, user User) ([]Connection, error) {
	var connections []Connection
	var preloadedUser User
	err := db.Model(&User{}).Preload("Connections").Where("id=?", user.ID).First(&preloadedUser).Error

	if err != nil {
		return connections, err
	}

	connections = preloadedUser.Connections
	return connections, nil
}

func CreateUserConnection(db *gorm.DB, user User, connStr string) error {
	conn := Connection{ConnString: connStr, UserID: user.ID}
	if err := db.Create(&conn).Error; err != nil {
		return err
	}

	userConnections, err := GetUserConnections(db, user)
	if err != nil {
		return err
	}

	user.Connections = append(userConnections, conn)

	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))

	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
