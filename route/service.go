// Package route a bone-based http server
//
// By taka@cmwang.net
//
package route

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/taka-wang/mb-web/conf"
)

var (
	// defaultMbPort default modbus slave port number
	defaultMbPort string
	// minConnTimeout minimal modbus tcp connection timeout
	minConnTimeout int64
	// minPollInterval minimal modbus tcp poll interval
	minPollInterval uint64
	// requestHandler HTTP request handler callback
	requestHandler RequestHandler
)

// setDefaults set default values
func setDefaults() {
	// set default psmbtcp values
	conf.SetDefault(keyTCPDefaultPort, defaultTCPDefaultPort)
	conf.SetDefault(keyMinConnectionTimout, defaultMinConnectionTimout)
	conf.SetDefault(keyPollInterval, defaultPollInterval)

	// set default web values
	conf.SetDefault(keyWebPrefix, defaultWebPrefix)
	conf.SetDefault(keyWebIP, defaultWebIP)
	conf.SetDefault(keyWebPort, defaultWebPort)
}

func init() {
	setDefaults()

	// set null handler
	requestHandler = func(int64, string, string, http.ResponseWriter) {}

	// get mbtcp default
	defaultMbPort = conf.GetString(keyTCPDefaultPort)
	minConnTimeout = conf.GetInt64(keyMinConnectionTimout)
	minPollInterval = uint64(conf.GetInt(keyPollInterval))
}

// NewService create http server
func NewService() Service {
	return Service{
		mux:  bone.New(),
		ip:   conf.GetString(keyWebIP),
		port: conf.GetString(keyWebPort),
	}
}

// getMux get mux instance
func (s *Service) getMux() *bone.Mux {
	return s.mux
}

// start start http server
func (s *Service) start() {
	http.ListenAndServe(s.ip+":"+s.port, s.mux)
}

// SetRoute set route middleware
func (s *Service) setRoute() {

	// Set prefix
	s.mux.Prefix(conf.GetString(keyWebPrefix))
	// ==================== general handlers ====================
	// 404 Not Found
	s.mux.NotFoundFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("404 NOT FOUND!"))
	})

	// ROOT
	s.mux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("version: " + version))
	}))

	// ==================== one-off command ====================
	s.mux.Get(pmbOnceRead, http.HandlerFunc(handleMbOnceRead))
	s.mux.Post(pmbOnceWrite, http.HandlerFunc(handleMbOnceWrite))
	s.mux.Get(pmbGetTimeout, http.HandlerFunc(handleMbGetTimeout))
	s.mux.Post(pmbSetTimeout, http.HandlerFunc(handleMbSetTimeout))

	// ==================== poll command ====================
	s.mux.Get(pmbGetPoll, http.HandlerFunc(handleMbGetPoll))
	s.mux.Post(pmbCreatePoll, http.HandlerFunc(handleMbCreatePoll))
	s.mux.Put(pmbUpdatePoll, http.HandlerFunc(handleMbUpdatePoll))
	s.mux.Delete(pmbDeletePoll, http.HandlerFunc(handleMbDeletePoll))
	s.mux.Post(pmbTogglePoll, http.HandlerFunc(handleMbTogglePoll))
	s.mux.Get(pmbGetPollHistory, http.HandlerFunc(handleMbGetPollHistory))
	s.mux.Get(pmbGetPolls, http.HandlerFunc(handleMbGetPolls))
	s.mux.Delete(pmbDeletePolls, http.HandlerFunc(handleMbDeletePolls))
	s.mux.Post(pmbTogglePolls, http.HandlerFunc(handleMbTogglePolls))
	s.mux.Get(pmbExportPolls, http.HandlerFunc(handleMbExportPolls))
	s.mux.Post(pmbImportPolls, http.HandlerFunc(handleMbImportPolls))

	// ==================== filter command ====================
	s.mux.Get(pmbGetFilter, http.HandlerFunc(handleMbGetFilter))
	s.mux.Post(pmbCreateFilter, http.HandlerFunc(handleMbCreateFilter))
	s.mux.Put(pmbUpdateFilter, http.HandlerFunc(handleMbUpdateFilter))
	s.mux.Delete(pmbDeleteFilter, http.HandlerFunc(handleMbDeleteFilter))
	s.mux.Post(pmbToggleFilter, http.HandlerFunc(handleMbToggleFilter))
	s.mux.Get(pmbGetFilters, http.HandlerFunc(handleMbGetFilters))
	s.mux.Delete(pmbDeleteFilters, http.HandlerFunc(handleMbDeleteFilters))
	s.mux.Post(pmbToggleFilters, http.HandlerFunc(handleMbToggleFilters))
	s.mux.Get(pmbExportFilters, http.HandlerFunc(handleMbExportFilters))
	s.mux.Post(pmbImportFilters, http.HandlerFunc(handleMbImportFilters))
}
