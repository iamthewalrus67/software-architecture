package main

import (
	"app/internal/common"
	"app/internal/logging"
	"fmt"
	"net/http"
)

func serverHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	logging.ErrorLog.Fatal(http.ListenAndServe(common.MessageServicePort, http.HandlerFunc(serverHandler)))
}
