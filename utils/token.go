package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ikariiin/dbvr-go/models"
	"gorm.io/gorm"
)

func GenerateToken(user models.User) (string, error) {
	tokenLifeSpan, err := strconv.Atoi(os.Getenv("JWT_TOKEN_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifeSpan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SIGN")))
}

func getTokenFromRequestContext(ctx *gin.Context) string {
	bearerToken := ctx.Request.Header.Get("Authorization")
	split := strings.Split(bearerToken, " ")

	if len(split) == 2 {
		return split[1]
	}

	return ""
}

func GetToken(ctx *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequestContext(ctx)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SIGN")), nil
	})

	return token, err
}

func ValidateToken(ctx *gin.Context) error {
	token, err := GetToken(ctx)

	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}

	return errors.New("Invalid authorization token")
}

func CurrentUser(ctx *gin.Context, db *gorm.DB) (models.User, error) {
	err := ValidateToken(ctx)
	if err != nil {
		return models.User{}, err
	}

	token, _ := GetToken(ctx)
	claims, _ := token.Claims.(jwt.MapClaims)

	userId := uint(claims["id"].(float64))

	user, err := models.GetUserById(userId, db)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
