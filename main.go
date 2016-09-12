package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/taka-wang/mb-web/route"
	"github.com/taka-wang/mb-web/worker"
)

func main() {

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// start cleanup
		worker.Stop()
		// stop cleanup
		os.Exit(0)
	}()

	worker.Start()
	// link worker and route packages
	route.SetRequestHandler(worker.RequestHandler)
	route.SetRoute()
	route.Start()
}
