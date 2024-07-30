package kafka

import (
	"context"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *kafka.Producer
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewProducer initializes a new Kafka producer
func NewProducer() (*Producer, error) {
	// Create a new Kafka producer with the provided bootstrap servers
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS")})
	if err != nil {
		return nil, err
	}

	// Create a context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	producer := &Producer{
		producer: p,
		ctx:      ctx,
		cancel:   cancel,
	}

	// Start a goroutine to handle delivery reports
	go producer.handleEvents()

	return producer, nil
}

// handleEvents processes delivery reports from the producer
func (p *Producer) handleEvents() {
	for {
		select {
		case ev := <-p.producer.Events():
			switch e := ev.(type) {
			case *kafka.Message:
				if e.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", e.TopicPartition.Error)
				} else {
					log.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset)
				}
			}
		case <-p.ctx.Done():
			return
		}
	}
}

// SendMessage sends a message to a specified Kafka topic
func (p *Producer) SendMessage(topic string, message string) error {
	deliveryChan := make(chan kafka.Event)

	// Produce the message to the topic
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		log.Printf("Failed to produce message: %v\n", err)
		return err
	}

	// Handle delivery report asynchronously
	go func() {
		e := <-deliveryChan
		m := e.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		} else {
			log.Printf("Delivered message to %v\n", m.TopicPartition)
		}

		close(deliveryChan)
	}()

	return nil
}

// Close gracefully shuts down the producer
func (p *Producer) Close() {
	p.producer.Close()
	p.cancel()
}
