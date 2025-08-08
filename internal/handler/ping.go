package handler

import "github.com/gin-gonic/gin"

// Ping godoc
// @Summary      Ping server
// @Description  do ping
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
