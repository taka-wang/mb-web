package route

import "github.com/go-zoo/bone"

var defaultServer = NewService()

//
// Exported API
//

// GetMux get mux instance
func GetMux() *bone.Mux {
	return defaultServer.getMux()
}

// SetRequestHandler set request handler from external
func SetRequestHandler(f RequestHandler) {
	requestHandler = f
}

// Start start http server
func Start() {
	defaultServer.setRoute()
	defaultServer.start()
}

// SetRoute set route middleware
func SetRoute() {
	defaultServer.setRoute()
}
