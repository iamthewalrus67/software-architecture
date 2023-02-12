package common

var FacadeServicePort string = ":8080"
var LoggingServicePort string = ":8081"
var MessageServicePort string = ":8082"

var FacadeServiceAddress string = "http://localhost" + FacadeServicePort
var LoggingServiceAddress string = "http://localhost" + LoggingServicePort
var MessageServiceAddress string = "http://localhost" + MessageServicePort
