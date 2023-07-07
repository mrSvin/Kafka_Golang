package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"kafka/provider"
	"time"
)

const topicName = "test_topic"

func main() {

	client := provider.CreateClientKafka([]string{"localhost:9092"})
	defer client.Close()

	producer := provider.CreateProducer(client)
	defer producer.Close()

	provider.SendMessage(producer, topicName, "Hello")

	consumer := provider.CreateConsumer(client)
	defer consumer.Close()

	partitionConsumer := provider.CreatePartitionConsumer(consumer, topicName)
	defer partitionConsumer.Close()

	go listenerMessages(partitionConsumer)

	time.Sleep(5 * time.Second)

}

func listenerMessages(partitionConsumer sarama.PartitionConsumer) {
	for message := range partitionConsumer.Messages() {
		fmt.Printf("message: %s, offset: %v \n", string(message.Value), message.Offset)

	}
}
