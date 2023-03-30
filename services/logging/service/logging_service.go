package service

import (
	"app/internal/common"
	"app/services/logging/repository"
)

type LoggingService struct {
	repo repository.LoggingRepository
}

func NewLoggingService() *LoggingService {
	return NewLoggingServiceWithRepository(repository.NewHazelcastRepository())
}

func NewLoggingServiceWithRepository(repo repository.LoggingRepository) *LoggingService {
	return &LoggingService{repo: repo}
}

func (l *LoggingService) AddMessage(msg common.Message) {
	l.repo.AddMessage(msg)
}

func (l *LoggingService) GetAllMessages() []common.Message {
	return l.repo.GetAllMessages()
}
