// Package worker a zmq backend worker queue
//
// By taka@cmwang.net
//
package worker

import "net/http"

var defaultSrv, _ = NewService()

// RequestHandler http request handler
func RequestHandler(ts int64, command, payload string, w http.ResponseWriter) {
	(*defaultSrv).requestHandler(ts, command, payload, w)
}

// Start start service
func Start() {
	(*defaultSrv).start()
}

// Stop stop service
func Stop() {
	(*defaultSrv).stop()
}
