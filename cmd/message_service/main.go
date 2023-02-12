package main

import (
	"app/internal/common"
	"fmt"
	"log"
	"net/http"
)

func serverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Not implemented")

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Incorrect request")
	}

}

func main() {
	log.Fatal(http.ListenAndServe(common.MessageServicePort, http.HandlerFunc(serverHandler)))
}
