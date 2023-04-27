package service

import (
	"app/internal/common"
	"app/services/message/consumer"
	"app/services/message/repository"
)

type MessageService struct {
	repo            repository.MessageRepository
	messageConsumer consumer.Consumer
}

func NewMessageService() *MessageService {
	return &MessageService{repo: repository.NewMemoryRepository(), messageConsumer: consumer.NewKafkaConsumer()}
}

func NewMessageServiceWithRepository(repo repository.MessageRepository) *MessageService {
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
