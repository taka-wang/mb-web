package route

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-zoo/bone"
)

func setJSONheader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
}

func handleMbOnceRead(rw http.ResponseWriter, req *http.Request) {
	setJSONheader(&rw)
	fc, _ := strconv.Atoi(bone.GetValue(req, "fc"))
	switch fc {
	case 1, 2, 3, 4:
		// do nothing
	default:
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, `{"status": "bad request"}`)
		return
	}

	io.WriteString(rw, `{"status": "once read"}`)
}

func handleMbOnceWrite(rw http.ResponseWriter, req *http.Request) {
	setJSONheader(&rw)
	fc, _ := strconv.Atoi(bone.GetValue(req, "fc"))
	switch fc {
	case 5, 6, 15, 16:
		// do nothing
	default:
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, `{"status": "bad request"}`)
		return
	}

	rw.Write([]byte("once write"))
}

// handleMbGetTimeout
// GET handler
func handleMbGetTimeout(rw http.ResponseWriter, req *http.Request) {
	setJSONheader(&rw)

	resp := timeoutRes{Status: "ok", Data: 123456}

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	rw.Write(js)
}

// handleMbSetTimeout
// POST handler
func handleMbSetTimeout(rw http.ResponseWriter, req *http.Request) {
	setJSONheader(&rw)
	// check content-type
	if strings.HasPrefix(req.Header["Content-Type"][0], "application/json") {
		var body timeoutReq
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			js, _ := json.Marshal(simpleRes{Status: err.Error()})
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

		// post to backend
		//service.requestHandler("hello", "world", rw, req)

		fmt.Println(body.Data)
		rw.Write([]byte("set timeout"))
	} else {
		rw.WriteHeader(http.StatusUnsupportedMediaType)
		io.WriteString(rw, `{"status": "Invalid content type"}`)
		return
	}
}

func handleMbCreatePoll(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("create poll")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("create poll"))
}

func handleMbUpdatePoll(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("update poll")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("update poll"))
}

func handleMbGetPoll(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("get poll")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("get poll"))
}

func handleMbDeletePoll(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("delete poll")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("delete poll"))
}

func handleMbTogglePoll(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("toggle poll")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("toggle poll"))
}

func handleMbGetPollHistory(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("get poll history")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("get poll history"))
}

func handleMbGetPolls(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("get polls")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("get polls"))
}

func handleMbDeletePolls(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("delete polls")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("delete polls"))
}

func handleMbTogglePolls(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("toggle polls")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("toggle polls"))
}

func handleMbImportPolls(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("import polls")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("import polls"))
}

func handleMbExportPolls(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("export polls")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("export polls"))
}

func handleMbCreateFilter(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("create filter")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("create filter"))
}

func handleMbUpdateFilter(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("update filter")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("update filter"))
}

func handleMbGetFilter(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("get filter")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("get filter"))
}

func handleMbDeleteFilter(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("delete filter")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("delete filter"))
}

func handleMbToggleFilter(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("toggle filter")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("toggle filter"))
}

func handleMbGetFilters(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("get filters")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("get filters"))
}

func handleMbDeleteFilters(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("delete filters")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("delete filters"))
}

func handleMbToggleFilters(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("toggle filters")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("toggle filters"))
}

func handleMbImportFilters(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("import filters")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("import filters"))
}

func handleMbExportFilters(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("export filters")

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
otherwise, timeout timer will trigger timeout reponse

*/
