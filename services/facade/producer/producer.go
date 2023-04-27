package producer

import "app/internal/common"

type Producer interface {
	SendMessage(msg common.Message)
}
