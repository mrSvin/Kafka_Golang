package tests

import (
	"fmt"
	"github.com/Shopify/sarama"
	"kafka/provider"
	"log"
	"testing"
	"time"
)

const topicName2 = "test_topic2"

func Test_Kafka_Two_Consumer(t *testing.T) {

	client := provider.CreateClientKafka([]string{"localhost:9092", "localhost:9093", "localhost:9094"})
	defer client.Close()

	producer := provider.CreateProducer(client)
	defer producer.Close()

	consumer1 := provider.CreateConsumer(client)
	defer consumer1.Close()

	consumer2 := provider.CreateConsumer(client)
	defer consumer2.Close()

	offset, err := client.GetOffset(topicName, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed to get offset: ", err)
	}

	offset2, err := client.GetOffset(topicName2, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed to get offset: ", err)
	}

	consumer1PartitionConsumer, err := consumer1.ConsumePartition(topicName, 0, offset)
	if err != nil {
		log.Fatal("Failed to create consumer 1 partition consumer: ", err)
	}
	consumer2PartitionConsumer, err := consumer2.ConsumePartition(topicName2, 0, offset2)
	if err != nil {
		log.Fatal("Failed to create consumer 2 partition consumer: ", err)
	}

	// Создаем каналы для обработки сообщений
	consumer1Messages := consumer1PartitionConsumer.Messages()
	consumer2Messages := consumer2PartitionConsumer.Messages()

	// Создаем канал для ошибок
	errors := consumer1PartitionConsumer.Errors()

	go listenerConsumer(consumer1Messages, 1)
	go listenerConsumer(consumer2Messages, 2)
	go listenerErrors(errors)

	provider.SendMessage(producer, topicName, "Test 1")
	provider.SendMessage(producer, topicName2, "Test 2")

	time.Sleep(5 * time.Second)

}

func listenerConsumer(consumerMessages <-chan *sarama.ConsumerMessage, numberConsumer uint) {
	for message := range consumerMessages {
		fmt.Println("Consumer", numberConsumer, string(message.Value))
	}
}

func listenerErrors(consumerErrors <-chan *sarama.ConsumerError) {
	for errors := range consumerErrors {
		fmt.Println("Consumer error: ", errors.Error())
	}
}
