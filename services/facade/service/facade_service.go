package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	consul "github.com/hashicorp/consul/api"

	"app/internal/common"
	"app/internal/inet"
	"app/internal/logging"
	"app/services/facade/producer"
)

type FacadeService struct {
	prod         producer.Producer
	consulClient *consul.Client
}

func NewFacadeService() *FacadeService {
	config := consul.DefaultConfig()
	config.Address = "consul:8500"
	consulClient, err := consul.NewClient(config)
	if err != nil {
		logging.ErrorLog.Fatal("failed to create consul client")
	}

	reg := &consul.AgentServiceRegistration{
		ID:      common.MyAddress,
		Name:    "facade",
		Port:    8080,
		Address: "http://" + common.MyAddress,
	}

	err = consulClient.Agent().ServiceRegister(reg)
	if err != nil {
		logging.ErrorLog.Fatal(err)
	}

	logging.InfoLog.Printf("Service %s registered with Consul\n", common.MyAddress)

	return &FacadeService{prod: producer.NewKafkaProducer(consulClient), consulClient: consulClient}
}

func (f *FacadeService) SendMessage(msg common.Message) {
	f.prod.SendMessage(msg)
}

func (f *FacadeService) LogMessage(msg common.Message) error {
	address := inet.GetRandomLoggingIp(f.consulClient)
	_, err := http.Post(address, "text", bytes.NewReader(msg.ToJSON()))

	if err != nil {
		return err
	}

	return nil
}

func (f *FacadeService) GetAllMessages() (string, error) {
	address := inet.GetRandomMessageIp(f.consulClient)
	res, err := getRequestToService(address)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (f *FacadeService) GetAllLogs() ([]common.Message, error) {
	address := inet.GetRandomLoggingIp(f.consulClient)
	res, err := getRequestToService(address)

	if err != nil {
		logging.ErrorLog.Println("Failed to get logs")
		return make([]common.Message, 0), err
	}

	messages := make([]common.Message, 1)
	err = json.Unmarshal(res, &messages)

	if err != nil {
		logging.ErrorLog.Println("Failed to unmarshal")
		return make([]common.Message, 0), err
	}

	return messages, nil
}

func (f *FacadeService) GetAllLogsText() ([]string, error) {
	res, err := f.GetAllLogs()

	if err != nil {
		return make([]string, 0), err
	}

	values := make([]string, len(res))
	for i, msg := range res {
		values[i] = msg.Text
	}

	return values, nil
}

func getRequestToService(address string) ([]byte, error) {
	resp, err := http.Get(address)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return make([]byte, 0), err
	}

	result, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return make([]byte, 0), err
	}

	return result, nil
}
