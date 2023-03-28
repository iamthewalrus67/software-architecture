package main

import (
	"app/services/facade/web"
)

func main() {
	web := web.NewFacadeWeb()
	web.Start()
}
