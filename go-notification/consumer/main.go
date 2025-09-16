package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"github.com/rs/zerolog/log"
)

const (
	ConsumerServer    = ":8082"
	ConsumerTopic     = "notification"
	ConsumerGroupName = "go-notification-consumer-group"
	KafkaServer       = "localhost:9092"
)

// ============ SARAMA ===============
// create handler from sarama ConsumerGroupHandler interface
// ref: https://pkg.go.dev/github.com/IBM/sarama#ConsumerGroupHandler
type consumerGroupHandler struct {
	Store []*sarama.ConsumerMessage
}

func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d Key:%q Value:%q\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		c.Store = append(c.Store, msg)
		session.MarkMessage(msg, "")
	}

	return nil
}

func setupGroup() sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	// note: ini hanya berfungsi jika nama consumer-nya berbeda
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	cgrupname := ConsumerGroupName
	// for debug only
	uuid, _ := uuid.GenerateUUID()
	cgrupname = fmt.Sprintf("%s-%s", ConsumerGroupName, uuid)

	group, err := sarama.NewConsumerGroup([]string{KafkaServer}, cgrupname, config)
	if err != nil {
		panic(err)
	}

	return group
}

// ============ GIN ===============
type RequestUserUri struct {
	User string `uri:"user" binding:"required"`
}

func getMessageHandler(c *gin.Context, group sarama.ConsumerGroup, handler *consumerGroupHandler) {
	var user RequestUserUri
	if err := c.ShouldBindUri(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	err := group.Consume(context.Background(), []string{ConsumerTopic}, handler)
	if err != nil {
		log.Err(err).Msg("Consumer Group Error")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "OK", "user": handler.Store})
}

func main() {

	// setup kafka consumer
	group := setupGroup()
	defer func() {
		_ = group.Close()
	}()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// setup gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/api/message/:user", func(c *gin.Context) {
		handler := consumerGroupHandler{
			Store: make([]*sarama.ConsumerMessage, 0),
		}

		getMessageHandler(c, group, &handler)
	})
	router.Run(ConsumerServer)

}
