package main

import (
	"app/services/message/web"
)

func main() {
	messageWeb := web.NewMessageWeb()
	messageWeb.Start()
}
