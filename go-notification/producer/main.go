package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

/**

Sebagai API, menerima post dan menerima data
{user: string, message:string}
dan push ke kafka

**/

const (
	KafkaServer   = "localhost:9092"
	ProducerPort  = ":8001"
	ProducerTopic = "notification"
)

type responseMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type formMessage struct {
	FromUser string `form:"fromUser" json:"fromUser" binding:"required"`
	ToUser   string `form:"toUser" json:"toUser" binding:"required"`
	Message  string `form:"message" json:"message" binding:"required"`
}

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{KafkaServer}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup kafka producer : %w", err)
	}

	return producer, nil

}

func sendKafkaMessage(producer sarama.SyncProducer, message *formMessage) error {

	messageJson, _ := json.Marshal(message)

	msg := sarama.ProducerMessage{
		Topic: ProducerTopic,
		Key:   sarama.StringEncoder(message.ToUser),
		Value: sarama.StringEncoder(messageJson),
	}

	_, _, err := producer.SendMessage(&msg)

	return err
}

func messageHandler(producer sarama.SyncProducer) gin.HandlerFunc {

	return func(c *gin.Context) {

		// get form input
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

		// send kafka message
		err := sendKafkaMessage(producer, &data)
		if err != nil {
			log.Info().Err(err).Msg("Send Kafka Message Failed")
			c.IndentedJSON(http.StatusBadRequest, responseMessage{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}

		c.IndentedJSON(http.StatusOK, responseMessage{
			Status:  http.StatusOK,
			Message: "success",
		})
	}
}

func main() {

	// setup kafka producer
	producer, err := setupProducer()
	if err != nil {
		log.Err(err).Msg("Kafka Producer Connection Error")
		return
	}

	defer producer.Close()

	// setup gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/api/message", messageHandler(producer))

	router.Run(ProducerPort)
}
