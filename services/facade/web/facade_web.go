package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"app/internal/common"
	"app/internal/logging"
	"app/services/facade/service"
)

type FacadeWeb struct {
	port string
	serv *service.FacadeService
}

func NewFacadeWeb() FacadeWeb {
	return FacadeWeb{port: common.FacadeServicePort, serv: service.NewFacadeService()}
}

func (f *FacadeWeb) Start() {
	logging.InfoLog.Println("Started facade service")
	log.Fatal(http.ListenAndServe(f.port, http.HandlerFunc(f.serverHandler)))
}

func (f *FacadeWeb) serverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		logging.InfoLog.Println("Received POST request")

		body, err := ioutil.ReadAll(r.Body)

		msg := common.GenerateNewMessage(string(body))

		err = f.serv.LogMessage(msg)

		if err != nil {
			logging.ErrorLog.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		f.serv.SendMessage(msg)

		w.WriteHeader(http.StatusOK)

	} else if r.Method == http.MethodGet {
		logging.InfoLog.Println("Received GET request")

		logs, err := f.serv.GetAllLogsText()

		if err != nil {
			logging.ErrorLog.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		message, err := f.serv.GetAllMessages()

		if err != nil {
			logging.ErrorLog.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, message+": "+fmt.Sprintf("%v", logs))
	} else {
		logging.WarningLog.Println("Received other request")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Incorrect request")
	}
}
