package channelcontext_test

import (
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

	if !errors.Is(ctx.Err(), channelcontext.ReceivedMessage) || ctx.DoneValue() != 5 {
		t.Fatalf("expected a value of 5 in the done value")
	}
}
