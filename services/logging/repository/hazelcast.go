package repository

import (
	"app/internal/common"
	"context"

	"github.com/hazelcast/hazelcast-go-client"
)

type HazelcastRepository struct {
	hzClient *hazelcast.Client
}

func NewHazelcastRepository() *HazelcastRepository {
	ctx := context.TODO()
	config := hazelcast.NewConfig()
	config.Cluster.Network.Addresses = append(config.Cluster.Network.Addresses, "hazelcast")
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	common.PanicIfErr(err)
	return &HazelcastRepository{hzClient: client}
}
