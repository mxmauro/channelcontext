package channelcontext_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mxmauro/channelcontext"
)

// -----------------------------------------------------------------------------

func TestChannelCtx(t *testing.T) {
	intCh := make(chan int)
	ctx, cancelCtx := channelcontext.New[int](intCh)
	defer cancelCtx()

	go func() {
		intCh <- 5
	}()

	<-ctx.Done()

	if ctx.Err() != nil {
		t.Fatalf("expected nil error after receiving a value, got %v", ctx.Err())
	}
	if ctx.DoneValue() != 5 {
		t.Fatalf("expected a value of 5 in the done value")
	}
}

func TestChannelCtxClosedChannel(t *testing.T) {
	intCh := make(chan int)
	close(intCh)

	ctx, cancelCtx := channelcontext.New[int](intCh)
	defer cancelCtx()

	<-ctx.Done()

	if !errors.Is(ctx.Err(), channelcontext.ClosedChannel) {
		t.Fatalf("expected a closed channel error, got %v", ctx.Err())
	}
	if ctx.DoneValue() != 0 {
		t.Fatalf("expected zero value after channel close, got %d", ctx.DoneValue())
	}
}

func TestChannelCtxCancel(t *testing.T) {
	intCh := make(chan int)
	ctx, cancelCtx := channelcontext.New[int](intCh)

	cancelCtx()
	cancelCtx()

	<-ctx.Done()

	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatalf("expected canceled context, got %v", ctx.Err())
	}
	if ctx.DoneValue() != 0 {
		t.Fatalf("expected zero value after cancellation, got %d", ctx.DoneValue())
	}
}

func TestChannelCtxNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic for nil channel")
		}
	}()

	channelcontext.New[int](nil)
}
