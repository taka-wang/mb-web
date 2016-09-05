package dummy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-zoo/bone"
	route "github.com/taka-wang/mb-web/route"
	worker "github.com/taka-wang/mb-web/worker"
	"github.com/taka-wang/psmb"
	"github.com/takawang/sugar"
)

var mux *bone.Mux
var hostName string
var portNum1 = "502"
var portNum2 = "503"

func init() {
	// start worker
	worker.Start()
	// setup route
	route.SetRoute()
	mux = route.GetMux()
	// link worker and route packages
	route.SetRequestHandler(worker.RequestHandler)

	// setup docker
	time.Sleep(2000 * time.Millisecond)
	// generalize host reslove for docker/local env
	host, err := net.LookupHost("slave")
	if err != nil {
		fmt.Println("Local run")
		hostName = "127.0.0.1"
	} else {
		fmt.Println("Docker run")
		hostName = host[0] //docker
	}
}

func TestTimeoutOps(t *testing.T) {

	s := sugar.New(t)

	s.Assert("Get timeout `1st` round", func(logf sugar.Log) bool {
		req, _ := http.NewRequest("GET", "/api/mb/tcp/timeout", nil)
		writer := httptest.NewRecorder()
		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("Set `valid` timeout", func(logf sugar.Log) bool {
		var jsonStr = []byte(`{ "timeout": 212345 }`)
		req, _ := http.NewRequest("POST", "/api/mb/tcp/timeout", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("Get timeout `2nd` round", func(logf sugar.Log) bool {
		req, _ := http.NewRequest("GET", "/api/mb/tcp/timeout", nil)
		writer := httptest.NewRecorder()
		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("Set `invalid` timeout", func(logf sugar.Log) bool {
		var jsonStr = []byte(`{ "timeout": 123 }`)
		req, _ := http.NewRequest("POST", "/api/mb/tcp/timeout", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

}

func TestOneOffWriteFC5(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`mbtcp.once.write FC5` write bit test: port 502 - invalid value(2) - (1/4)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"data": 2
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/5", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC5` write bit test: port 502 - miss port - (2/4)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"slave": 1,
				"addr": 10,
				"data": 2
			}`, hostName))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/5", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC5` write bit test: port 502 - valid value(0) - (3/4)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"data": 0
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/5", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC5` write bit test: port 502 - valid value(1) - (4/4)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"data": 1
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/5", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write - invalid function code", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"data": 1
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/7", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})
}

func TestOneOffWriteFC6(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`mbtcp.once.write FC6` write `DEC` register test: port 502 - valid value (22) - (1/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"hex": false,
				"data": "22"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/6", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `DEC` register test: port 502 - miss hex type & port - (2/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"slave": 1,
				"addr": 10,
				"data": "22"
			}`, hostName))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/6", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `DEC` register test: port 502 - invalid value (array) - (3/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"hex": false,
				"data": "22,11"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/6", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `DEC` register test: port 502 - invalid hex type - (4/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"hex": false,
				"data": "ABCD1234"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/6", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `HEX` register test: port 502 - valid value (ABCD) - (5/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"hex": true,
				"data": "ABCD"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/6", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `HEX` register test: port 502 - miss port (ABCD) - (6/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"hex": true,
				"data": "ABCD"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/6", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `HEX` register test: port 502 - invalid value (ABCD1234) - (7/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"hex": true,
				"data": "ABCD1234"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/6", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `HEX` register test: port 502 - invalid hex type - (8/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"hex": true,
				"data": "22,11"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/6", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})
}

func TestOneOffWriteFC15(t *testing.T) {

	s := sugar.New(t)

	s.Assert("`mbtcp.once.write FC15` write bit test: port 502 - invalid json type - (1/5)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
                "data": [-1,0,-1,0]
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/15", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC15` write bit test: port 502 - invalid json type - (2/5)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
                "data": "1,0,1,0"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/15", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC15` write bit test: port 502 - invalid value(2) - (3/5)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
                "data": [2,0,2,0]
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/15", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC15` write bit test: port 502 - miss from & port - (4/5)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
                "data": [2,0,2,0]
			}`, hostName))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/15", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC15` write bit test: port 502 - valid value(0) - (5/5)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
                "data": [0,1,1,0]
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/15", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})
}

func TestOneOffWriteFC16(t *testing.T) {

	s := sugar.New(t)

	s.Assert("`mbtcp.once.write FC16` write `DEC` register test: port 502 - valid value (11,22,33,44) - (1/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
				"hex": false,
				"data": "11,22,33,44"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/16", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `DEC` register test: port 502 - miss hex type & port - (2/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
				"data": "11,22,33,44"
			}`, hostName))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/16", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `DEC` register test: port 502 - invalid hex type - (3/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
				"hex": false,
				"data": "ABCD1234"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/16", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `DEC` register test: port 502 - invalid length - (4/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"len": 8,
				"hex": false,
				"data": "11,22,33,44"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/16", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `HEX` register test: port 502 - valid value (ABCD1234) - (5/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
				"hex": true,
				"data": "ABCD1234"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/16", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `HEX` register test: port 502 - miss port (ABCD) - (6/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
				"hex": true,
				"data": "ABCD1234"
			}`, hostName))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/16", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `HEX` register test: port 502 - invalid hex type (11,22,33,44) - (7/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"slave": 1,
				"addr": 10,
				"len": 4,
				"hex": true,
				"data": "11,22,33,44"
			}`, hostName))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/16", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `HEX` register test: port 502 - invalid length - (8/8)", func(logf sugar.Log) bool {
		var jsonStr = []byte(fmt.Sprintf(`{
				"ip": "%s",
				"port": "%s",
				"slave": 1,
				"addr": 10,
				"len": 8,
				"hex": true,
				"data": "ABCD1234"
			}`, hostName, portNum1))

		req, _ := http.NewRequest("POST", "/api/mb/tcp/fc/16", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		logf("POST `%v`", string(jsonStr))

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

}

func TestOneOffOperations(t *testing.T) {
	s := sugar.New(t)

	s.Assert("Get `valid` function code - fc1", func(logf sugar.Log) bool {
		return true
	})
	s.Assert("Get `invalid` function code - fc5", func(logf sugar.Log) bool {
		return true
	})
	/*
		s.Assert("Test handleMbOnceRead", func(logf sugar.Log) bool {

			var url = fmt.Sprintf("/api/mb/tcp/fc/1?ip=%s&port=%s&slave=1&addr=10&len=4", hostName, portNum1)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Content-Type", "application/json")
			writer := httptest.NewRecorder()

			time.Sleep(time.Duration(2) * time.Second)
			mux.ServeHTTP(writer, req)
			// wait for response
			time.Sleep(time.Duration(2) * time.Second)

			// Check the status code is what we expect.
			if status := writer.Code; status != http.StatusOK {
				logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
			}

			// Check the response body is what we expect.
			var res psmb.MbtcpTimeoutRes
			if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
				logf("Fail to Unmarshal: `%s`", writer.Body.String())
				return false
			}

			if res.Status != "ok" {
				logf("Handler returned unexpected body: got `%s`", writer.Body.String())
				return false
			}
			logf("Got desired response: `%s`", writer.Body.String())
			return true
		})
	*/
	/*
		s.Assert("Test handleMbOnceRead", func(logf sugar.Log) bool {
			req, _ := http.NewRequest("GET", "/api/mb/tcp/fc/1", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			// Check the status code is what we expect.
			if status := w.Code; status != http.StatusOK {
				logf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			// Check the response body is what we expect.
			expected := `{"status": "once read"}`
			//logf(rr.Body.String())
			if w.Body.String() != expected {
				logf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
				return false
			}
			return true
		})

		s.Assert("Test handleMbOnceRead2", func(logf sugar.Log) bool {
			req, _ := http.NewRequest("GET", "/api/mb/tcp/fc/5", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			// Check the status code is what we expect.
			if status := w.Code; status != http.StatusOK {
				logf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			// Check the response body is what we expect.
			expected := `{"status": "once read"}`
			//logf(rr.Body.String())
			if w.Body.String() != expected {
				logf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
				return true
			}
			return false
		})
	*/
}

func TestOneOffReadFC1(t *testing.T) {

	s := sugar.New(t)

	s.Assert("`FC1` read bits test: port 502 - miss ip - (1/5)", func(logf sugar.Log) bool {
		//var url = fmt.Sprintf("/api/mb/tcp/fc/1?ip=%s&port=%s&slave=1&addr=3&len=1", hostName, portNum1)
		var url = fmt.Sprintf("/api/mb/tcp/fc/1?&port=%s&slave=1&addr=3&len=1", portNum1)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC1` read bits test: port 502 - length 1 - (2/5)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/1?ip=%s&port=%s&slave=1&addr=8&len=1", hostName, portNum1)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC1` read bits test: port 502 - length 7 - (3/5)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/1?ip=%s&port=%s&slave=1&addr=3&len=7", hostName, portNum1)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC1` read bits test: port 502 - Illegal data address - (4/5)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/1?ip=%s&port=%s&slave=1&addr=20000&len=7", hostName, portNum1)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC1` read bits test: port 503 - length 7 - (5/5)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/1?ip=%s&port=%s&slave=1&addr=3&len=7", hostName, portNum2)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}

		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})
}

func TestOneOffReadFC2(t *testing.T) {

	s := sugar.New(t)

	s.Assert("`FC2` read bits test: port 502 - length 1 - (1/4)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/2?ip=%s&port=%s&slave=1&addr=3&len=1", hostName, portNum1)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC2` read bits test: port 502 - length 7 - (2/4)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/2?ip=%s&port=%s&slave=1&addr=3&len=7", hostName, portNum1)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC2` read bits test: port 502 - Illegal data address - (3/4)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/2?ip=%s&port=%s&slave=1&addr=20000&len=7", hostName, portNum1)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC2` read bits test: port 503 - length 7 - (4/4)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/2?ip=%s&port=%s&slave=1&addr=3&len=7", hostName, portNum2)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})
}

func TestOneOffReadFC3(t *testing.T) {

	s := sugar.New(t)

	s.Assert("`FC3` read bytes Type 1 test: port 502 - (1/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d", hostName, portNum1, psmb.RegisterArray)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 2 test: port 502 - (2/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d", hostName, portNum1, psmb.HexString)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 3 length 4 test: port 502 - (3/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d&a=0&b=65535&c=100&d=500", hostName, portNum1, psmb.Scale)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 3 length 7 test: port 502 - invalid length - (4/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d&a=0&b=65535&c=100&d=500", hostName, portNum1, psmb.Scale)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 4 length 4 test: port 502 - Order: AB - (5/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d&order=%d", hostName, portNum1, psmb.UInt16, psmb.AB)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 4 length 4 test: port 502 - Order: BA - (6/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d&order=%d", hostName, portNum1, psmb.UInt16, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 4 length 4 test: port 502 - miss order - (7/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d", hostName, portNum1, psmb.UInt16)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 5 length 4 test: port 502 - Order: AB - (8/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d&order=%d", hostName, portNum1, psmb.Int16, psmb.AB)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 5 length 4 test: port 502 - Order: BA - (9/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d&order=%d", hostName, portNum1, psmb.Int16, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 5 length 4 test: port 502 - miss order - (10/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d", hostName, portNum1, psmb.Int16)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 6 length 8 test: port 502 - Order: AB - (11/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.UInt32, psmb.AB)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 6 length 8 test: port 502 - Order: BA - (12/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.UInt32, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 6 length 8 test: port 502 - miss order - (13/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d", hostName, portNum1, psmb.UInt32)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 6 length 7 test: port 502 - invalid length - (14/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d&order=%d", hostName, portNum1, psmb.UInt32, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 7 length 8 test: port 502 - Order: AB - (15/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Int32, psmb.AB)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 7 length 8 test: port 502 - Order: BA - (16/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Int32, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 7 length 8 test: port 502 - miss order - (17/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d", hostName, portNum1, psmb.Int32)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 7 length 7 test: port 502 - invalid length - (18/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d&order=%d", hostName, portNum1, psmb.Int32, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 8 length 8 test: port 502 - order: ABCD - (19/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Float32, psmb.ABCD)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 8 length 8 test: port 502 - order: DCBA - (20/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Float32, psmb.DCBA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 8 length 8 test: port 502 - order: BADC - (21/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Float32, psmb.BADC)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 8 length 8 test: port 502 - order: CDAB - (22/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Float32, psmb.CDAB)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes Type 8 length 7 test: port 502 - invalid length - (23/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d&order=%d", hostName, portNum1, psmb.Float32, psmb.BigEndian)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC3` read bytes: port 502 - invalid type - (24/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/3?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, 9, psmb.BigEndian)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

}

func TestOneOffReadFC4(t *testing.T) {

	s := sugar.New(t)

	s.Assert("`FC4` read bytes Type 1 test: port 502 - (1/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d", hostName, portNum1, psmb.RegisterArray)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 2 test: port 502 - (2/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d", hostName, portNum1, psmb.HexString)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 3 length 4 test: port 502 - (3/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d&a=0&b=65535&c=100&d=500", hostName, portNum1, psmb.Scale)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 3 length 7 test: port 502 - invalid length - (4/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d&a=0&b=65535&c=100&d=500", hostName, portNum1, psmb.Scale)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 4 length 4 test: port 502 - Order: AB - (5/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d&order=%d", hostName, portNum1, psmb.UInt16, psmb.AB)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 4 length 4 test: port 502 - Order: BA - (6/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d&order=%d", hostName, portNum1, psmb.UInt16, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 4 length 4 test: port 502 - miss order - (7/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d", hostName, portNum1, psmb.UInt16)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 5 length 4 test: port 502 - Order: AB - (8/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d&order=%d", hostName, portNum1, psmb.Int16, psmb.AB)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 5 length 4 test: port 502 - Order: BA - (9/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d&order=%d", hostName, portNum1, psmb.Int16, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 5 length 4 test: port 502 - miss order - (10/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=4&type=%d", hostName, portNum1, psmb.Int16)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 6 length 8 test: port 502 - Order: AB - (11/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.UInt32, psmb.AB)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 6 length 8 test: port 502 - Order: BA - (12/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.UInt32, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 6 length 8 test: port 502 - miss order - (13/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d", hostName, portNum1, psmb.UInt32)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 6 length 7 test: port 502 - invalid length - (14/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d&order=%d", hostName, portNum1, psmb.UInt32, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 7 length 8 test: port 502 - Order: AB - (15/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Int32, psmb.AB)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 7 length 8 test: port 502 - Order: BA - (16/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Int32, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 7 length 8 test: port 502 - miss order - (17/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d", hostName, portNum1, psmb.Int32)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 7 length 7 test: port 502 - invalid length - (18/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d&order=%d", hostName, portNum1, psmb.Int32, psmb.BA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 8 test: port 502 - order: ABCD - (19/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Float32, psmb.ABCD)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 8 test: port 502 - order: DCBA - (20/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Float32, psmb.DCBA)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 8 test: port 502 - order: BADC - (21/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Float32, psmb.BADC)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 8 test: port 502 - order: CDAB - (22/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, psmb.Float32, psmb.CDAB)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 7 test: port 502 - invalid length - (23/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=7&type=%d&order=%d", hostName, portNum1, psmb.Float32, psmb.BigEndian)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		// we expect not ok response
		if res.Status == "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired error response: `%s`", writer.Body.String())
		return true
	})

	s.Assert("`FC4` read bytes: port 502 - invalid type - (24/24)", func(logf sugar.Log) bool {
		var url = fmt.Sprintf("/api/mb/tcp/fc/4?ip=%s&port=%s&slave=1&addr=3&len=8&type=%d&order=%d", hostName, portNum1, 9, psmb.BigEndian)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()

		time.Sleep(time.Duration(2) * time.Second)
		mux.ServeHTTP(writer, req)
		// wait for response
		time.Sleep(time.Duration(2) * time.Second)

		// Check the status code is what we expect.
		if status := writer.Code; status != http.StatusOK {
			logf("Handler returned wrong status code: got `%v` want `%v`", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var res psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(writer.Body.String()), &res); err != nil {
			logf("Fail to Unmarshal: `%s`", writer.Body.String())
			return false
		}

		if res.Status != "ok" {
			logf("Handler returned unexpected body: got `%s`", writer.Body.String())
			return false
		}
		logf("Got desired response: `%s`", writer.Body.String())
		return true
	})

}
