package inet

import (
	"app/internal/logging"
	"fmt"
	"math/rand"

	consul "github.com/hashicorp/consul/api"
)

func getServiceEndpoint(consulClient *consul.Client, service string) string {
	instances, _, err := consulClient.Health().Service(service, "", true, nil)
	if err != nil {
		logging.ErrorLog.Fatal(err)
	}

	if len(instances) == 0 {
		logging.ErrorLog.Fatalf("No instances of service %s found \n", service)
	}
	randomInstance := instances[rand.Intn(len(instances))]

	return randomInstance.Service.Address + ":" + fmt.Sprint(randomInstance.Service.Port)
}

func GetRandomMessageIp(consulClient *consul.Client) string {
	return getServiceEndpoint(consulClient, "message")
}

func GetRandomLoggingIp(consulClient *consul.Client) string {
	return getServiceEndpoint(consulClient, "logging")
}

func GetRandomFacadeIp(consulClient *consul.Client) string {
	return getServiceEndpoint(consulClient, "facade")
}

func GetRandomIp(ipList []string) string {
	randIdx := rand.Intn(len(ipList))
	return ipList[randIdx]
}
