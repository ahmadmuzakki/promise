package main

import (
	"context"
	"fmt"
	"github.com/ahmadmuzakki/promise"
	"log"
	"net/http"
	"time"
)

// run this file and go to your browser
// visit localhost:9000/dosomething?timeout=4s
func main() {
	http.HandleFunc("/dosomething", handleDoSomething)
	http.ListenAndServe(":9000", nil)
}

func handleDoSomething(w http.ResponseWriter, r *http.Request) {
	timeoutStr := r.FormValue("timeout")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Println("cannot parse duration", err)
		return
	}

	ctx := r.Context()
	// timeout for the whole operation is depends on "timeout" param
	ctx, done := context.WithTimeout(ctx, timeout)
	defer done()

	// return result with Do()
	val, err := promise.Promise(ctx, dbTx).Do()
	if err != nil {
		w.Write([]byte(fmt.Sprint("dbTx ", "error:", err)))
	} else {
		w.Write([]byte(fmt.Sprint("dbTx result:", val)))
	}
	w.Write([]byte("\n"))

	// or we can do callback style like this
	promise.Promise(ctx, httpTx).Then(func(val interface{}) {
		w.Write([]byte(fmt.Sprint("httpTx result:", val)))
	}, func(err error) {
		w.Write([]byte(fmt.Sprint("httpTx error:", err)))
	})

}

// let's say we have transaction that takes sometimes here
func dbTx() (interface{}, error) {
	time.Sleep(time.Duration(1) * time.Second)
	return "success", nil
}

// another transaction
func httpTx() (interface{}, error) {
	time.Sleep(time.Duration(2) * time.Second)
	return "success", nil
}
