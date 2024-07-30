package kafka

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"zota_payment/repositories"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// StartKafkaConsumer initializes and starts the Kafka consumer
func StartKafkaConsumer(topic string, callbackRepo *repositories.CallbackRepository) {
	// Create a new Kafka consumer with the provided bootstrap servers and group ID
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          "mygroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer c.Close()

	// Subscribe to the specified topic
	c.SubscribeTopics([]string{topic}, nil)

	// Handle graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	run := true

	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			// Read messages from Kafka
			msg, err := c.ReadMessage(-1)
			if err == nil {
				log.Printf("Consumer Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				err := callbackRepo.InsertMessage(string(msg.Value))
				if err != nil {
					log.Printf("Failed to insert message: %v", err)
				}
			} else {
				log.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}

	log.Printf("Consumer shutting down")
}
