# channelcontext

Wraps an input channel into a Golang `context.Context` object.

### Behavior

The returned context will act according to the following rules:

1. It is fulfilled as soon as one of the components is fulfilled.
2. If it involves an error, the merged object will return the same error.
3. The context object has an additional method that returns the component index that was signalled.


## Example

```golang
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

    // We can retrieve the value that fulfilled the context if no error.
    if ctx.Err() != nil || ctx.DoneValue() != 5 {
        // Signal error
    }

    // ....
}
```

## LICENSE

See [LICENSE](/LICENSE) file for details.
