package backend

import (
	"fmt"
	"net/http"

	zmq "github.com/takawang/zmq3"
)

const (
	_ dataSource = iota
	upstream
	downstream
)

type (
	// dataSource data source: http or zmq
	dataSource int

	// job
	job struct {
		source dataSource
		msg    []string
	}

	// worker in worker pool
	worker struct {
		id      int
		service *Service
	}

	// Service service
	Service struct {
		jobChan chan job
		pub     *zmq.Socket
		sub     *zmq.Socket
	}
)

// NewService create service
func NewService() (*Service, error) {
	sender, _ := zmq.NewSocket(zmq.PUB)
	receiver, _ := zmq.NewSocket(zmq.SUB)
	return &Service{
		pub: sender,
		sub: receiver,
	}, nil
}

// Start start service
func (b *Service) Start() {
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

	// start sub
	b.subscriber()
}

func (b *Service) startZMQ() {
	fmt.Println("Start ZMQ")
	// pub
	if err := b.pub.Connect(pubEndpoint); err != nil {
		fmt.Println(err)
		return
	}
	// sub
	if err := b.sub.Connect(subEndpoint); err != nil {
		fmt.Println(err)
		return
	}
	if err := b.sub.SetSubscribe(""); err != nil {
		fmt.Println(err)
		return
	}
}

func (w worker) process(j job) {
	switch j.source {
	case upstream:
		fmt.Println("upstream")
		for {
			w.service.pub.Send(j.msg[0], zmq.SNDMORE)
			w.service.pub.Send(j.msg[1], 0)
			break
		}
	default:
		// downstream
		fmt.Println("downstream")
	}
}

// dispatch create job and push it to the job channel
func (b *Service) dispatch(source dataSource, msg []string) {
	job := job{source, msg}
	go func() {
		b.jobChan <- job
	}()
}

// generic subscribe
func (b *Service) subscriber() {
	doLoop := true
	fmt.Println("listen..")
	for doLoop {
		if msg, err := b.sub.RecvMessage(0); err == nil {
			if len(msg) == 2 {
				fmt.Println(msg[0], msg[1])
				b.dispatch(downstream, msg)
			} else {
				fmt.Println("strange response")
				fmt.Println(msg)
				// discard
			}
		}
	}
}

// generic tcp publisher
func (b *Service) requestHandler(cmd, json string, w http.ResponseWriter, r *http.Request) {
	msg := []string{cmd, json}
	b.dispatch(upstream, msg)
}

/*
func handle(msg []string) {
	switch msg[0] {
	case tcp.mbOnceRead:
		//
	case tcp.mbOnceWrite:
		//
	case tcp.mbGetTimeout:
		//
	case tcp.mbSetTimeout:
		//
	default:
		//
	}
}
*/
