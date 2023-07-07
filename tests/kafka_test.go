package tests

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/Shopify/sarama"
	"kafka/provider"
	"log"
	"testing"
	"time"
)

const topicName = "test_topic"

var checkMessage = false

func Test_Kafka(t *testing.T) {

	client := provider.CreateClientKafka([]string{"localhost:9092"})
	defer client.Close()

	producer := provider.CreateProducer(client)
	defer producer.Close()

	testMessage := generateRandomHash()

	provider.SendMessage(producer, topicName, testMessage)

	consumer := provider.CreateConsumer(client)
	defer consumer.Close()

	partitionConsumer := provider.CreatePartitionConsumer(consumer, topicName)
	defer partitionConsumer.Close()

	go listenerMessages(partitionConsumer, testMessage)

	time.Sleep(5 * time.Second)

	if checkMessage == false {
		log.Fatal("Test fail, message not found")
	}

}

func generateRandomHash() string {
	hash := make([]byte, 16)
	_, err := rand.Read(hash)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hash)
}

func listenerMessages(partitionConsumer sarama.PartitionConsumer, searchMessage string) {
	for message := range partitionConsumer.Messages() {
		fmt.Printf("message: %s, offset: %v \n", string(message.Value), message.Offset)
		if searchMessage == string(message.Value) {
			checkMessage = true
		}

	}
}
