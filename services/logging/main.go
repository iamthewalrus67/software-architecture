package main

import (
	"app/internal/common"
	"app/internal/logging"

	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hazelcast/hazelcast-go-client"
)

var userMessages map[string]string = make(map[string]string)

func serverHandler(w http.ResponseWriter, r *http.Request, client *hazelcast.Client) {
	if r.Method == http.MethodPost {
		logging.InfoLog.Println("Got POST request")

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		bodyString := string(body)

		bodyString = bodyString[1 : len(bodyString)-1]
		splitIndex := strings.Index(bodyString, ",")
		id, msg := bodyString[:splitIndex], bodyString[splitIndex+2:]
		logging.InfoLog.Println("Received new message. Saving...")
		logging.InfoLog.Printf("UUID: %s\nMessage: %s\n", id, msg)
		userMessages[id] = msg

		w.WriteHeader(http.StatusOK)

	} else if r.Method == http.MethodGet {
		logging.InfoLog.Println("Received GET request")

		values := make([]string, len(userMessages))
		i := 0
		for _, val := range userMessages {
			values[i] = val
			i++
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, strings.Join(values, " | "))
	} else {
		logging.WarningLog.Println("Received other request")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Incorrect request")
	}
}

func main() {
	ctx := context.TODO()
	config := hazelcast.NewConfig()
	config.Cluster.Network.Addresses = append(config.Cluster.Network.Addresses, "hazelcast")
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	common.PanicIfErr(err)

	logging.ErrorLog.Fatal(http.ListenAndServe(common.LoggingServicePort, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHandler(w, r, client)
	})))
}
