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
		fmt.Println(stringToSend)
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
		resp, err := http.Get(common.LoggingServiceAddress)

		if err != nil {
			fmt.Printf("error: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp.Body.Close()
			return
		}

		result, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Printf("error: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp.Body.Close()
			return
		}

		fmt.Fprint(w, string(result))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	log.Fatal(http.ListenAndServe(common.FacadeServicePort, http.HandlerFunc(serverHandler)))
}
