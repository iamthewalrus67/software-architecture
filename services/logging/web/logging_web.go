package web

import (
	"app/internal/common"
	"app/services/logging/repository"
)

type LoggingWeb struct {
	port string
	repo repository.LoggingRepository
}

func NewLoggingWeb() LoggingWeb {
	return LoggingWeb{port: common.LoggingServicePort, repo: repository.NewLoggingRepository()}
}
