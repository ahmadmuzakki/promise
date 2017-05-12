package promise

import (
	"context"
	"errors"
)

type Action func() (interface{}, error)

func Promise(ctx context.Context, action Action) *callback {
	success := make(chan interface{})
	fail := make(chan error)

	go func() {
		defer func() {
			close(success)
			close(fail)
		}()

		val, err := action()
		if err != nil {
			fail <- err
			return
		}
		success <- val
	}()

	return &callback{
		ctx.Done(), success, fail,
	}
}

type callback struct {
	done    <-chan struct{}
	success <-chan interface{}
	fail    <-chan error
}

type PromiseSuccess func(interface{})

type PromiseFail func(error)

func (c *callback) Then(success PromiseSuccess, fail PromiseFail) {
	defer c.drain()
	select {
	case <-c.done:
		fail(TimeoutErr)
	case val := <-c.success:
		success(val)
	case err := <-c.fail:
		fail(err)
	}
}

func (c *callback) Do() (interface{}, error) {
	defer c.drain()
	select {
	case <-c.done:
		return nil, TimeoutErr
	case val := <-c.success:
		return val, nil
	case err := <-c.fail:
		return nil, err
	}
}

// drain the channel if timeout hits first
// if the process doesn't timeout then it will do nothing
// because its already closed
func (c *callback) drain() {
	go func() { <-c.success }()
	go func() { <-c.fail }()
}

var TimeoutErr = errors.New("Operation Timeout")
