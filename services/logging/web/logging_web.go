package web

import (
	"app/internal/common"
	"app/internal/logging"
	"app/services/logging/service"
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"
)

type LoggingWeb struct {
	port string
	serv *service.LoggingService
}

func NewLoggingWeb() LoggingWeb {
	return LoggingWeb{port: common.LoggingServicePort, serv: service.NewLoggingService()}
}

func (l *LoggingWeb) serverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		logging.InfoLog.Println("Got POST request")

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		msg, err := common.MessageFromBytes(body)

		if err != nil {
			logging.ErrorLog.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logging.InfoLog.Println("Received new message. Saving...")
		l.serv.AddMessage(msg)
		logging.InfoLog.Printf("UUID: %s\nMessage: %s\n", msg.UUID, msg.Text)

		w.WriteHeader(http.StatusOK)

	} else if r.Method == http.MethodGet {
		logging.InfoLog.Println("Received GET request")

		data, err := json.Marshal(l.serv.GetAllMessages())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logging.ErrorLog.Println("Failed to marshal json")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(data))
	} else {
		logging.WarningLog.Println("Received other request")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Incorrect request")
	}

}

func (l *LoggingWeb) Start() {
	logging.InfoLog.Println("Starting logging service")
	logging.ErrorLog.Fatal(http.ListenAndServe(common.LoggingServicePort, http.HandlerFunc(l.serverHandler)))
}
