package repository

import (
	"app/internal/common"
	"app/internal/logging"
	"context"
	"encoding/gob"

	"github.com/hazelcast/hazelcast-go-client"
)

type HazelcastRepository struct {
	client *hazelcast.Client
	mp     *hazelcast.Map
}

func NewHazelcastRepository() *HazelcastRepository {
	gob.Register(common.Message{}) // For some reason Hazelcast client breaks without this
	ctx := context.TODO()
	config := hazelcast.NewConfig()
	config.Cluster.Network.Addresses = append(config.Cluster.Network.Addresses, "hazelcast")
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	common.PanicIfErr(err)
	mp, err := client.GetMap(ctx, "userMessages")
	return &HazelcastRepository{client: client, mp: mp}
}

func (h *HazelcastRepository) AddMessage(msg common.Message) {
	ctx := context.TODO()
	h.mp.Lock(ctx, msg.UUID)
	h.mp.Set(ctx, msg.UUID, msg)
	h.mp.Unlock(ctx, msg.UUID)
}

func (h *HazelcastRepository) GetAllMessages() []common.Message {
	ctx := context.TODO()
	values, err := h.mp.GetValues(ctx)

	if err != nil {
		logging.ErrorLog.Println(err)
		return make([]common.Message, 0)
	}

	messages := make([]common.Message, len(values))

	for i, val := range values {
		messages[i] = val.(common.Message)
	}

	if err != nil {
		logging.ErrorLog.Println(err)
		return make([]common.Message, 0)
	}

	return messages
}
