package promise_test

import (
	"context"
	"github.com/ahmadmuzakki29/promise"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPromiseDo(t *testing.T) {
	ctx := context.Background()
	timeout := time.Duration(2) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)

	val, err := promise.Promise(ctx, dbTx).Do()
	assert.Equal(t, "success", val.(string))
	assert.NoError(t, err)

	val, err = promise.Promise(ctx, httpTx).Do()
	// should be timeout now
	assert.Nil(t, val)
	assert.Equal(t, err, promise.TimeoutErr)

	cancel()
}

func TestPromiseThen(t *testing.T) {
	ctx := context.Background()
	timeout := time.Duration(2) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)

	promise.Promise(ctx, dbTx).Then(
		// it should be passed
		func(val interface{}) {
			assert.Equal(t, "success", val.(string))
		},
		func(err error) {
			t.Fatal()
		},
	)

	promise.Promise(ctx, httpTx).Then(
		func(val interface{}) {
			t.Fatal()
		},
		// it should be timeout
		func(err error) {
			assert.Equal(t, promise.TimeoutErr, err)
		},
	)

	cancel()
}

func dbTx() (interface{}, error) {
	time.Sleep(time.Duration(1) * time.Second)
	return "success", nil
}

func httpTx() (interface{}, error) {
	time.Sleep(time.Duration(2) * time.Second)
	return "success", nil
}
