package channelcontext

import (
	"context"
	"errors"
	"sync"
	"time"
)

// -----------------------------------------------------------------------------

type Context[T any] interface {
	context.Context

	DoneValue() T
}

// -----------------------------------------------------------------------------

var ReceivedMessage = errors.New("received message")

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

// New creates a new context object from the given channel.
func New[T any](ch <-chan T) (Context[T], context.CancelFunc) {
	if ch == nil {
		return nil, nil
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

func (cc *channelContext[T]) Deadline() (deadline time.Time, ok bool) {
	return
}

func (cc *channelContext[T]) Done() <-chan struct{} {
	return cc.doneCh
}

func (cc *channelContext[T]) DoneValue() T {
	cc.lock.RLock()
	defer cc.lock.RUnlock()

	return cc.doneValue
}

func (cc *channelContext[T]) Err() error {
	cc.lock.RLock()
	defer cc.lock.RUnlock()

	return cc.err
}

func (_ *channelContext[T]) Value(_ any) any {
	return nil
}

func (cc *channelContext[T]) monitor() {
	select {
	case v := <-cc.ch:
		cc.lock.Lock()
		cc.doneValue = v
		cc.err = ReceivedMessage
		cc.lock.Unlock()

	case <-cc.cancelCh:
		cc.lock.Lock()
		cc.err = context.Canceled
		cc.lock.Unlock()
	}

	close(cc.doneCh)
	go cc.cancel()
}

func (cc *channelContext[T]) cancel() {
	cc.cancelOnce.Do(func() {
		close(cc.cancelCh)
	})
}
