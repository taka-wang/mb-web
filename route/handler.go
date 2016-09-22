package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-zoo/bone"
	"github.com/taka-wang/psmb"
	"github.com/taka-wang/psmb/viper-conf"
)

// handle404NotFound handle 404 not found
func handle404NotFound(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNotFound)
	io.WriteString(rw, `{"status": "404 NOT FOUND!"}`)
}

// handleRoot handle root
func handleRoot(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	io.WriteString(rw, fmt.Sprintf(`{"version": "%s"}`, version))
}

// handleMbOnceRead [GET] request
func handleMbOnceRead(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	queries := bone.GetAllQueries(req)

	// common queries
	var (
		ip, port  string
		slave     uint8
		addr, len uint16
	)

	// check ip
	if queries["ip"] == nil {
		conf.Log.Debug("Slave IP is empty")
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, `{"status": "Invalid query string: ip"}`)
		return
	}
	ip = queries["ip"][0]

	// check port
	if queries["port"] == nil {
		port = "502" // set default slave port
	} else {
		port = queries["port"][0]
	}

	// check slave
	if queries["slave"] == nil {
		slave = 1 // set default slave id
	} else {
		if v, err := strconv.ParseUint(queries["slave"][0], 10, 8); err == nil {
			slave = uint8(v)
		} else {
			conf.Log.WithError(err).Warn("Invalid slave value")
			rw.WriteHeader(http.StatusBadRequest)
			io.WriteString(rw, `{"status": "Invalid query string: slave"}`)
			return
		}
	}

	// check addr
	if queries["addr"] == nil {
		addr = 0 // set default start address
	} else {
		if v, err := strconv.ParseUint(queries["addr"][0], 10, 16); err == nil {
			addr = uint16(v)
		} else {
			conf.Log.WithError(err).Warn("Invalid address value")
			rw.WriteHeader(http.StatusBadRequest)
			io.WriteString(rw, `{"status": "Invalid query string: addr"}`)
			return
		}
	}

	// check len
	if queries["len"] == nil {
		len = 1 // set default length
	} else {
		if v, err := strconv.ParseUint(queries["len"][0], 10, 16); err == nil {
			len = uint16(v)
		} else {
			conf.Log.WithError(err).Warn("Invalid length value")
			rw.WriteHeader(http.StatusBadRequest)
			io.WriteString(rw, `{"status": "Invalid query string: length"}`)
			return
		}
	}

	// get function code from path and check the boundary
	fc, _ := strconv.Atoi(bone.GetValue(req, "fc"))
	switch fc {
	case 1, 2:
		// call dispatcher
		tid := time.Now().UTC().UnixNano()
		js, _ := json.Marshal(psmb.MbtcpReadReq{
			From:  iam,
			Tid:   tid,
			FC:    fc,
			IP:    ip,
			Port:  port,
			Slave: slave,
			Addr:  addr,
			Len:   len,
		})
		requestHandler(tid, psmb.CmdMbtcpOnceRead, string(js), rw)
	case 3, 4:
		// FC3, FC4 only query
		var (
			rtype      psmb.RegValueType
			order      psmb.Endian
			a, b, c, d float64
		)

		// check type
		if queries["type"] == nil {
			rtype = psmb.RegisterArray // set default register type
		} else {
			if v, err := strconv.Atoi(queries["type"][0]); err == nil {
				if v < 1 || v > 8 {
					v = 1 // set to valid range
				}
				rtype = psmb.RegValueType(v)
			} else {
				conf.Log.WithError(err).Warn("Invalid type value")
				rw.WriteHeader(http.StatusBadRequest)
				io.WriteString(rw, `{"status": "Invalid query string: type"}`)
				return
			}
		}

		invalid := false
		// handle different register type
		switch rtype {
		case psmb.Scale:
			// get scale ranges
			if queries["a"] == nil {
				invalid = true
				goto READ_ERROR
			} else {
				if v, err := strconv.ParseFloat(queries["a"][0], 64); err == nil {
					a = v
				} else {
					conf.Log.WithError(err).Warn("Invalid type value")
					invalid = true
					goto READ_ERROR
				}
			}

			if queries["b"] == nil {
				invalid = true
				goto READ_ERROR
			} else {
				if v, err := strconv.ParseFloat(queries["b"][0], 64); err == nil {
					b = v
				} else {
					conf.Log.WithError(err).Warn("Invalid type value")
					invalid = true
					goto READ_ERROR
				}
			}

			if queries["c"] == nil {
				invalid = true
				goto READ_ERROR
			} else {
				if v, err := strconv.ParseFloat(queries["c"][0], 64); err == nil {
					c = v
				} else {
					conf.Log.WithError(err).Warn("Invalid type value")
					invalid = true
					goto READ_ERROR
				}
			}

			if queries["d"] == nil {
				invalid = true
				goto READ_ERROR
			} else {
				if v, err := strconv.ParseFloat(queries["d"][0], 64); err == nil {
					d = v
				} else {
					conf.Log.WithError(err).Warn("Invalid type value")
					invalid = true
					goto READ_ERROR
				}
			}

		READ_ERROR: // error label
			if invalid {
				rw.WriteHeader(http.StatusBadRequest)
				io.WriteString(rw, `{"status": "Invalid query string: scale type without args"}`)
				return
			}

			// call dispatcher
			tid := time.Now().UTC().UnixNano()
			js, _ := json.Marshal(psmb.MbtcpReadReq{
				From:  iam,
				Tid:   tid,
				FC:    fc,
				IP:    ip,
				Port:  port,
				Slave: slave,
				Addr:  addr,
				Len:   len,
				Type:  rtype,
				Range: &psmb.ScaleRange{
					DomainLow:  a,
					DomainHigh: b,
					RangeLow:   c,
					RangeHigh:  d,
				},
			})
			requestHandler(tid, psmb.CmdMbtcpOnceRead, string(js), rw)
		case psmb.UInt16, psmb.Int16, psmb.UInt32, psmb.Int32, psmb.Float32:
			// check order
			if queries["order"] == nil {
				order = psmb.BigEndian
			} else {
				if v, err := strconv.Atoi(queries["order"][0]); err == nil {
					if v < 1 || v > 4 {
						v = 1 // set to valid range
					}
					order = psmb.Endian(v)
				} else {
					conf.Log.WithError(err).Warn("Invalid order value")
					rw.WriteHeader(http.StatusBadRequest)
					io.WriteString(rw, `{"status": "Invalid query string: order"}`)
					return
				}
			}
			// call dispatcher
			tid := time.Now().UTC().UnixNano()
			js, _ := json.Marshal(psmb.MbtcpReadReq{
				From:  iam,
				Tid:   tid,
				FC:    fc,
				IP:    ip,
				Port:  port,
				Slave: slave,
				Addr:  addr,
				Len:   len,
				Type:  rtype,
				Order: order,
			})
			requestHandler(tid, psmb.CmdMbtcpOnceRead, string(js), rw)
		default: // psmb.RegisterArray, psmb.HexString
			// call dispatcher
			tid := time.Now().UTC().UnixNano()
			js, _ := json.Marshal(psmb.MbtcpReadReq{
				From:  iam,
				Tid:   tid,
				FC:    fc,
				IP:    ip,
				Port:  port,
				Slave: slave,
				Addr:  addr,
				Len:   len,
				Type:  rtype,
			})
			requestHandler(tid, psmb.CmdMbtcpOnceRead, string(js), rw)
		}
	default: // invalid function code
		conf.Log.WithField("FC", fc).Debug("Invalid function code for once read")
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, `{"status": "bad request"}`)
	}
}

// handleMbOnceWrite [POST] request
func handleMbOnceWrite(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	fc, _ := strconv.Atoi(bone.GetValue(req, "fc"))
	switch fc {
	case 5, 6, 15, 16:
		// check content-type
		if strings.HasPrefix(req.Header["Content-Type"][0], "application/json") {
			// partial
			var data json.RawMessage
			body := psmb.MbtcpWriteReq{Data: &data}
			// convert stream to byte
			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)

			if err := json.Unmarshal(buf.Bytes(), &body); err != nil {
				conf.Log.WithError(err).Warn("Fail to unmarshal POST body")
				js, _ := json.Marshal(psmb.MbtcpSimpleRes{Status: err.Error()})
				rw.WriteHeader(http.StatusUnsupportedMediaType)
				rw.Write(js)
				return
			}

			tid := time.Now().UTC().UnixNano()
			body.From = iam
			body.Tid = tid
			body.FC = fc
			js, _ := json.Marshal(body)
			requestHandler(tid, psmb.CmdMbtcpOnceWrite, string(js), rw)
		} else {
			rw.WriteHeader(http.StatusUnsupportedMediaType)
			io.WriteString(rw, `{"status": "Invalid content type"}`)
		}
	default: // invalid function code
		conf.Log.WithField("FC", fc).Debug("Invalid function code for once write")
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, `{"status": "bad request"}`)
	}
}

// handleMbGetTimeout [GET] request
func handleMbGetTimeout(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	// call dispatcher
	tid := time.Now().UTC().UnixNano()
	js, _ := json.Marshal(psmb.MbtcpTimeoutReq{From: iam, Tid: tid})
	requestHandler(tid, psmb.CmdMbtcpGetTimeout, string(js), rw)
}

// handleMbSetTimeout [POST] request
func handleMbSetTimeout(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	// check content-type
	if strings.HasPrefix(req.Header["Content-Type"][0], "application/json") {
		var body psmb.MbtcpTimeoutReq
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			conf.Log.WithError(err).Warn("Fail to unmarshal POST body")
			js, _ := json.Marshal(psmb.MbtcpSimpleRes{Status: err.Error()})
			rw.WriteHeader(http.StatusUnsupportedMediaType)
			rw.Write(js)
			return
		}

		// check timeout
		if body.Data < minConnTimeout {
			rw.WriteHeader(http.StatusBadRequest)
			io.WriteString(rw, `{"status": "bad request"}`)
			return
		}

		// call dispatcher
		tid := time.Now().UTC().UnixNano()
		body.From = iam
		body.Tid = tid
		js, _ := json.Marshal(body)
		requestHandler(tid, psmb.CmdMbtcpSetTimeout, string(js), rw)
	} else {
		rw.WriteHeader(http.StatusUnsupportedMediaType)
		io.WriteString(rw, `{"status": "Invalid content type"}`)
	}
}

func handleMbCreatePoll(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("create poll"))
}

func handleMbUpdatePoll(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("update poll"))
}

func handleMbGetPoll(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("get poll"))
}

func handleMbDeletePoll(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("delete poll"))
}

func handleMbTogglePoll(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("toggle poll"))
}

func handleMbGetPollHistory(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("get poll history"))
}

func handleMbGetPolls(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("get polls"))
}

func handleMbDeletePolls(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("delete polls"))
}

func handleMbTogglePolls(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("toggle polls"))
}

func handleMbImportPolls(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("import polls"))
}

func handleMbExportPolls(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("export polls"))
}

func handleMbCreateFilter(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("create filter"))
}

func handleMbUpdateFilter(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("update filter"))
}

func handleMbGetFilter(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("get filter"))
}

func handleMbDeleteFilter(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("delete filter"))
}

func handleMbToggleFilter(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("toggle filter"))
}

func handleMbGetFilters(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("get filters"))
}

func handleMbDeleteFilters(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("delete filters"))
}

func handleMbToggleFilters(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("toggle filters"))
}

func handleMbImportFilters(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("import filters"))
}

func handleMbExportFilters(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte("export filters"))
}

// ========================

/*
func Handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("trigger1")
	// Get the value of the "id" parameters.
	//val := bone.GetValue(req, "id")
	val := bone.GetQuery(req, "id")
	fmt.Println(val)
	val2 := bone.GetAllQueries(req)
	fmt.Println(val2)
	val3 := bone.GetValue(req, "hello")
	fmt.Println(val3)
	//rw.Write([]byte(val))
	rw.Write([]byte("hello"))
}

func Handler2(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("trigger2")

	// set response type
	rw.Header().Set("Content-Type", "application/json")

	if strings.HasPrefix(req.Header["Content-Type"][0], "application/json") {
		decoder := json.NewDecoder(req.Body)
		var t timeoutReq
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(t.Data)
		rw.Write([]byte("world"))
	} else {
		resp := simpleRes{Status: "Invalid content type"}
		js, _ := json.Marshal(resp)
		rw.WriteHeader(http.StatusUnsupportedMediaType)
		rw.Write(js)
	}
}

func Handler3(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("trigger3")

	resp := timeoutRes{Status: "ok", Data: 123456}

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)

}
*/
/*

func Another(rw *http.ResponseWriter) {
	rw2 := *rw
	rw2.Header().Set("Content-Type", "application/json")
	rw2.Write([]byte("hello another"))
}

pass rw to zmq queue
once got the response from psmb, send response
otherwise, timeout timer will trigger timeout response

*/
