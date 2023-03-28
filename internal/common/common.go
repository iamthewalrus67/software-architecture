package common

import (
	"fmt"
	"net/http"
	"os"
)

var FacadeServicePort string = ":8080"
var LoggingServicePort string = ":8081"
var MessageServicePort string = ":8082"

var FacadeServiceAddress string = "http://" + endpointOrDefault("FACADE_ENDPOINT") + FacadeServicePort
var LoggingServiceAddress string = "http://" + endpointOrDefault("LOGGING_ENDPOINT") + LoggingServicePort
var MessageServiceAddress string = "http://" + endpointOrDefault("MESSAGE_ENDPOINT") + MessageServicePort

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func endpointOrDefault(endpoint string) string {
	res := os.Getenv(endpoint)

	if res == "" {
		return "http://localhost"
	}

	return res
}

func DummyServerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "dummy post response")

	} else if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "dummy get response")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "dummy other request")
	}
}
