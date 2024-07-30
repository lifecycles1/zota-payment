package kafka

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	adminClient *kafka.AdminClient
	once        sync.Once
)

// getAdminClient creates or returns a singleton Kafka Admin client instance
func getAdminClient() *kafka.AdminClient {
	once.Do(func() {
		var err error
		adminClient, err = kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS")})
		if err != nil {
			log.Fatalf("Failed to create admin client: %v", err)
		}
	})
	return adminClient
}

// CreateKafkaTopic creates a Kafka topic with the specified number of partitions and replication factor
func CreateKafkaTopic(topic string, partitions, replicationFactor int) {
	admin := getAdminClient()
	existingTopics, err := admin.GetMetadata(nil, true, 5000)
	if err != nil {
		log.Fatalf("Failed to get Kafka metadata: %v", err)
	}

	if _, exists := existingTopics.Topics[topic]; exists {
		log.Printf("Topic %s already exists", topic)
		return
	}

	topicSpec := kafka.TopicSpecification{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}

	ctx := context.Background()

	_, err = admin.CreateTopics(ctx, []kafka.TopicSpecification{topicSpec}, nil)

	if err != nil {
		log.Fatalf("Failed to create topic %s: %v", topic, err)
	}

	log.Printf("Topic %s created with %d partitions and replication factor %d", topic, partitions, replicationFactor)
}
