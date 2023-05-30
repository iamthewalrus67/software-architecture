package main

import (
	"app/services/logging/web"
)

func main() {
	loggingWeb := web.NewLoggingWeb()
	loggingWeb.Start()
}
