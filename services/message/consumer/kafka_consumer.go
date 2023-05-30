package consumer

import (
	"app/internal/common"
	"app/internal/logging"
	"app/services/message/repository"
	"context"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	cons      *kafka.Consumer
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func NewKafkaConsumer() *KafkaConsumer {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")

	if kafkaAddress == "" {
		kafkaAddress = "localhost:29092"
	}

	cons, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaAddress,
		"group.id":          "group"})

	if err != nil {
		logging.ErrorLog.Println(err)
		common.PanicIfErr(err)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())

	return &KafkaConsumer{cons: cons, ctx: ctx, cancelCtx: cancelCtx}
}

func (k *KafkaConsumer) ReceiveMessages(repo repository.MessageRepository) {
	err := k.cons.Subscribe("messages", nil)
	if err != nil {
		logging.ErrorLog.Println(err)
	}

	go func() {
		run := true
		for run {
			select {
			case <-k.ctx.Done():
				run = false
			default:
				ev := k.cons.Poll(100)
				switch e := ev.(type) {
				case *kafka.Message:
					msg, err := common.MessageFromBytes(e.Value)
					if err != nil {
						logging.ErrorLog.Println(err)
					}
					logging.InfoLog.Printf("Got message: %s\n", msg)

					repo.AddMessage(msg)

				case kafka.PartitionEOF:
					logging.ErrorLog.Printf("Reached %v\n", e)
				case kafka.Error:
					logging.ErrorLog.Printf("Error: %v\n", e)
					run = false
				default:

				}
			}
		}
	}()

}

func (k *KafkaConsumer) Stop() {
	k.cancelCtx()
	k.cons.Close()
}
