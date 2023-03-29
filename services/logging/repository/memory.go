package repository

import (
	"app/internal/common"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	userMessages map[uuid.UUID]common.Message
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{userMessages: make(map[uuid.UUID]common.Message)}
}

func (m *MemoryRepository) AddMessage(msg common.Message) {
	m.userMessages[msg.UUID] = msg
}

func (m *MemoryRepository) GetAllMessages() []common.Message {
	values := make([]common.Message, len(m.userMessages))
	i := 0
	for _, val := range m.userMessages {
		values[i] = val
		i++
	}

	return values
}
