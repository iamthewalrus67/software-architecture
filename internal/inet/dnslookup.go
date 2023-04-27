package inet

import (
	"app/internal/common"
	"app/internal/logging"
	"math/rand"
	"net"
)

var (
	MessageIps []string = nil
	LoggingIps []string = nil
)

func getIpsOfService(service string) []string {
	ips, err := net.LookupIP(service)
	if err != nil {
		logging.ErrorLog.Fatalf("Could not get IPs: %v\n", err)
	}

	ipsStr := make([]string, len(ips))

	for i, ip := range ips {
		ipsStr[i] = ip.String()
	}

	return ipsStr
}

func GetRandomMessageIp() string {
	if MessageIps == nil {
		MessageIps = getIpsOfService("message")
	}

	randIp := "http://" + GetRandomIp(MessageIps) + common.MessageServicePort
	return randIp
}

func GetRandomLoggingIp() string {
	if LoggingIps == nil {
		LoggingIps = getIpsOfService("logging")
	}

	randIp := "http://" + GetRandomIp(LoggingIps) + common.LoggingServicePort
	return randIp
}

func GetRandomIp(ipList []string) string {
	randIdx := rand.Intn(len(ipList))
	return ipList[randIdx]
}
