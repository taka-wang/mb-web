package worker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-zoo/bone"
	"github.com/takawang/sugar"
	zmq "github.com/takawang/zmq3"
)

var sender, receiver *zmq.Socket

func init() {
	// setup sender
	sender, _ = zmq.NewSocket(zmq.PUB)
	sender.Bind("ipc:///tmp/from.psmb")
	// setup receiver
	receiver, _ = zmq.NewSocket(zmq.SUB)
	receiver.Bind("ipc:///tmp/to.psmb")
	filter := ""
	receiver.SetSubscribe(filter) // filter frame 1
	// wait for zmq ready
	time.Sleep(time.Duration(1) * time.Second)
}

// generic tcp publisher
func publisher(cmd, json string) {
	for {
		time.Sleep(time.Duration(10) * time.Millisecond)
		t := time.Now()
		fmt.Println("pub:", t.Format("2006-01-02 15:04:05.000"))
		sender.Send(cmd, zmq.SNDMORE) // frame 1
		sender.Send(json, 0)          // convert to string; frame 2
		// send the exit loop
		break
	}
}

// generic subscribe
func subscriber() {
	for {
		msg, _ := receiver.RecvMessage(0)
		t := time.Now()
		fmt.Println("sub:", t.Format("2006-01-02 15:04:05.000"))
		fmt.Println(msg)
		// send
		publisher(msg[0], msg[1])
		break
	}
}

func TestWorker(t *testing.T) {
	s := sugar.New(t)

	s.Assert("E2E test", func(logf sugar.Log) bool {
		go subscriber()

		Start()
		time.Sleep(time.Duration(3) * time.Second)

		// send http request
		r, _ := http.NewRequest("GET", "/api/mb/tcp/fc/1", nil)
		w := httptest.NewRecorder()
		RequestHandler(12345, "hello", "world", w)
		time.Sleep(time.Duration(3) * time.Second)
		mux := bone.New()
		// check response
		mux.ServeHTTP(w, r)
		// Check the status code is what we expect.
		if status := w.Code; status != http.StatusOK {
			logf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `{"status": "once read"}`
		//logf(rr.Body.String())
		if w.Body.String() != expected {
			logf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}

		Stop()
		return true
	})

	s.Assert("test start service", func(logf sugar.Log) bool {
		Start()
		time.Sleep(time.Duration(3) * time.Second)
		Stop()
		time.Sleep(time.Duration(3) * time.Second)
		Start()
		time.Sleep(time.Duration(3) * time.Second)
		Stop()
		return true
	})
}
