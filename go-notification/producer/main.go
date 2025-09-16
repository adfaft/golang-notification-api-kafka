package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/**

Sebagai API, menerima post dan menerima data
{user: string, message:string}
dan push ke kafka

**/

func postMessage(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, gin.H{"test": "tersting"})
}

func main() {

	router := gin.Default()
	router.POST("/api/message", postMessage)

	router.Run("localhost:8001")
}
