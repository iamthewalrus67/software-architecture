package main

import (
	"app/internal/common"
	"app/internal/logging"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var userMessages map[string]string = make(map[string]string)

func serverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		bodyString := string(body)

		bodyString = bodyString[1 : len(bodyString)-1]
		splitIndex := strings.Index(bodyString, ",")
		id, msg := bodyString[:splitIndex], bodyString[splitIndex+2:]
		logging.Log("Received new message. Saving...")
		logging.Logf("UUID: %s\nMessage: %s", id, msg)
		userMessages[id] = msg

		w.WriteHeader(http.StatusOK)

	} else if r.Method == http.MethodGet {
		values := make([]string, len(userMessages))
		i := 0
		for _, val := range userMessages {
			values[i] = val
			i++
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, strings.Join(values, " | "))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Incorrect request")
	}
}

func main() {
	log.Fatal(http.ListenAndServe(common.LoggingServicePort, http.HandlerFunc(serverHandler)))
}
