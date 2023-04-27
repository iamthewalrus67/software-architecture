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
		"bootstrap.servers": kafkaAddress})

	if err != nil {
		logging.ErrorLog.Println(err)
		common.PanicIfErr(err)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())

	return &KafkaConsumer{cons: cons, ctx: ctx, cancelCtx: cancelCtx}
}

func (k *KafkaConsumer) ReceiveMessages(repo repository.MessageRepository) {

	go func() {
		msg_count := 0
		run := true
		for run {
			select {
			case <-k.ctx.Done():
				run = false
			default:
				ev := k.cons.Poll(100)
				switch e := ev.(type) {
				case *kafka.Message:
					msg_count += 1
					if msg_count%5 == 0 {
						k.cons.Commit()
					}
					logging.InfoLog.Printf("Message on %s:\n%s\n",
						e.TopicPartition, string(e.Value))

				case kafka.PartitionEOF:
					logging.WarningLog.Printf("Reached %v\n", e)
				case kafka.Error:
					logging.ErrorLog.Printf("Error: %v\n", e)
					run = false
				default:
					logging.WarningLog.Printf("Ignored %v\n", e)
				}
			}
		}
	}()

}

func (k *KafkaConsumer) Stop() {
	k.cancelCtx()
	// k.ctx.Done()
}
