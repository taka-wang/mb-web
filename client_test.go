package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	route "github.com/taka-wang/mb-web/route"
	"github.com/takawang/sugar"
)

func TestHandleMbOnceRead(t *testing.T) {

	s := sugar.New(t)
	route.Middleware()
	mux := route.Mux

	s.Assert("Test handleMbOnceRead", func(log sugar.Log) bool {
		fmt.Println("1")
		req, _ := http.NewRequest("GET", "/mb/tcp/fc/1", nil)
		fmt.Println("2")
		w := httptest.NewRecorder()
		fmt.Println("3")
		mux.ServeHTTP(w, req)
		fmt.Println("4")
		// Check the status code is what we expect.
		if status := w.Code; status != http.StatusOK {
			log("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `{"status": "once read"}`
		//log(rr.Body.String())
		if w.Body.String() != expected {
			log("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
			return false
		}
		return true
	})

	s.Assert("Test handleMbOnceRead2", func(log sugar.Log) bool {
		req, _ := http.NewRequest("GET", "/mb/tcp/fc/5", nil)

		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		// Check the status code is what we expect.
		if status := w.Code; status != http.StatusOK {
			log("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `{"status": "once read"}`
		//log(rr.Body.String())
		if w.Body.String() != expected {
			log("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
			return false
		}
		return true
	})
}
