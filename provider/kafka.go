package provider

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

func CreateClientKafka(brokers []string) sarama.Client {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatalf("Failed create client kafka: %v", err)
	}
	return client
}

func SendMessage(producer sarama.SyncProducer, topicName string, messageText string) {
	message := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(messageText),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}

func CreateProducer(client sarama.Client) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	return producer
}

func CreateConsumer(client sarama.Client) sarama.Consumer {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	return consumer
}

func CreatePartitionConsumer(consumer sarama.Consumer, topicName string) sarama.PartitionConsumer {
	partitionConsumer, err := consumer.ConsumePartition(topicName, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to create partition consumer: %v", err)
	}
	return partitionConsumer
}
