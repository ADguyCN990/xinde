package account

import (
	"net/http"
	"time"
	"xinde/internal/model"
	"xinde/internal/store"
	"xinde/pkg/util"

	"github.com/gin-gonic/gin"
)

// Register handles user registration.
func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 重复输入密码需一致

	// Hash the password before saving
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now().Unix()

	// Save the user to the database
	if err := store.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
