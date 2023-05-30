package producer

import (
	"app/internal/common"
	"app/internal/logging"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	consul "github.com/hashicorp/consul/api"
)

type KafkaProducer struct {
	prod  *kafka.Producer
	topic string
}

func NewKafkaProducer(consul *consul.Client) *KafkaProducer {
	kafkaAddress, _, err := consul.KV().Get("kafka/server", nil)
	if err != nil {
		logging.ErrorLog.Fatal(err)
	}

	confMap := &kafka.ConfigMap{
		"bootstrap.servers": string(kafkaAddress.Value)}

	prod, err := kafka.NewProducer(confMap)

	if err != err {
		logging.ErrorLog.Println(err)
		common.PanicIfErr(err)
	}

	topic, _, err := consul.KV().Get("kafka/topic_name", nil)
	if err != nil {
		logging.ErrorLog.Fatal(err)
	}

	return &KafkaProducer{prod: prod, topic: string(topic.Value)}
}

func (k *KafkaProducer) SendMessage(msg common.Message) {
	deliveryChan := make(chan kafka.Event, 10000)
	err := k.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.topic, Partition: kafka.PartitionAny},
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
