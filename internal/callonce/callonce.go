// Package callonce implements concurrent call suppression and cache mechanism.
package callonce

import (
	"context"
	"sync"
)

// call is used to store information about in-flight data and result of a call of a user function.
type call struct {
	// result contains the value that was produced by calling the user function.
	result interface{}
	err    error

	// waiterCount is current count of callers of Do that are waiting for the result.
	// This value is only valid when done is false.
	waiterCount uint64

	// cancel is used to cancel the context provided to the user function.
	cancel context.CancelFunc

	// done is set to true once the user function has returned.
	done bool

	// doneChan is closed after done has been set to true.
	doneChan chan struct{}
}

// Group allows to call a function and cache its result once per key.
type Group struct {
	// mu protects the calls map
	mu sync.Mutex

	// calls contains information about all in-flight and completed calls, by key.
	// once a call is created for the given key it stays there, unless the call has been abandoned
	// before completing.
	calls map[string]*call
}

// Do calls the given fn and caches the returned value.
// If fn successfully returns a value for the given key and shouldCache is true, Do will return the same value forever
// without calling fn for the key again.
// If fn returns a value with shouldCache false, then only Do calls that spawned fn and calls that run concurrently to
// fn will receive the value.
// If the provided context is canceled Do returns immediately, but fn might continue running to fulfill
// concurrent requests from other goroutines.
// The context passed to fn is canceled in case all contexts passed to Do for a given key are canceled or after
// fn returns, whichever happens first. If the context passed to fn was canceled, Do won't cache the return
// value of fn.
func (g *Group) Do(ctx context.Context, key string, fn func(context.Context) (interface{}, error)) (interface{}, error) {
	g.mu.Lock()

	// We need to initialize the map the first time Group is used.
	if g.calls == nil {
		g.calls = make(map[string]*call)
	}

	// If no call for this group is in-flight or has completed, we need to spawn new call with its own
	// context.
	c, ok := g.calls[key]
	if !ok {
		c = &call{
			doneChan: make(chan struct{}),
		}
		g.calls[key] = c

		callCtx, cancel := context.WithCancel(context.Background())
		c.cancel = cancel

		go func() {
			defer cancel()
			result, err := fn(callCtx)
			g.mu.Lock()
			c.result = result
			c.err = err
			c.done = true
			close(c.doneChan)
			g.mu.Unlock()
		}()
	}

	// we are starting to wait for the call to complete.
	c.waiterCount++

	g.mu.Unlock()

	// Now wait until the call is complete or the caller abandons the wait.
	select {
	case <-ctx.Done():
		// The caller has abandoned the call.
		// If the the call is still running and nobody is waiting for it anymore, we need to cancel it.
		g.mu.Lock()
		if !c.done {
			c.waiterCount--
			if c.waiterCount == 0 {
				c.cancel()
				delete(g.calls, key)
			}
		}
		g.mu.Unlock()
		return nil, ctx.Err()
	case <-c.doneChan:
		// The call has successfully completed, return the result.
		return c.result, c.err
	}
}
