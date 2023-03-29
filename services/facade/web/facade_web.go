package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"app/internal/common"
	"app/services/facade/service"
)

type FacadeWeb struct {
	port string
	serv service.FacadeService
}

func NewFacadeWeb() FacadeWeb {
	return FacadeWeb{port: common.FacadeServicePort, serv: service.NewFacadeService()}
}

func (f *FacadeWeb) Start() {
	fmt.Println(f.port)
	log.Fatal(http.ListenAndServe(f.port, http.HandlerFunc(f.serverHandler)))
}

func (f *FacadeWeb) serverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("Received POST request")
		body, err := ioutil.ReadAll(r.Body)

		err = f.serv.LogMessage(string(body))

		if err != nil {
			fmt.Printf("error: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		w.WriteHeader(http.StatusOK)

	} else if r.Method == http.MethodGet {
		log.Println("Received GET request")
		logginServiceResult, err := f.serv.GetAllLogs()

		if err != nil {
			fmt.Printf("error: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		messageServiceResult, err := f.serv.GetAllMessages()

		if err != nil {
			fmt.Printf("error: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, messageServiceResult+": "+logginServiceResult)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Incorrect request")
	}
}
