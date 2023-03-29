package repository

import (
	"app/internal/common"
	"context"

	"github.com/hazelcast/hazelcast-go-client"
)

type LoggingRepository struct {
	hzClient *hazelcast.Client
}

func NewLoggingRepository() LoggingRepository {
	ctx := context.TODO()
	config := hazelcast.NewConfig()
	config.Cluster.Network.Addresses = append(config.Cluster.Network.Addresses, "hazelcast")
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	common.PanicIfErr(err)
	return LoggingRepository{hzClient: client}
}
