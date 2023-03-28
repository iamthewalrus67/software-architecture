package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"app/internal/common"
)

type FacadeService struct {
}

func NewFacadeService() FacadeService {
	return FacadeService{}
}

func (f *FacadeService) LogMessage(msg string) error {
	id := uuid.New()

	stringToSend := fmt.Sprintf("{%s, %s}", id.String(), msg)
	_, err := http.Post(common.LoggingServiceAddress, "text", strings.NewReader(stringToSend))

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

	return res, nil
}

func (f *FacadeService) GetAllLogs() (string, error) {
	res, err := getRequestToService(common.LoggingServiceAddress)

	if err != nil {
		return "", err
	}

	return res, nil
}

func getRequestToService(address string) (string, error) {
	resp, err := http.Get(address)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(result), nil

}
