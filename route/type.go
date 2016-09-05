package route

import (
	"net/http"

	"github.com/go-zoo/bone"
)

type (
	// RequestHandler request handler type
	RequestHandler func(int64, string, string, http.ResponseWriter)

	// Service web server
	Service struct {
		// mux mux instance
		mux *bone.Mux
		// ip web ip
		ip string
		// port web port
		port string
	}
)

// ==================== Request Types ====================

type (

/*
	// ReqMbtcpTimeout set timeout request
	ReqMbtcpTimeout struct {
		Data int64 `json:"timeout,omitempty"`
	}

	// MbtcpWriteReq write coil/register request
	MbtcpWriteReq struct {
		IP    string      `json:"ip"`
		Port  string      `json:"port,omitempty"`
		Slave uint8       `json:"slave"`
		Addr  uint16      `json:"addr"`
		Len   uint16      `json:"len,omitempty"`
		Hex   bool        `json:"hex,omitempty"`
		Data  interface{} `json:"data"`
	}
*/
)

// ==================== Response Types ====================
