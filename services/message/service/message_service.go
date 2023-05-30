package service

import (
	"app/internal/common"
	"app/internal/logging"
	"app/services/message/consumer"
	"app/services/message/repository"

	consul "github.com/hashicorp/consul/api"
)

type MessageService struct {
	repo            repository.MessageRepository
	messageConsumer consumer.Consumer
}

func NewMessageService() *MessageService {
	return NewMessageServiceWithRepository(repository.NewMemoryRepository())
}

func NewMessageServiceWithRepository(repo repository.MessageRepository) *MessageService {
	config := consul.DefaultConfig()
	config.Address = "consul:8500"
	consulClient, err := consul.NewClient(config)
	if err != nil {
		logging.ErrorLog.Fatal("failed to create consul client")
	}

	reg := &consul.AgentServiceRegistration{
		ID:      common.MyAddress,
		Name:    "message",
		Port:    8082,
		Address: "http://" + common.MyAddress,
	}

	err = consulClient.Agent().ServiceRegister(reg)
	if err != nil {
		logging.ErrorLog.Fatal(err)
	}

	logging.InfoLog.Printf("Service %s registered with Consul\n", common.MyAddress)

	return &MessageService{repo: repo, messageConsumer: consumer.NewKafkaConsumer()}
}

func (m *MessageService) StopConsumer() {
	m.messageConsumer.Stop()
}

func (m *MessageService) StartConsumer() {
	m.messageConsumer.ReceiveMessages(m.repo)
}

func (m *MessageService) GetMessages() []common.Message {
	return m.repo.GetAllMessages()
}

func (m *MessageService) GetMessagesText() []string {
	msgs := m.GetMessages()

	texts := make([]string, len(msgs))
	for i, msg := range msgs {
		texts[i] = msg.Text
	}

	return texts
}
