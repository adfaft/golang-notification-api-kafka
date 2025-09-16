package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
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
	// todo: Store ini dibutuhkan semakin complex
	// karena menggunakan goroutine, setelah dipanggil oleh API harus dihapus dari memory
	Store map[string][]*sarama.ConsumerMessage
}

func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Debug().
			Str("data", fmt.Sprintf("Message topic:%q partition:%d offset:%d Key:%q Value:%q\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)).
			Msg("Message")

		c.Store[string(msg.Key)] = append(c.Store[string(msg.Key)], msg)
		session.MarkMessage(msg, "")
	}

	return nil
}

func setupGroupConsumer(ctx context.Context, user string, handler *consumerGroupHandler) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	// note: ini hanya berfungsi jika nama consumer-nya berbeda
	// config.Consumer.Offsets.Initial = sarama.OffsetOldest

	cgrupname := ConsumerGroupName
	// for debug only
	// uuid, _ := uuid.GenerateUUID()
	// cgrupname = fmt.Sprintf("%s-%s", ConsumerGroupName, uuid)

	group, err := sarama.NewConsumerGroup([]string{KafkaServer}, cgrupname, config)
	if err != nil {
		log.Err(err).Msg("Consumer Group Error GrupInit")
		// panic(err)
	}
	defer group.Close()

	topicName := fmt.Sprintf("%s-%s", ConsumerTopic, user)
	fmt.Printf("topic: %s, consumer ID: %s\n", topicName, cgrupname)

	for {
		errConsume := group.Consume(ctx, []string{topicName}, handler)
		if errConsume != nil {
			log.Err(errConsume).Msg("Consumer Group Error Consume")
		}

		if ctx.Err() != nil {
			return
		}
	}

}

// ============ GIN ===============
type RequestUserUri struct {
	User string `uri:"user" binding:"required"`
}

func getMessageHandler(c *gin.Context, ctx context.Context, handler *consumerGroupHandler) {
	var user RequestUserUri
	if err := c.ShouldBindUri(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// TODO: jika cancellable context disini, maka langsung berakhir tanpa consume data
	// bagaimana caranya agar goroutine ini berjalan hanya beberapa detik saat dibutuhkan saja? beberapa "detik" apakah konsep ini sesuai?
	// setup kafka consumer and run in separate thread in background
	go setupGroupConsumer(ctx, user.User, handler)

	if ctx.Err() != nil {
		log.Err(ctx.Err()).Msg("Consumer Group Error Context")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": ctx.Err().Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "OK", "user": handler.Store})
}

func main() {

	// cancel ini dari main program
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := consumerGroupHandler{
		Store: make(map[string][]*sarama.ConsumerMessage, 0),
	}

	// setup gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/api/message/:user", func(c *gin.Context) {
		getMessageHandler(c, ctx, &handler)
	})
	fmt.Printf("begin consumer on %s", ConsumerServer)
	router.Run(ConsumerServer)

}
