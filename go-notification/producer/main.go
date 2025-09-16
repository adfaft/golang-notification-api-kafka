package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**

Sebagai API, menerima post dan menerima data
{user: string, message:string}
dan push ke kafka

**/

func postMessage(c *gin.Context) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(`{"test": "tersting"}`), &data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "input is wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"test": "testing"})
}

func main() {

	type responseMessage struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	json, _ := json.Marshal(&responseMessage{
		Status:  200,
		Message: "success",
	})
	fmt.Println(string(json))

	// router := gin.Default()
	// router.POST("/api/message", postMessage)

	// router.Run("localhost:8001")
}
