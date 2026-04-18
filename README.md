# channelcontext

Wraps an input channel into a Golang `context.Context` object.

### Behavior

The returned context will act according to the following rules:

1. It is fulfilled when either a value is received from the wrapped channel, the wrapped channel is closed, or the returned cancel function is called.
2. When a value is received, `Err()` returns `nil` and `DoneValue()` returns the received value.
3. When the channel is closed before a value is received, `Err()` returns `ClosedChannel`.
4. When canceled, `Err()` returns `context.Canceled`.
5. Calling `New` with a `nil` channel panics.


## Example

```go
import (
    "github.com/mxmauro/channelcontext"
)

func main() {
    // Create a channel with an integer element
    intCh := make(chan int)

    // Create the context to wrap the channel
    ctx, cancelCtx := channelcontext.New[int](intCh)
    defer cancelCtx()

    // Send, in background, some value to the channel
    go func() {
        intCh <- 5
    }()

    // Wait for the context to be fulfilled
    <-ctx.Done()

    // We can retrieve the value that fulfilled the context when no error is reported.
    if ctx.Err() != nil || ctx.DoneValue() != 5 {
        // Signal error
    }

    // ....
}
```

## LICENSE

See [LICENSE](/LICENSE) file for details.
