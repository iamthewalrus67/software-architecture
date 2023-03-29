package main

import (
	"app/services/facade/web"
)

func main() {
	facadeWeb := web.NewFacadeWeb()
	facadeWeb.Start()
}
