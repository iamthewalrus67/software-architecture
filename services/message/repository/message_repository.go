package repository

import (
	"app/internal/common"
)

type MessageRepository interface {
	AddMessage(msg common.Message)
	GetAllMessages() []common.Message
}
