package common

import (
	"fmt"
	"net/http"
	"os"
)

var FacadeServicePort string = ":8080"
var LoggingServicePort string = ":8081"
var MessageServicePort string = ":8082"

var MyAddress string = os.Getenv("MY_ADDRESS")

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
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
