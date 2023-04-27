package producer

import (
	"app/internal/common"
	"app/internal/logging"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	prod *kafka.Producer
}

func NewKafkaProducer() *KafkaProducer {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")

	if kafkaAddress == "" {
		kafkaAddress = "localhost:29092"
	}

	confMap := &kafka.ConfigMap{
		"bootstrap.servers": kafkaAddress}

	prod, err := kafka.NewProducer(confMap)

	logging.ErrorLog.Println(err)
	common.PanicIfErr(err)

	return &KafkaProducer{prod: prod}
}

func (k *KafkaProducer) SendMessage(msg common.Message) {
	deliveryChan := make(chan kafka.Event, 10000)
	topic := "messages"
	err := k.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg.ToJSON())},
		deliveryChan,
	)

	if err != nil {
		logging.ErrorLog.Println(err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		logging.ErrorLog.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		logging.InfoLog.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	close(deliveryChan)
}
