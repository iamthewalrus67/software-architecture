package repository

import (
	"app/internal/common"
	"sync"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	userMessages map[uuid.UUID]common.Message
	mu           sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{userMessages: make(map[uuid.UUID]common.Message), mu: sync.RWMutex{}}
}

func (m *MemoryRepository) AddMessage(msg common.Message) {
	m.mu.Lock()
	m.userMessages[msg.UUID] = msg
	m.mu.Unlock()
}

func (m *MemoryRepository) GetAllMessages() []common.Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	values := make([]common.Message, len(m.userMessages))
	i := 0
	for _, val := range m.userMessages {
		values[i] = val
		i++
	}

	return values
}
