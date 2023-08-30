package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikariiin/dbvr-go/models"
	"github.com/ikariiin/dbvr-go/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterUserDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthRoutes struct {
	db     *gorm.DB
	router *gin.Engine
}

func NewAuthRoutes(db *gorm.DB, router *gin.Engine) *AuthRoutes {
	return &AuthRoutes{db: db, router: router}
}

func (r *AuthRoutes) registerUser(ctx *gin.Context) {
	var input RegisterUserDTO

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
	}

	user := models.User{Username: input.Username, Password: input.Password}
	user.HashPassword()

	if err := r.db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func (r *AuthRoutes) loginCheck(username, password string) (string, error) {
	var err error

	user := models.User{}

	if err = r.db.Model(models.User{}).Where("username=?", username).Take(&user).Error; err != nil {
		return "", err
	}

	err = models.VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := utils.GenerateToken(user)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *AuthRoutes) login(ctx *gin.Context) {
	var input LoginUserDTO

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user := models.User{Username: input.Username, Password: input.Password}

	token, err := r.loginCheck(user.Username, user.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username-password pair not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *AuthRoutes) RegisterAuthRoutes() {
	group := r.router.Group("auth")

	group.POST("register", r.registerUser)
	group.POST("login", r.login)
}
