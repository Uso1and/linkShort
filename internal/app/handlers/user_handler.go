package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"lnkshrt/internal/domain/config"
	"lnkshrt/internal/domain/models"
	"lnkshrt/internal/domain/repo"
)

type UserHandler struct {
	userRepo repo.UserRepoInterface
	cfg      *config.ConfigDB
}

func NewUserHandler(userRepo repo.UserRepoInterface, cfg *config.ConfigDB) *UserHandler {
	return &UserHandler{userRepo: userRepo, cfg: cfg}
}

func (h *UserHandler) CreateUserHandler(c *gin.Context) {

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid required"})
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, email and password required"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
	}

	if err := h.userRepo.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not create user"})
		return
	}
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "user create succes",
		"user":    user,
	})

}

func (h *UserHandler) LoginHandler(c *gin.Context) {
	var cread struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&cread); err != nil {
		log.Printf("invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	user, err := h.userRepo.GetUserByUsername(c.Request.Context(), cread.Username)

	if err != nil {
		log.Printf("error getting user: %v", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid cread"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cread.Password))
	if err != nil {
		log.Printf("error password mismatch:%v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid cread"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))

	if err != nil {
		log.Printf("error generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	log.Printf("Generated token for user %s: %s", user.Username, tokenString)
	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"expires": expirationTime,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
