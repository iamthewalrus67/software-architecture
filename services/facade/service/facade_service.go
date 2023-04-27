package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"app/internal/common"
	"app/internal/inet"
	"app/internal/logging"
	"app/services/facade/producer"
)

type FacadeService struct {
	prod producer.Producer
}

func NewFacadeService() *FacadeService {
	return &FacadeService{prod: producer.NewKafkaProducer()}
}

func (f *FacadeService) SendMessage(msg common.Message) {
	f.prod.SendMessage(msg)
}

func (f *FacadeService) LogMessage(msg common.Message) error {
	_, err := http.Post(inet.GetRandomLoggingIp(), "text", bytes.NewReader(msg.ToJSON()))

	if err != nil {
		return err
	}

	return nil
}

func (f *FacadeService) GetAllMessages() (string, error) {
	res, err := getRequestToService(inet.GetRandomMessageIp())
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (f *FacadeService) GetAllLogs() ([]common.Message, error) {
	res, err := getRequestToService(common.LoggingServiceAddress)

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
