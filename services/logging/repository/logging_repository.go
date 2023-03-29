package repository

import (
	"app/internal/common"
)

type LoggingRepository interface {
	AddMessage(msg common.Message)
	GetAllMessages() []common.Message
}
