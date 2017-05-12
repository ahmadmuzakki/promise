# Promise

When you need to restrict execution time of an API you can make use of Go `context.Context` and set the timeout using `context.WithTimeout(Context,time.Duration)` but it's not very handy if you always listen to `ctx.Done()` channel everytime you need to execute function that could take long time.
With Promise you can easily wrap the desired function to respect the deadline of the context.

# Install
`go get -u github.com/ahmadmuzakki29/promise`

# Usage

```go
package main

import (
    "time"
    "fmt"
    "context"
    "github.com/ahmadmuzakki29/promise"
)

func main() {
	ctx := context.Background()
	// timeout for the whole operation is 2 seconds
	timeout := time.Duration(2) * time.Second
	ctx, done := context.WithTimeout(ctx, timeout)

	// return result with Do()
	val, err := promise.Promise(ctx, dbTx).Do()
	fmt.Println("dbTx result:",val,"error:",err)

	// or we can do callback style like this
	promise.Promise(ctx, httpTx).Then(func(val interface{}){
		fmt.Println("httpTx result:",val)
	},func(err error){
		fmt.Println("httpTx error:",err)
	})

	done()
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

```
