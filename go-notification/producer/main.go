package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

/**

Sebagai API, menerima post dan menerima data
{user: string, message:string}
dan push ke kafka

**/

type responseMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type formMessage struct {
	FromUser string `form:"fromUser" json:"fromUser" binding:"required"`
	ToUser   string `form:"toUser" json:"toUser" binding:"required"`
	Message  string `form:"message" json:"message" binding:"required"`
}

func postMessage(c *gin.Context) {

	data := formMessage{}

	if err := c.ShouldBind(&data); err != nil {
		log.Info().Err(err).Msg("Form Request Failed")
		c.IndentedJSON(http.StatusBadRequest, responseMessage{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	jsonData, _ := json.Marshal(data)
	log.Debug().RawJSON("data", jsonData).Msg("form request")

	c.IndentedJSON(http.StatusOK, data)
}

func main() {

	router := gin.Default()
	router.POST("/api/message", postMessage)

	router.Run("localhost:8001")
}
