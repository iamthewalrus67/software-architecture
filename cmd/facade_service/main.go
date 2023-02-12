package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"app/internal/common"

	"github.com/google/uuid"
)

func serverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := uuid.New()
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		stringToSend := fmt.Sprintf("{%s, %s}", id.String(), body)
		resp, err := http.Post(common.LoggingServiceAddress, "text", strings.NewReader(stringToSend))

		if err != nil {
			fmt.Printf("error: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp.Body.Close()
			return

		}

		resp.Body.Close()
		w.WriteHeader(http.StatusOK)

	} else if r.Method == http.MethodGet {
		logginServiceResult, err := getRequestToService(common.LoggingServiceAddress)

		if err != nil {
			fmt.Printf("error: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		messageServiceResult, err := getRequestToService(common.MessageServiceAddress)

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

func main() {
	log.Fatal(http.ListenAndServe(common.FacadeServicePort, http.HandlerFunc(serverHandler)))
}
