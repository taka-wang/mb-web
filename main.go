package main

import (
	route "github.com/taka-wang/mb-web/route"
	worker "github.com/taka-wang/mb-web/worker"
)

func main() {
	worker.Start()
	// link worker and route packages
	route.SetRequestHandler(worker.RequestHandler)
	route.SetRoute()
	route.Start()
}
