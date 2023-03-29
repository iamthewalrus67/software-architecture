package service

import (
	"bytes"
	"io/ioutil"
	"net/http"

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

	message := common.NewMessage(id, msg)

	// stringToSend := fmt.Sprintf("{%s, %s}", id.String(), msg)
	// _, err := http.Post(common.LoggingServiceAddress, "text", strings.NewReader(stringToSend))
	_, err := http.Post(common.LoggingServiceAddress, "text", bytes.NewReader(message.ToJSON()))

	if err != nil {
		return err
	}

	return nil
}

// func formJSON(message string, uuid uuid.UUID) ([]byte, error) {
// 	data := map[string]interface{}{
// 		"message": message,
// 		"uuid":    uuid,
// 	}
//
// 	jsonData, err := json.Marshal(data)
//
// 	if err != nil {
// 		logging.ErrorLog.Println("Failed to marshal json")
// 		return make([]byte, 0), err
// 	}
//
// 	return jsonData, nil
// }

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
