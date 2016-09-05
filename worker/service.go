// Package worker a zmq backend worker queue
//
// By taka@cmwang.net
//
package worker

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/taka-wang/mb-web/conf"
	"github.com/taka-wang/psmb"
	zmq "github.com/takawang/zmq3"
)

var (
	// pubEndpoint zmq pub endpoint
	pubEndpoint string
	// subEndpoint zmq sub endpoint
	subEndpoint string
	// maxQueueSize the size of job queue
	maxQueueSize int
	// maxWorkers the number of workers to start
	maxWorkers int
)

func setDefaults() {
	conf.SetDefault(keyWebPub, defaultWebPub)
	conf.SetDefault(keyWebSub, defaultWebSub)
	conf.SetDefault(keyMaxWorker, defaultMaxWorker)
	conf.SetDefault(keyMaxQueue, defaultMaxQueue)
}

func init() {
	// set default values
	setDefaults()

	// get init values from config file
	pubEndpoint = conf.GetString(keyWebPub)
	subEndpoint = conf.GetString(keyWebSub)
	maxWorkers = conf.GetInt(keyMaxWorker)
	maxQueueSize = conf.GetInt(keyMaxQueue)
}

type (

	// dataSource data source: http or zmq
	dataSource int

	// job
	job struct {
		source  dataSource
		cmd     string
		payload string
	}

	// worker in worker pool
	worker struct {
		id      int
		service *Service
	}

	// Service service
	Service struct {
		sync.RWMutex
		// jobChan job channel
		jobChan chan job
		// pub ZMQ publisher endpoints
		pub *zmq.Socket
		// sub ZMQ subscriber endpoints
		sub *zmq.Socket
		// httpMap http req map
		httpMap map[int64]http.ResponseWriter
		// isRunning running flag
		isRunning bool
		// isStopped stop channel
		isStopped chan bool
	}
)

// NewService create service
func NewService() (*Service, error) {
	sender, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		conf.Log.WithError(err).Error("Fail to create sender")
		return nil, err
	}

	receiver, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		conf.Log.WithError(err).Error("Fail to create receiver")
		return nil, err
	}

	return &Service{
		isRunning: false,
		pub:       sender,
		sub:       receiver,
		isStopped: make(chan bool),
		httpMap:   make(map[int64]http.ResponseWriter),
	}, nil
}

// Start start service
func (b *Service) start() {
	b.Lock()
	defer b.Unlock()

	conf.Log.Debug("Start Worker Queue")

	// only start the service if it hasn't been started yet.
	if b.isRunning {
		conf.Log.Warn("Already Running")
		return
	}

	b.isRunning = true
	b.startZMQ()

	// init the job channel
	b.jobChan = make(chan job, maxQueueSize)

	// create workers
	for i := 0; i < maxWorkers; i++ {
		w := worker{i, b}
		go func(w worker) {
			for j := range b.jobChan {
				w.process(j)
			}
		}(w)
	}

	// start listening
	ticker := time.NewTicker(200 * time.Millisecond)
	go func() {
		for {
			select {
			case <-b.isStopped:
				b.isRunning = false
				return
			case <-ticker.C:
				if b.isRunning {
					if msg, err := b.sub.RecvMessage(zmq.DONTWAIT); err == nil {
						if len(msg) == 2 {
							conf.Log.WithFields(conf.Fields{
								"cmd":     msg[0],
								"payload": msg[1],
							}).Debug("Recv response from psmb")
							b.dispatch(downstream, msg[0], msg[1]) // send to worker queue
						} else {
							conf.Log.WithField("msg", msg).Error("Invalid message length, discard!!")
						}
					}
					//conf.Log.Debug("waiting..")
				}
			}
		}
	}()
}

func (b *Service) stop() {
	b.Lock()
	defer b.Unlock()

	conf.Log.Debug("Stop Worker Queue")

	if b.isRunning {
		b.isStopped <- true
		b.stopZMQ()
		close(b.jobChan) // close job channel and wait for workers to complete
	}
}

func (b *Service) stopZMQ() {
	conf.Log.Debug("Stop ZMQ")
	if err := b.pub.Disconnect(pubEndpoint); err != nil {
		conf.Log.WithError(err).Debug("Fail to disconnect from publisher endpoint")
	}
	if err := b.sub.Disconnect(subEndpoint); err != nil {
		conf.Log.WithError(err).Debug("Fail to disconnect from subscribe endpoint")
	}
}

func (b *Service) startZMQ() {
	conf.Log.Debug("Start ZMQ")

	// publish to psmb
	if err := b.pub.Connect(pubEndpoint); err != nil {
		conf.Log.WithError(err).Error("Fail to connect to publisher endpoint")
		return
	}
	// subscribe from psmb
	if err := b.sub.Connect(subEndpoint); err != nil {
		conf.Log.WithError(err).Error("Fail to connect to subscriber endpoint")
		return
	}
	// filter
	if err := b.sub.SetSubscribe(""); err != nil {
		conf.Log.WithError(err).Error("Fail to set subscriber's filter")
		return
	}
}

// dispatch create job and push it to the job channel
func (b *Service) dispatch(source dataSource, command, payload string) {
	job := job{source, command, payload}
	go func() {
		b.Lock()
		defer b.Unlock()
		b.jobChan <- job
	}()
}

// process handle request
func (w worker) process(j job) {
	switch j.source {
	case upstream:
		conf.Log.WithFields(conf.Fields{
			"cmd":     j.cmd,
			"payload": j.payload,
		}).Debug("Processing frontend request")
		w.service.sendRequest(j.cmd, j.payload)
	default:
		conf.Log.WithFields(conf.Fields{
			"cmd":     j.cmd,
			"payload": j.payload,
		}).Debug("Processing psmb response")
		w.service.sendResponse(j.cmd, j.payload)
	}
}

// requestHandler generic tcp publisher for callback setting
func (b *Service) requestHandler(ts int64, command, payload string, w http.ResponseWriter) {
	// set http response handler to map with mutex lock
	b.Lock()
	defer b.Unlock()

	b.httpMap[ts] = w
	b.dispatch(upstream, command, payload) // send to worker queue
}

func (b *Service) sendRequest(command, payload string) {
	for {
		b.pub.Send(command, zmq.SNDMORE)
		b.pub.Send(payload, 0)
		break
	}
}

func (b *Service) sendResponse(command, payload string) {
	b.Lock()
	defer b.Unlock()

	switch command {
	case psmb.CmdMbtcpOnceRead:

		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(payload), &res); err != nil {
			conf.Log.WithError(err).Error(ErrUnmarshal.Error())
			return
		}

		// try to retrieve writer from map
		if writer, ok := b.httpMap[res.Tid]; !ok {
			conf.Log.Warn("Fail to retrieve http writer from map.")
		} else {
			// remove http writer from map
			delete(b.httpMap, res.Tid)

			// check response status
			var resp psmb.MbtcpReadRes
			if res.Status != "ok" {
				writer.WriteHeader(http.StatusBadRequest)
				resp = psmb.MbtcpReadRes{Status: res.Status}
			} else {
				resp = psmb.MbtcpReadRes{
					Status: res.Status,
					Type:   res.Type,
					Bytes:  res.Bytes,
					Data:   res.Data,
				}
			}

			// send response to http client
			js, _ := json.Marshal(resp)
			writer.Write(js)
		}
	case psmb.CmdMbtcpOnceWrite:
		// unmarshal response from psmb
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(payload), &res); err != nil {
			conf.Log.WithError(err).Error(ErrUnmarshal.Error())
			return
		}

		// try to retrieve writer from map
		if writer, ok := b.httpMap[res.Tid]; !ok {
			conf.Log.Warn("Fail to retrieve http writer from map.")
		} else {
			// remove http writer from map
			delete(b.httpMap, res.Tid)

			// check response status
			var resp psmb.MbtcpSimpleRes
			if res.Status != "ok" {
				writer.WriteHeader(http.StatusBadRequest)
				resp = psmb.MbtcpSimpleRes{Status: res.Status}
			} else {
				resp = psmb.MbtcpSimpleRes{Status: res.Status}
			}

			// send response to http client
			js, _ := json.Marshal(resp)
			writer.Write(js)
		}
	case psmb.CmdMbtcpSetTimeout:
		// unmarshal response from psmb
		var res psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(payload), &res); err != nil {
			conf.Log.WithError(err).Error(ErrUnmarshal.Error())
			return
		}

		// try to retrieve writer from map
		if writer, ok := b.httpMap[res.Tid]; !ok {
			conf.Log.Warn("Fail to retrieve http writer from map.")
		} else {
			// remove http writer from map
			delete(b.httpMap, res.Tid)

			// check response status
			var resp psmb.MbtcpSimpleRes
			if res.Status != "ok" {
				writer.WriteHeader(http.StatusBadRequest)
				resp = psmb.MbtcpSimpleRes{Status: res.Status}
			} else {
				resp = psmb.MbtcpSimpleRes{Status: res.Status}
			}

			// send response to http client
			js, _ := json.Marshal(resp)
			writer.Write(js)
		}
	case psmb.CmdMbtcpGetTimeout:
		// unmarshal response from psmb
		var res psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(payload), &res); err != nil {
			conf.Log.WithError(err).Error(ErrUnmarshal.Error())
			return
		}

		// try to retrieve writer from map
		if writer, ok := b.httpMap[res.Tid]; !ok {
			conf.Log.Warn("Fail to retrieve http writer from map.")
		} else {
			// remove http writer from map
			delete(b.httpMap, res.Tid)

			// check response status
			var resp psmb.MbtcpTimeoutRes
			if res.Status != "ok" {
				writer.WriteHeader(http.StatusBadRequest)
				resp = psmb.MbtcpTimeoutRes{Status: res.Status}
			} else {
				resp = psmb.MbtcpTimeoutRes{Status: res.Status, Data: res.Data}
			}

			// send response to http client
			js, _ := json.Marshal(resp)
			writer.Write(js)
		}
	default:
		conf.Log.Warn(ErrResponseNotSupport.Error())
	}
}
