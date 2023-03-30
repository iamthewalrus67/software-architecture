package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"app/internal/common"
	"app/internal/logging"
)

type FacadeService struct {
}

func NewFacadeService() *FacadeService {
	return &FacadeService{}
}

func (f *FacadeService) LogMessage(msg string) error {
	id := uuid.New()

	message := common.NewMessage(id, msg)

	_, err := http.Post(common.LoggingServiceAddress, "text", bytes.NewReader(message.ToJSON()))

	if err != nil {
		return err
	}

	return nil
}

func (f *FacadeService) GetAllMessages() (string, error) {
	res, err := getRequestToService(common.MessageServiceAddress)
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
