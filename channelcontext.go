package channelcontext

import (
	"context"
	"errors"
	"sync"
	"time"
)

// -----------------------------------------------------------------------------

// Context exposes a channel-backed completion signal through the context.Context interface.
type Context[T any] interface {
	context.Context

	// DoneValue returns the value received from the wrapped channel after completion.
	DoneValue() T
}

// -----------------------------------------------------------------------------

// ClosedChannel reports that the wrapped channel closed before producing a value.
var ClosedChannel = errors.New("closed channel")

// -----------------------------------------------------------------------------

type channelContext[T any] struct {
	ch         <-chan T
	lock       sync.RWMutex
	doneCh     chan struct{}
	doneValue  T
	cancelCh   chan struct{}
	cancelOnce sync.Once
	err        error
}

// -----------------------------------------------------------------------------

// New creates a channel-backed context and returns it with a cancel function.
func New[T any](ch <-chan T) (Context[T], context.CancelFunc) {
	if ch == nil {
		panic("channelcontext: nil channel")
	}

	// Create new context
	cc := channelContext[T]{
		ch:         ch,
		lock:       sync.RWMutex{},
		doneCh:     make(chan struct{}),
		cancelCh:   make(chan struct{}),
		cancelOnce: sync.Once{},
	}
	go cc.monitor()

	// Done
	return &cc, func() {
		cc.cancel()
	}
}

// Deadline reports that this context does not define a deadline.
func (cc *channelContext[T]) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done returns a channel that closes when the wrapped channel completes or the context is canceled.
func (cc *channelContext[T]) Done() <-chan struct{} {
	return cc.doneCh
}

// DoneValue returns the value received from the wrapped channel after completion.
func (cc *channelContext[T]) DoneValue() T {
	cc.lock.RLock()
	defer cc.lock.RUnlock()

	return cc.doneValue
}

// Err returns the completion error, if the wrapped channel closed or the context was canceled.
func (cc *channelContext[T]) Err() error {
	cc.lock.RLock()
	defer cc.lock.RUnlock()

	return cc.err
}

// Value reports that this context does not carry request-scoped values.
func (_ *channelContext[T]) Value(_ any) any {
	return nil
}

func (cc *channelContext[T]) monitor() {
	select {
	case v, ok := <-cc.ch:
		cc.lock.Lock()
		if ok {
			cc.doneValue = v
		} else {
			cc.err = ClosedChannel
		}
		cc.lock.Unlock()

	case <-cc.cancelCh:
		cc.lock.Lock()
		cc.err = context.Canceled
		cc.lock.Unlock()
	}

	close(cc.doneCh)
	cc.cancel()
}

func (cc *channelContext[T]) cancel() {
	cc.cancelOnce.Do(func() {
		close(cc.cancelCh)
	})
}
