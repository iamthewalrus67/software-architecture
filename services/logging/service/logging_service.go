package service

import (
	consul "github.com/hashicorp/consul/api"

	"app/internal/common"
	"app/internal/logging"
	"app/services/logging/repository"
)

type LoggingService struct {
	repo         repository.LoggingRepository
	consulClient *consul.Client
}

func NewLoggingService() *LoggingService {
	config := consul.DefaultConfig()
	config.Address = "consul:8500"
	consulClient, err := consul.NewClient(config)
	if err != nil {
		logging.ErrorLog.Fatal("failed to create consul client")
	}

	reg := &consul.AgentServiceRegistration{
		ID:      common.MyAddress,
		Name:    "logging",
		Port:    8081,
		Address: "http://" + common.MyAddress,
	}

	err = consulClient.Agent().ServiceRegister(reg)
	if err != nil {
		logging.ErrorLog.Fatal(err)
	}

	logging.InfoLog.Printf("Service %s registered with Consul\n", common.MyAddress)

	return &LoggingService{repo: repository.NewHazelcastRepository(consulClient), consulClient: consulClient}
}

func (l *LoggingService) AddMessage(msg common.Message) {
	l.repo.AddMessage(msg)
}

func (l *LoggingService) GetAllMessages() []common.Message {
	return l.repo.GetAllMessages()
}
