package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/graph/model"
	"github.com/speedrun-website/leaderboard-backend/middleware"
	"gorm.io/gorm"
)

func MeHandler(c *gin.Context) {
	user, _ := c.Get(middleware.JwtConfig.IdentityKey)
	db, err := database.GetDatabase()

	// todo error handler or middleware?
	if err != nil {
		log.Println("Unable to connect to database", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	var me model.User
	result := db.Where(model.User{
		Email: user.(*model.User).Email,
	}).First(&me)

	if result.Error != nil {
		var code = http.StatusInternalServerError

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		}

		c.AbortWithStatusJSON(code, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, me)
}
