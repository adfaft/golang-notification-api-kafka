package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ConsumerServer = ":8082"
)

type RequestUserUri struct {
	User string `uri:"user" binding:"required"`
}

func getMessageHandler(c *gin.Context) {
	var user RequestUserUri
	if err := c.ShouldBindUri(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"status": "OK", "user": user.User})
}

func main() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/api/message/:user", getMessageHandler)
	router.Run(ConsumerServer)

}
