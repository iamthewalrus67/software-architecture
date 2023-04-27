package consumer

import "app/services/message/repository"

type Consumer interface {
	ReceiveMessages(repo repository.MessageRepository)
	Stop()
}
