package web

import (
	"app/internal/common"
	"app/internal/logging"
	"fmt"
	"net/http"
)

type MessageWeb struct {
	port string
}

func NewMessageWeb() *MessageWeb {
	return &MessageWeb{port: common.MessageServicePort}
}

func (m *MessageWeb) serverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		logging.InfoLog.Println("Received GET request")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Not implemented")

	} else {
		logging.InfoLog.Println("Received other request")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Incorrect request")
	}

}

func (m *MessageWeb) Start() {
	logging.ErrorLog.Fatal(http.ListenAndServe(common.MessageServicePort, http.HandlerFunc(m.serverHandler)))
}
